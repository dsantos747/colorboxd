package main

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
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/rs/cors"

	prominentcolor "github.com/EdlinOrg/prominentcolor"
	"github.com/lucasb-eyer/go-colorful"
)

type RequestBody struct{}

type ImageInfo struct {
	Path  string
	Hex   string
	Color colorful.Color
	Hue   float64
}

type Image struct {
	img  image.Image
	info ImageInfo
}

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

	functions.HTTP("HelloAOC", HelloWorld)
	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloWorld)

	// Use CORS middleware
	handler := cors.Default().Handler(mux)

	// Start the server
	http.ListenAndServe(":8080", handler)
}

// This is a simple tracer bullet function deployed to Cloud Functions
// Use the same process to create API route endpoints
func HelloWorld(w http.ResponseWriter, r *http.Request) {

	// var reqBody RequestBody
	// decoder := json.NewDecoder(r.Body)
	// err := decoder.Decode(&reqBody)
	// if err != nil {
	// 	http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
	// 	return
	// }

	response := map[string]interface{}{
		"world":    "Hello World!",
		"darkness": "My old friend",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	fmt.Fprintf(w, "Hello world!")
}

func AuthUser() {
	// Access users info with Auth0
	// If possible, do this only once and pass an auth token around?
}

func GetUserLists() {
	// Hit up letterboxd api to get all users lists
	// Return: id, name, date modified
}

func SortListById() {
	// get list images
	// sort
	// return a data structure with picture links, sorting number
}

func WriteSortedList() {
	// receive list id and sorting info from client
	// Write sorting info to users list
}

func getListImagesById() {}

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

func main() {
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
