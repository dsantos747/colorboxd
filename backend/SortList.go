package colorboxd

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
	"sync"

	// Accepted image formats in loadImage
	_ "image/jpeg"
	_ "image/png"

	"github.com/dsantos747/letterboxd_hue_sort/backend/redis"

	prominentcolor "github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
	"go.uber.org/ratelimit"
	"golang.org/x/sync/errgroup"
)

var rc redis.Redis

// HTTPSortListById is the serverless function for computing the color information of each movie poster in
// a user's Letterboxd list and consequently computing the different sort rankings.
func HTTPSortListById(w http.ResponseWriter, r *http.Request) {
	var err error
	ctx := context.Background() // Hack for now

	l := slog.Default()

	// Read env variables
	err = LoadEnv()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
		return
	}

	rc = redis.New(os.Getenv("REDIS_URL"))

	// Set necessary headers for CORS and cache policy
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("BASE_URL"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "private, max-age=3600")

	// Read accessToken from query url - return error if not present
	accessToken := r.URL.Query().Get("accessToken")
	if accessToken == "" {
		ReturnError(w, "Missing or empty 'accessToken' query parameter", http.StatusBadRequest)
		return
	}

	// Get List id
	listId := r.URL.Query().Get("listId")
	if listId == "" {
		ReturnError(w, "Missing or empty 'listId' query parameter", http.StatusBadRequest)
		return
	}

	// Get Entries from List
	listEntries, err := getListEntries(ctx, accessToken, listId)
	if err != nil {
		l.Error("failed to retrieve entries from list", "err", err)
		ReturnError(w, "failed to retrieve entries from list", http.StatusInternalServerError)
		return
	}

	entriesWithImageInfo, err := processListImagesV3(ctx, listEntries)
	if err != nil {
		l.Error("failed to process posters for list entries", "err", err)
		ReturnError(w, "failed to process posters for list entries", http.StatusInternalServerError)
		return
	}

	entriesWithRanking, err := assignListRankings(entriesWithImageInfo)
	if err != nil {
		l.Error("failed assigning sort rankings for list", "err", err)
		ReturnError(w, "failed assigning sort rankings for list", http.StatusInternalServerError)
		return
	}

	slices.SortFunc[[]Entry](*entriesWithRanking, func(a, b Entry) int {
		return cmp.Compare[int](AlgoHue(a.ImageInfo.Colors), AlgoHue(b.ImageInfo.Colors))
	})

	response := map[string][]Entry{
		"items": *entriesWithRanking,
	}

	// Return response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getFilmCount(token, id string) (int, error) {
	method := "GET"
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
	endpoint := fmt.Sprintf("%s/list/%s", os.Getenv("LBOXD_BASEURL"), id)

	response, err := MakeHTTPRequest(method, endpoint, nil, headers)
	if err != nil {
		return 0, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer response.Body.Close()

	var responseData List
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return 0, fmt.Errorf("error decoding letterboxd list metadata JSON response: %v", err)
	}

	return int(responseData.FilmCount), nil
}

