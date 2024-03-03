package colorboxd

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	prominentcolor "github.com/EdlinOrg/prominentcolor"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/lucasb-eyer/go-colorful"
)

func init() {
	functions.HTTP("AuthUser", HTTPAuthUser)
	functions.HTTP("GetLists", HTTPGetLists)
	functions.HTTP("SortList", HTTPSortListById)
	functions.HTTP("WriteList", HTTPWriteList)
}

func HTTPAuthUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received call to HTTPAuthUser")

	// Read env variables - local development only
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Printf("Could not load environment variables from .env file: %v\n", err)
	// 	return
	// }

	// Set necessary headers for CORS and cache policy
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("BASE_URL"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "private, max-age=3595") // Expire time of token (-5s for safety)

	// Read authCode from query url - return error if not present
	authCode := r.URL.Query().Get("authCode")
	if authCode == "" {
		http.Error(w, "Missing or empty 'authCode' query parameter", http.StatusBadRequest)
		return
	}

	// Get Access Token
	accessTokenResponse, err := GetAccessToken(authCode)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting access token: %v", err), http.StatusInternalServerError)
		return
	}
	if accessTokenResponse.AccessToken == "" {
		// NOTE - are we sure we want to error this out? Or just return empty?
		// http.Error(w, "No access token in response", http.StatusInternalServerError)
		return
	}

	member, err := GetMemberId(accessTokenResponse.AccessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting member id: %v", err), http.StatusInternalServerError)
		return
	}

	response := AuthUserResponse{
		Token:          accessTokenResponse.AccessToken,
		TokenType:      accessTokenResponse.TokenType,
		TokenRefresh:   accessTokenResponse.RefreshToken,
		TokenExpiresIn: accessTokenResponse.ExpiresIn,
		UserId:         member.ID,
		Username:       member.Username,
		UserGivenName:  member.GivenName,
	}

	// Return response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HTTPGetLists(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received call to HTTPGetLists")

	// Read env variables - local development only
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Printf("Could not load environment variables from .env file: %v\n", err)
	// 	return
	// }

	// Set necessary headers for CORS and cache policy
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("BASE_URL"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "private, max-age=3600")

	// Read accessToken from query url - return error if not present
	accessToken := r.URL.Query().Get("accessToken")
	if accessToken == "" {
		http.Error(w, "Missing or empty 'accessToken' query parameter", http.StatusBadRequest)
		return
	}

	// Get Member id
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "Missing or empty 'userId' query parameter", http.StatusBadRequest)
		return
	}

	// Get User Lists
	userLists, err := GetUserLists(accessToken, userId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting user lists: %v", err), http.StatusInternalServerError)
		return
	}

	// Return response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userLists)
}

func GetAccessToken(authCode string) (*AccessTokenResponse, error) {
	// Prepare endpoint and body for POST request
	method := "POST"
	endpoint := fmt.Sprintf("%s/auth/token", os.Getenv("LBOXD_BASEURL"))
	formData := url.Values{
		"grant_type":    {"authorization_code"},
		"code":          {authCode},
		"redirect_uri":  {os.Getenv("LBOXD_REDIRECT_URL")},
		"client_id":     {os.Getenv("LBOXD_KEY")},
		"client_secret": {os.Getenv("LBOXD_SECRET")},
	}
	headers := map[string]string{"Content-Type": "application/x-www-form-urlencoded", "Accept": "application/json"}

	response, err := MakeHTTPRequest(method, endpoint, strings.NewReader(formData.Encode()), headers)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Decode data into struct - handle error cases
	var responseData AccessTokenResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

func GetMemberId(token string) (*Member, error) {
	method := "GET"
	endpoint := fmt.Sprintf("%s/me", os.Getenv("LBOXD_BASEURL"))
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)} // Is this actually necessary?

	response, err := MakeHTTPRequest(method, endpoint, nil, headers)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseData map[string]json.RawMessage
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		fmt.Println(err)
		return nil, err
	}
	var member Member
	if err = json.Unmarshal(responseData["member"], &member); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &member, nil
}

func GetUserLists(token, id string) (*[]ListSummary, error) {
	method := "GET"
	endpoint := fmt.Sprintf("%s/lists", os.Getenv("LBOXD_BASEURL"))
	query := fmt.Sprintf("?member=%s&memberRelationship=Owner&perPage=100", id)
	url := endpoint + query
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)} // Is this actually necessary?

	response, err := MakeHTTPRequest(method, url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseData ListsResponse
	err = json.NewDecoder(response.Body).Decode(&responseData)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var lists []ListSummary = responseData.Items

	return &lists, nil
}

