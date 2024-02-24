package colorboxd

import (
	"cmp"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	prominentcolor "github.com/EdlinOrg/prominentcolor"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/joho/godotenv"
	"github.com/lucasb-eyer/go-colorful"
)

type Mode struct{ sort, test string }
type ImgSource struct{ url, local string }

var Modes = Mode{"sort", "test"}
var ImgSources = ImgSource{"url", "local"}

var mode *string
var imageSource *string

func init() {
	// Define global parameters to use during dev
	mode = &Modes.sort
	imageSource = &ImgSources.local

	functions.HTTP("AuthUserGetLists", HTTPAuthUserGetLists)
	functions.HTTP("AuthUser", HTTPAuthUser)
	functions.HTTP("GetLists", HTTPGetLists)
}

func HTTPAuthUserGetLists(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received call to HTTPAuthUserGetLists")

	// Handle CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	// Read env variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
		return
	}

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
		http.Error(w, "No access token in response", http.StatusInternalServerError)
		return
	}

	// Get Member id
	member, err := GetMemberId(accessTokenResponse.AccessToken)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting member id: %v", err), http.StatusInternalServerError)
		return
	}

	// Get User Lists
	userLists, err := GetUserLists(accessTokenResponse.AccessToken, member.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting user lists: %v", err), http.StatusInternalServerError)
		return
	}

	response := map[string]any{
		"token": accessTokenResponse,
		"user":  member,
		"lists": userLists,
	}

	// Return response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func HTTPAuthUser(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received call to HTTPAuthUser")

	// Handle CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	// Read env variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
		return
	}

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
		http.Error(w, "No access token in response", http.StatusInternalServerError)
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

	// Handle CORS
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")

	// Read env variables
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
		return
	}

	// Read authCode from query url - return error if not present
	accessToken := r.URL.Query().Get("accessToken")
	if accessToken == "" {
		http.Error(w, "Missing or empty 'authCode' query parameter", http.StatusBadRequest)
		return
	}

	// Get Member id
	userId := r.URL.Query().Get("userId")
	if err != nil {
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
	query := fmt.Sprintf("?client_id=%s&client_secret=%s", os.Getenv("LBOXD_KEY"), os.Getenv("LBOXD_SECRET")) // Handle this in a better way
	url := endpoint + query
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)} // Is this actually necessary?

	response, err := MakeHTTPRequest(method, url, nil, headers)
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
	query := fmt.Sprintf("?client_id=%s&client_secret=%s&member=%s&memberRelationship=Owner", os.Getenv("LBOXD_KEY"), os.Getenv("LBOXD_SECRET"), id) // Handle client key/secret in a better way
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

func SortListById() {
	// get list images
	getListImagesById()
	// sort
	// return a data structure with picture links, sorting number
}

func WriteSortedList() {
	// receive list id and sorting info from client
	// Write sorting info to users list
}

func getListImagesById() {
	// return slice of Image (or ImageInfo?)
}

func getTestImageUrls() ([]string, error) {
	var imageUrlSlice []string

	response, err := http.Get("https://picsum.photos/v2/list/")
	if err != nil {
		fmt.Printf("Error accessing Picsum API: %v\n", err)
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: Unexpected status code %d\n", response.StatusCode)
		return nil, err
	}

	res, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return nil, err
	}

	var responseData []map[string]interface{}
	err = json.Unmarshal(res, &responseData)
	if err != nil {
		fmt.Printf("Error decoding JSON data: %v\n", err)
		return nil, err
	}

	for _, image := range responseData {
		if url, ok := image["download_url"].(string); ok {
			imageUrlSlice = append(imageUrlSlice, url)
		}
	}

	return imageUrlSlice, nil
}

func getImagePaths(dir string) ([]string, error) {
	var imagePaths []string

	var fileTypes = []string{".png", ".jpg", ".jpeg"}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() {
			if slices.Contains[[]string](fileTypes, filepath.Ext(file.Name())) {
				imagePaths = append(imagePaths, filepath.Join(dir, file.Name()))
			}
		}
	}

	return imagePaths, nil
}

