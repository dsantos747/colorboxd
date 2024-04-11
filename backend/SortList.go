package colorboxd

import (
	"cmp"
	"encoding/json"
	"fmt"
	"image"
	"net/http"
	"os"
	"slices"
	"sync"

	// Accepted image formats in loadImage
	_ "image/jpeg"
	_ "image/png"

	prominentcolor "github.com/EdlinOrg/prominentcolor"
	"github.com/disintegration/imaging"
	"github.com/lucasb-eyer/go-colorful"
)

// Endpoint for determining the sorting of a list, given a list id
func HTTPSortListById(w http.ResponseWriter, r *http.Request) {
	var err error

	// Read env variables
	err = LoadEnv()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
		return
	}

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
	listEntries, err := getListEntries(accessToken, listId)
	if err != nil {
		ReturnError(w, fmt.Errorf("failed to retrieve entries from list: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	// start := time.Now()
	entriesWithImageInfo, err := processListImages(listEntries)
	if err != nil {
		ReturnError(w, fmt.Errorf("failed to process posters for list entries: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	// fmt.Println("Process list images took: ", time.Since(start))

	entriesWithRanking, err := assignListRankings(entriesWithImageInfo)
	if err != nil {
		ReturnError(w, fmt.Errorf("failed assigning sort rankings for list: %w", err).Error(), http.StatusInternalServerError)
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

// For a given list id, returns a slice of each entry in the list
func getListEntries(token, id string) (*[]Entry, error) {
	nextCursor := "start=0"
	method := "GET"
	endpoint := fmt.Sprintf("%s/list/%s/entries", os.Getenv("LBOXD_BASEURL"), id)
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}
	var listEntriesData []ListEntries

	// Loop until no "next" pagination cursor is present in response
	for len(nextCursor) > 0 {
		query := fmt.Sprintf("?cursor=%s&perPage=100", nextCursor)
		url := endpoint + query

		response, err := MakeHTTPRequest(method, url, nil, headers)
		if err != nil {
			return nil, fmt.Errorf("error making HTTP request: %v", err)
		}
		defer response.Body.Close()

		var responseData ListEntriesResponse
		if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
			fmt.Println(err)
			return nil, fmt.Errorf("error decoding letterboxd list entries JSON response: %v", err)
		}

		listEntriesData = append(listEntriesData, responseData.Items...)
		nextCursor = responseData.Next
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

		entries[i] = Entry{
			EntryID:            item.EntryID,
			FilmID:             item.Film.ID,
			Name:               item.Film.Name,
			ReleaseYear:        item.Film.ReleaseYear,
			Adult:              item.Film.Adult,
			PosterCustomisable: item.Film.PosterCustomisable,
			PosterURL:          item.Film.Poster.Sizes[0].URL,
			AdultPosterURL:     adultUrl,
			ImageInfo:          ImageInfo{Path: imgPath},
		}
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
		// Fetch image and process
		img, err := loadImage(image.info.ImageInfo.Path)
		if err != nil {
			errChan <- fmt.Errorf("error loading image %s: %v", image.info.ImageInfo.Path, err)
			wg.Done()
			continue
		}

		image.img = img

		entry, err := getImageInfo(image.info, image.img)
		if err != nil {
			errChan <- fmt.Errorf("error getting image color info for poster for %s: %v", image.info.Name, err)
			wg.Done()
			continue
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
		return nil, err
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