func HTTPSortListById(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received call to HTTPSortListById")

	// Read env variables - local development only
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Printf("Could not load environment variables from .env file: %v\n", err)
	// 	return
	// }

	// Set necessary headers for CORS and cache policy
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("BASE_URL"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Cache-Control", "private, max-age=3600")

	// Read accessToken from query url - return error if not present
	accessToken := r.URL.Query().Get("accessToken")
	if accessToken == "" {
		http.Error(w, "Missing or empty 'accessToken' query parameter", http.StatusBadRequest)
		return
	}

	// Get List id
	listId := r.URL.Query().Get("listId")
	if listId == "" {
		http.Error(w, "Missing or empty 'listId' query parameter", http.StatusBadRequest)
		return
	}

	// Get Entries from List
	listEntries, err := GetListEntries(accessToken, listId)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting entries from list: %v", err), http.StatusInternalServerError)
		return
	}

	entriesWithImageInfo, err := processListImages(listEntries)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error processing posters for list entries: %v", err), http.StatusInternalServerError)
		return
	}

	slices.SortFunc[[]Entry](*entriesWithImageInfo, func(a, b Entry) int { return cmp.Compare[float64](a.ImageInfo.Hue, b.ImageInfo.Hue) })

	response := map[string][]Entry{
		"items": *entriesWithImageInfo,
	}

	// Return response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func GetListEntries(token, id string) (*[]Entry, error) {
	method := "GET"
	endpoint := fmt.Sprintf("%s/list/%s/entries", os.Getenv("LBOXD_BASEURL"), id)
	query := "?perPage=100"
	url := endpoint + query
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	response, err := MakeHTTPRequest(method, url, nil, headers)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	// Raw decode
	var responseData map[string]json.RawMessage
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Cast items to type
	var listEntriesData []ListEntriesResponse
	if err = json.Unmarshal(responseData["items"], &listEntriesData); err != nil {
		fmt.Println(err)
		return nil, err
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

	for _, entry := range *listEntries {
		wg.Add(1)
		go loadImageConcurrent(entry, &wg, imageChan)
	}

	go func() {
		wg.Wait()
		close(imageChan)
	}()

	for img := range imageChan {
		images = append(images, img)
	}

	if len(images) != n {
		err := errors.New("Error: Images slice length does not match image paths slice length.")
		return nil, err
	}

	for i := 0; i < n; i++ {
		wg.Add(1)
		go getImageInfoConcurrent(images[i], &wg, colorChan)
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

func HTTPWriteList(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received call to HTTPWriteList")

	// Read env variables - local development only
	// err := godotenv.Load()
	// if err != nil {
	// 	fmt.Printf("Could not load environment variables from .env file: %v\n", err)
	// 	return
	// }

	// Set necessary headers for CORS
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("BASE_URL"))
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var responseData WriteListRequest
	err := json.NewDecoder(r.Body).Decode(&responseData)
	if err != nil {
		fmt.Println("Here")

		http.Error(w, "Missing or empty 'accessToken' query parameter", http.StatusBadRequest)
		return
	}

	listUpdateRequest, err := prepareListUpdateRequest(responseData.List, responseData.Offset, "")
	if err != nil {
		http.Error(w, "Error preparing list update request body", http.StatusBadRequest)
		return
	}

	message, err := writeListSorting(responseData.AccessToken, responseData.List.ID, *listUpdateRequest)
	if err != nil {
		http.Error(w, "Error updating user list", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

// Taking a sorted list, return a ListUpdateRequest, as required by Letterboxd endpoint
func prepareListUpdateRequest(list ListWithEntries, offset int, sortMethod string) (*ListUpdateRequest, error) {
	request := ListUpdateRequest{Version: list.Version}

	//
	//
	//
	//
	// THIS STILL CONTAINS A BUG. Sorting as seen on colorboxd doesn't always match how the list ends up
	//
	//
	//
	//

	// n := len(list.Entries)
	currentPositions := make(map[string]int)
	finishPositions := make(map[string]int)
	for i, entry := range list.Entries {
		initPos, err := strconv.Atoi(entry.EntryID)
		if err != nil {
			fmt.Println("Error parsing entryId to int")
			return nil, err
		}
		// endPos := (i - offset) % n // FIX: Offset doesn't seem to be working correctly

		currentPositions[entry.FilmID] = initPos
		finishPositions[entry.FilmID] = i
	}

	updateEntries := []listUpdateEntry{}

	for film, finPos := range finishPositions {
		currPos := currentPositions[film]

		delta := finPos - currPos

		updateEntries = append(updateEntries, listUpdateEntry{Action: "UPDATE", Position: currPos, NewPosition: finPos})

		currentPositions[film] = finPos

		for f, cP := range currentPositions {
			if f != film {
				if delta > 0 && currPos < cP && cP <= finPos {
					currentPositions[f]--
				} else if delta < 0 && finPos <= cP && cP < currPos {
					currentPositions[f]++
				}
			}
		}
	}

	request.Entries = updateEntries
	return &request, nil
}

// Send request to Letterboxd endpoint to update list
func writeListSorting(token, id string, listUpdateRequest ListUpdateRequest) (*[]string, error) {

	// Prepare endpoint and body for PATCH request
	method := "PATCH"
	endpoint := fmt.Sprintf("%s/list/%s", os.Getenv("LBOXD_BASEURL"), id)
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token), "Content-Type": "application/json", "X-HTTP-Method-Override": "PATCH"}
	body, err := json.Marshal(listUpdateRequest)
	if err != nil {
		return nil, err
	}

	response, err := MakeHTTPRequest(method, endpoint, bytes.NewReader(body), headers)
	if err != nil {
		return nil, err
	}
	fmt.Println(response.Status)
	defer response.Body.Close()

	var responseData ListUpdateResponse
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(responseData.Messages)

	var message []string
	if len(responseData.Messages) != 0 {
		for _, m := range responseData.Messages {
			message = append(message, fmt.Sprintf("%s: %s - %s", m.Type, m.Code, m.Title))
		}
		errorStr := "The letterboxd API responded with the following errors: " + strings.Join(message, "; ")
		return &message, errors.New(errorStr)
	}

	//
	//
	// Consider using the ListUpdateMessage type, even for a successful response, to
	// maintain standard response type between errors and successes.
	//
	//
	message = []string{"List updated successfully"}

	return &message, nil
}

func LoadImage(path string) (image.Image, error) {
	var file io.ReadCloser
	var err error = nil

	response, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, err
	}
	file = response.Body

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func loadImageConcurrent(entry Entry, wg *sync.WaitGroup, imageChan chan<- Image) {
	defer wg.Done()

	img, err := LoadImage(entry.ImageInfo.Path)
	if err != nil {
		fmt.Printf("Error loading image %s: %v\n", entry.ImageInfo.Path, err)
		return
	}

	image := Image{img: img, info: entry}

	imageChan <- image
}

func getDominantColor(k int, method int, img image.Image) (*string, error) {
	resizeSize := uint(prominentcolor.DefaultSize)
	// bgmasks := prominentcolor.GetDefaultMasks() // Default masks (black,white or green backgrounds)
	bgmasks := []prominentcolor.ColorBackgroundMask{} // No mask
	// Think it's best not to use a mask here. Can return a slice of dominant colours from
	// this function, then, if the most dominant is very high/low luminosity, for example,
	// can instead use the second most dominant colour

	res, err := prominentcolor.KmeansWithAll(k, img, method, resizeSize, bgmasks)
	if err != nil {
		return nil, err
	}

	stringResponse := res[0].AsString()
	// IMPORTANT - here we are choosing only the single most dominant colour. Might be useful to pass the three most dominant colours.
	// Then, we can round all hue values, and have three sub hue values. If any hues are directly equal, can use the second-most
	// dominant colour to determine how to sort the images.
	return &stringResponse, nil
}

func getImageInfo(entry Entry, img image.Image) (*Entry, error) {
	var method int = prominentcolor.ArgumentNoCropping // This is a constant

	domColor, err := getDominantColor(3, method, img)
	if err != nil {
		return nil, err
	}

	hex := "#" + *domColor
	color, _ := colorful.Hex(hex)
	hue, _, _ := color.Hsv()

	entry.ImageInfo.Hex = hex
	entry.ImageInfo.Color = color
	entry.ImageInfo.Hue = hue

	return &entry, nil
}

func getImageInfoConcurrent(image Image, wg *sync.WaitGroup, colorChan chan<- Entry) {
	defer wg.Done()

	imgColorInfo, err := getImageInfo(image.info, image.img)
	if err != nil {
		fmt.Printf("Error getting image color info for poster for %s: %v\n", image.info.Name, err)
		return
	}

	colorChan <- *imgColorInfo
}