func LoadImage(path string) (image.Image, error) {
	var file io.ReadCloser
	var err error = nil
	switch *imageSource {
	case ImgSources.local:
		file, err = os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	case ImgSources.url:
		response, err := http.Get(path)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			return nil, err
		}
		file = response.Body
	default:
		fmt.Printf("Invalid image source setting.")
		return nil, nil
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func loadImageConcurrent(path string, wg *sync.WaitGroup, imageChan chan<- Image) {
	defer wg.Done()

	img, err := LoadImage(path)
	if err != nil {
		fmt.Printf("Error loading image %s: %v\n", path, err)
		return
	}

	image := Image{img: img, info: ImageInfo{Path: path}}

	imageChan <- image
}

func getDominantColor(k int, method int, img image.Image) string {
	resizeSize := uint(prominentcolor.DefaultSize) // This is a constant
	bgmasks := prominentcolor.GetDefaultMasks()    // This is a constant

	res, err := prominentcolor.KmeansWithAll(k, img, method, resizeSize, bgmasks)
	if err != nil {
		log.Println(err)
	}

	return res[0].AsString()
}

func getImageInfo(filePath string, img image.Image) (ImageInfo, error) {
	var method int = prominentcolor.ArgumentNoCropping // This is a constant

	hex := "#" + getDominantColor(3, method, img)
	color, _ := colorful.Hex(hex)
	hue, _, _ := color.Hsv()

	return ImageInfo{filePath, hex, color, hue}, nil
}

func getImageInfoConcurrent(filePath string, img image.Image, wg *sync.WaitGroup, colorChan chan<- ImageInfo) {
	defer wg.Done()

	imgColorInfo, err := getImageInfo(filePath, img)
	if err != nil {
		fmt.Printf("Error getting image color info for %s: %v\n", filePath, err)
		return
	}

	colorChan <- imgColorInfo
}

func Run() {
	start := time.Now()
	var imagePaths []string
	var err error = nil
	if *imageSource == ImgSources.url {
		imagePaths, err = getTestImageUrls()
	} else if *imageSource == ImgSources.local {
		imagePaths, err = getImagePaths("./images/example")
	}
	if err != nil {
		fmt.Printf("Error getting test Image URLs: %v\n", err)
		return
	}

	if *mode == Modes.sort {
		var imageSlice []ImageInfo
		var images []Image
		n := len(imagePaths)

		var wg sync.WaitGroup
		imageChan := make(chan Image, n)
		colorChan := make(chan ImageInfo, n)

		for _, path := range imagePaths {
			wg.Add(1)
			go loadImageConcurrent(path, &wg, imageChan)
		}

		go func() {
			wg.Wait()
			close(imageChan)
		}()

		for img := range imageChan {
			images = append(images, img)
		}

		if len(images) != n {
			fmt.Println("Warning! Images slice length does not match image paths slice length.")
			return
		}

		for i := 0; i < n; i++ {
			wg.Add(1)
			go getImageInfoConcurrent(images[i].info.Path, images[i].img, &wg, colorChan)
		}

		go func() {
			wg.Wait()
			close(colorChan)
		}()

		for imageInfo := range colorChan {
			imageSlice = append(imageSlice, imageInfo)
		}

		slices.SortFunc[[]ImageInfo](imageSlice, func(a, b ImageInfo) int { return cmp.Compare[float64](a.Hue, b.Hue) })

		fmt.Printf("Sorted %d images in: %s.\n", n, time.Since(start))

		createHTMLOutput(imageSlice)
	} else if *mode == Modes.test {
		CreateImageColorSummary(imagePaths)
	}
}

func createHTMLOutput(imageSlice []ImageInfo) {
	var buff strings.Builder

	buff.WriteString("<html><body><h1>Ordered Images</h1><table border=\"1\">")

	for _, img := range imageSlice {

		buff.WriteString("<tr><td><img src=\"" + img.Path + "\" width=\"200\" border=\"1\"></td>")
		buff.WriteString(fmt.Sprintf("<td style=\"background-color: %s;width:200px;height:50px;text-align:center;\">Color: %s</td></tr>", img.Hex, img.Hex))

	}

	buff.WriteString("</table></body><html>")
	if err := os.WriteFile("./sortedImages.html", []byte(buff.String()), 0644); err != nil {
		panic(err)
	}
}
