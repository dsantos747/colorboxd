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
	"github.com/joho/godotenv"
	"github.com/lucasb-eyer/go-colorful"
)

func HTTPSortListById(w http.ResponseWriter, r *http.Request) {
	var err error

	// Read env variables - local development only - comment out for production
	err = godotenv.Load()
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

	entriesWithImageInfo, err := processListImages(listEntries)
	if err != nil {
		ReturnError(w, fmt.Errorf("failed to process posters for list entries: %w", err).Error(), http.StatusInternalServerError)
		return
	}

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

func processListImages(listEntries *[]Entry) (*[]Entry, error) {
	var entrySlice []Entry
	var images []Image
	n := len(*listEntries)

	var wg sync.WaitGroup
	imageChan := make(chan Image, n)
	colorChan := make(chan Entry, n)
	errChan := make(chan error, n)

	for _, entry := range *listEntries {
		wg.Add(1)
		go concurrentLoadImage(entry, &wg, imageChan, errChan)
	}

	go func() {
		wg.Wait()
		close(imageChan)
	}()

	for img := range imageChan {
		images = append(images, img)
	}

	select {
	case err := <-errChan:
		return nil, err
	default:
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go concurrentGetImageInfo(images[i], &wg, colorChan, errChan)
	}

	go func() {
		wg.Wait()
		close(colorChan)
	}()

	for entry := range colorChan {
		entrySlice = append(entrySlice, entry)
	}

	return &entrySlice, nil
}

func assignListRankings(listEntries *[]Entry) (*[]Entry, error) {
	for i, e := range *listEntries {
		(*listEntries)[i].SortVals.Hue = AlgoHue(e.ImageInfo.Colors)
		(*listEntries)[i].SortVals.Lum = AlgoLuminosity(e.ImageInfo.Colors)
		(*listEntries)[i].SortVals.BrightDomHue = AlgoBrightDominantHue(e.ImageInfo.Colors)
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

func loadImage(path string) (image.Image, error) {
	var err error = nil

	response, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, err
	}

	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func concurrentLoadImage(entry Entry, wg *sync.WaitGroup, imageChan chan<- Image, errChan chan<- error) {
	defer wg.Done()

	img, err := loadImage(entry.ImageInfo.Path)
	if err != nil {
		errChan <- fmt.Errorf("error loading image %s: %v", entry.ImageInfo.Path, err)
		return
	}

	image := Image{img: img, info: entry}

	imageChan <- image
}

func getDominantColors(k, method int, img image.Image) (*[]prominentcolor.ColorItem, error) {
	resizeSize := uint(prominentcolor.DefaultSize)
	var bgmasks []prominentcolor.ColorBackgroundMask // No mask
	// bgmasks := prominentcolor.GetDefaultMasks() // Default masks (black,white or green backgrounds)

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

func concurrentGetImageInfo(image Image, wg *sync.WaitGroup, colorChan chan<- Entry, errChan chan<- error) {
	defer wg.Done()

	imgColorInfo, err := getImageInfo(image.info, image.img)
	if err != nil {
		errChan <- fmt.Errorf("error getting image color info for poster for %s: %v", image.info.Name, err)
		return
	}

	colorChan <- *imgColorInfo
}