// For a given list id, returns a slice of each entry in the list
func getListEntries(ctx context.Context, token, id string) (*[]Entry, error) {
	method := "GET"
	endpoint := fmt.Sprintf("%s/list/%s/entries", os.Getenv("LBOXD_BASEURL"), id)
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
	var listEntriesData []ListEntries

	filmCount, err := getFilmCount(token, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get list length: %w", err)
	}

	errGroup, ctx := errgroup.WithContext(ctx)
	mu := sync.Mutex{}

	perPage := 100
	curr := 0
	done := false

	// Loop until no "next" pagination cursor is present in response
	for filmCount > 0 {
		query := fmt.Sprintf("?cursor=%s&perPage=%d", fmt.Sprintf("start=%d", curr), perPage)
		url := endpoint + query

		errGroup.Go(func() error {
			response, err := MakeHTTPRequest(method, url, nil, headers)
			if err != nil {
				return fmt.Errorf("error making HTTP request: %v", err)
			}
			defer response.Body.Close()

			var responseData ListEntriesResponse
			if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
				return fmt.Errorf("error decoding letterboxd list entries JSON response: %v", err)
			}

			mu.Lock()
			listEntriesData = append(listEntriesData, responseData.Items...)
			mu.Unlock()

			if len(responseData.Next) == 0 { // check we have done them all. Note this might not be the last routine to be processed
				done = true
			}
			return nil
		})

		filmCount -= perPage
		curr += perPage
	}

	err = errGroup.Wait()
	if err != nil {
		return nil, err
	}
	if !done {
		return nil, fmt.Errorf("failed to retrieve all list entries%s", "")
	}

	// Extract relevant info from each item into []Entry format
	n := len(listEntriesData)
	entries := make([]Entry, n)
	var adultUrl, imgPath string
	for i, item := range listEntriesData {
		adultUrl = ""
		imgPath = item.Film.Poster.Sizes[0].URL
		if item.Film.Adult {
			adultUrl = item.Film.AdultPoster.Sizes[0].URL
			imgPath = adultUrl
		}

		parsedURL, err := url.Parse(imgPath)
		if err != nil {
			return nil, fmt.Errorf("failed to parse img url: %w", err)
		}
		version := parsedURL.Query().Get("v")
		if version == "" {
			return nil, fmt.Errorf("failed to extract version from img url")
		}

		entries[i] = Entry{
			ListPosition:       i,
			EntryID:            item.EntryID,
			FilmID:             item.Film.ID,
			Name:               item.Film.Name,
			ReleaseYear:        item.Film.ReleaseYear,
			Adult:              item.Film.Adult,
			PosterCustomisable: item.Film.PosterCustomisable,
			PosterURL:          item.Film.Poster.Sizes[0].URL,
			AdultPosterURL:     adultUrl,
			CacheKey:           fmt.Sprintf("%s_%s", item.Film.ID, version), // underscore is important in key format
			ImageInfo:          ImageInfo{Path: imgPath},
		}
	}

	return &entries, nil
}

func processListImagesV3(ctx context.Context, listEntries *[]Entry) (*[]Entry, error) {
	// First we query Redis
	keys := []string{}
	for _, entry := range *listEntries {
		keys = append(keys, entry.CacheKey)
	}

	res, err := rc.GetBatch(keys)
	if err != nil {
		return nil, fmt.Errorf("failed to lookup keys in redis: %w", err)
	}

	// We pass through and append all cache hits
	var entries, entriesToLoad []Entry
	for _, e := range *listEntries {
		entry := e

		// Append entries fetched from cache
		if res[entry.CacheKey].Hit {
			entry.ImageInfo.Colors = parseColors(res[entry.CacheKey].Colors, res[entry.CacheKey].Counts)
			entries = append(entries, entry)
			continue
		}

		entriesToLoad = append(entriesToLoad, entry)
	}

	// Then we go through the process of fetch images that we are missing
	errGroup, ctx := errgroup.WithContext(ctx)
	rl := ratelimit.New(500)
	mu := sync.Mutex{}
	rlCtx, rlCancel := context.WithCancel(ctx) // This is a hack to cancel all goroutines if we get rate-limited when loading images

	var c_keys []string
	var c_colors [][]string
	var c_counts [][]int
	for _, e := range entriesToLoad {
		// Process any entries not available in cache
		errGroup.Go(func() error {
			rl.Take()

			// This block cancels all goroutines if we're getting rate-limited
			select {
			case <-rlCtx.Done():
				return rlCtx.Err()
			default:
			}

			img, err := loadImage(e.ImageInfo.Path)
			if err != nil {
				if strings.Contains(err.Error(), "error fetching image from letterboxd servers") {
					rlCancel()
				}
				return fmt.Errorf("error loading image %s: %v", e.ImageInfo.Path, err)
			}

			entry, err := getImageInfo(e, img)
			if err != nil {
				return fmt.Errorf("error getting image color info for poster for %s: %v", entry.Name, err)
			}

			colors, counts := []string{}, []int{}
			for _, c := range entry.ImageInfo.Colors {
				colors = append(colors, c.hex)
				counts = append(counts, c.count)
			}

			mu.Lock()
			c_keys = append(c_keys, entry.CacheKey)
			c_colors = append(c_colors, colors)
			c_counts = append(c_counts, counts)

			entries = append(entries, *entry)
			mu.Unlock()

			return nil
		})
	}

	egErr := errGroup.Wait()
	if len(keys) > 0 { // Even if we fail to process all, set to cache what we did manage
		go func() {
			rc.SetBatch(c_keys, c_colors, c_counts)
		}()
	}
	if egErr != nil { // Then handle the error
		return nil, egErr
	}

	return &entries, nil
}

// This v2 method bypasses the whole worker pattern and just uses a good old errgroup. NEEDS TO BE TESTED
func processListImagesV2(listEntries *[]Entry) (*[]Entry, error) {
	var entries []Entry

	ctx := context.Background() // Hack for now

	errGroup, ctx := errgroup.WithContext(ctx)

	for _, e := range *listEntries {
		errGroup.Go(func() error {

			res, err := rc.Get(e.CacheKey)
			if err != nil {
				return fmt.Errorf("failed to fetch from redis: %w", err)
			}

			if res.Hit {
				entry := e
				entry.ImageInfo.Colors = parseColors(res.Colors, res.Counts)
				entries = append(entries, entry)
			} else {
				img, err := loadImage(e.ImageInfo.Path)
				if err != nil {
					return fmt.Errorf("error loading image %s: %v", e.ImageInfo.Path, err)
				}

				entry, err := getImageInfo(e, img)
				if err != nil {
					return fmt.Errorf("error getting image color info for poster for %s: %v", entry.Name, err)
				}

				go func() {
					colors, counts := []string{}, []int{}
					for _, c := range entry.ImageInfo.Colors {
						colors = append(colors, c.hex)
						counts = append(counts, c.count)
					}
					rc.Set(entry.CacheKey, colors, counts)
				}()

				entries = append(entries, *entry)
			}
			return nil
		})
	}

	if err := errGroup.Wait(); err != nil {
		return nil, err
	}

	return &entries, nil
}

// For a slice of entries, this creates some goroutines which download the poster and extract colour
// information for each film. workerCount can be used to adjust the amount of goroutines.
func processListImages(listEntries *[]Entry) (*[]Entry, error) {
	var entrySlice []Entry
	n := len(*listEntries)

	var wg sync.WaitGroup
	imageChan := make(chan Image, n)
	colorChan := make(chan Entry, n)
	errChan := make(chan error, n)
	workerCount := 200 // Consider adjusting this based on list size; ask letterboxd team about rate limiting on their servers

	for i := 0; i < workerCount; i++ {
		go worker(imageChan, colorChan, &wg, errChan)
	}

	for _, entry := range *listEntries {
		wg.Add(1)
		imageChan <- Image{info: entry}
	}

	go func() {
		wg.Wait()
		close(imageChan)
		close(colorChan)
	}()

	for entry := range colorChan {
		entrySlice = append(entrySlice, entry)
	}

	return &entrySlice, nil
}

// This worker (pool size limited by workerCount) listens on imageChan, downloads the image and
// extracts the colour information, then returns the populated Entry to colorChan
func worker(imageChan <-chan Image, colorChan chan<- Entry, wg *sync.WaitGroup, errChan chan<- error) {
	for image := range imageChan {
		// Here need to first check redis cache for image info
		entry := &image.info

		res, err := rc.Get(image.info.CacheKey)
		if err != nil {
			errChan <- fmt.Errorf("failed to fetch from redis: %w", err)
			continue
		}

		if res.Hit {
			entry.ImageInfo.Colors = parseColors(res.Colors, res.Counts)
		} else {
			// Fetch image and process
			img, err := loadImage(image.info.ImageInfo.Path)
			if err != nil {
				errChan <- fmt.Errorf("error loading image %s: %v", image.info.ImageInfo.Path, err)
				wg.Done()
				continue
			}

			image.img = img

			entry, err = getImageInfo(image.info, image.img)
			if err != nil {
				errChan <- fmt.Errorf("error getting image color info for poster for %s: %v", image.info.Name, err)
				wg.Done()
				continue
			}

			// Set to redis
			// todo, NEED to ensure that the key we are using is the id of the film poster - not the id of the film itself
			// The letterboxd api has a posterPickerUrl, which is to do with the custom poster chosen in a list. Could we use this?
			// Maybe, check if that field is empty or not when the poster is standard, that could be useful
			// Once you have the api key back, can use that to make some postman calls and check
			colors, counts := []string{}, []int{}
			for _, c := range entry.ImageInfo.Colors {
				colors = append(colors, c.hex)
				counts = append(counts, c.count)
			}
			rc.Set(image.info.CacheKey, colors, counts)
		}

		colorChan <- *entry
		wg.Done()
	}
}

// Download and resize an image, given a source url
func loadImage(path string) (image.Image, error) {
	var err error = nil

	response, err := MakeHTTPRequest("GET", path, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error fetching image from letterboxd servers: %w", err)
	}
	defer response.Body.Close()

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	smallImg := imaging.Resize(img, 80, 0, imaging.NearestNeighbor)

	return smallImg, nil
}

// Populate an image with information about its dominant colours
func getImageInfo(entry Entry, img image.Image) (*Entry, error) {
	var method int = prominentcolor.ArgumentNoCropping

	domColors, err := getDominantColors(3, method, img)
	if err != nil {
		return nil, err
	}

	var currColor Color
	var colors []Color

	for _, c := range *domColors {
		hex := "#" + c.AsString()
		rgb, _ := colorful.Hex(hex) // This feels a bit backwards, going from rgb to hex to rgb
		hue, sat, lum := rgb.Hsl()
		_, _, val := rgb.Hsv() // Look into docs on using Clamped rgb values before converting to hsl/hsv

		currColor = Color{rgb: rgb, hex: hex, h: hue, s: sat, l: lum, v: val, count: c.Cnt}
		colors = append(colors, currColor)
	}

	entry.ImageInfo.Colors = colors
	return &entry, nil
}

// Run the k-means method to extract the top 3 dominant colours from a poster
func getDominantColors(k, method int, img image.Image) (*[]prominentcolor.ColorItem, error) {
	resizeSize := uint(1000) // larger to prevent re-resizing (we've already resized)
	// resizeSize := uint(prominentcolor.DefaultSize)

	var bgmasks []prominentcolor.ColorBackgroundMask // No mask

	res, err := prominentcolor.KmeansWithAll(k, img, method, resizeSize, bgmasks)
	if err != nil {
		return nil, err
	}

	// Limit to top 3 colors - Increasing this limit may cause "a lotta damage"
	if len(res) > 3 {
		res = res[0:2]
	}

	return &res, nil
}

// This function calculates each poster's ranking according to each sort method (see sortAlgorithms file)
func assignListRankings(listEntries *[]Entry) (*[]Entry, error) {
	for i, e := range *listEntries {
		(*listEntries)[i].SortVals.Hue = AlgoHue(e.ImageInfo.Colors)
		(*listEntries)[i].SortVals.Lum = AlgoLuminosity(e.ImageInfo.Colors)
		(*listEntries)[i].SortVals.InverseStep_8 = AlgoInverseStep(e.ImageInfo.Colors, 8)
		(*listEntries)[i].SortVals.InverseStep_12 = AlgoInverseStep(e.ImageInfo.Colors, 12)
		(*listEntries)[i].SortVals.InverseStep2_8 = AlgoInverseStepV2(e.ImageInfo.Colors, 8)
		(*listEntries)[i].SortVals.InverseStep2_12 = AlgoInverseStepV2(e.ImageInfo.Colors, 12)
		(*listEntries)[i].SortVals.BRBW1 = AlgoBRBW1(e.ImageInfo.Colors)
		(*listEntries)[i].SortVals.BRBW2 = AlgoBRBW2(e.ImageInfo.Colors)
	}

	// where error handling?

	return listEntries, nil
}

func parseColors(hexes []string, counts []int) []Color {
	var colors []Color
	for i, hex := range hexes {
		rgb, _ := colorful.Hex(hex)
		hue, sat, lum := rgb.Hsl()
		_, _, val := rgb.Hsv()

		colors = append(colors, Color{rgb: rgb, hex: hex, h: hue, s: sat, l: lum, v: val, count: counts[i]})
	}

	return colors
}
