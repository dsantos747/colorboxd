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

	prominentcolor "github.com/EdlinOrg/prominentcolor"
	"github.com/lucasb-eyer/go-colorful"
)

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

func loadImage(path string) (image.Image, error) {
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

	img, err := loadImage(path)
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

func outputColorRange(colorRange []prominentcolor.ColorItem) string { // WILL BE DELETED
	var buff strings.Builder
	buff.WriteString("<table><tr>")
	for _, color := range colorRange {
		buff.WriteString(fmt.Sprintf("<td style=\"background-color: #%s;width:200px;height:50px;text-align:center;\">#%s %d</td>", color.AsString(), color.AsString(), color.Cnt))
	}
	buff.WriteString("</tr></table>")
	return buff.String()
}

func processBatch(k int, bitarr []int, img image.Image) string { // WILL BE DELETED
	var buff strings.Builder

	prefix := fmt.Sprintf("K=%d, ", k)
	resizeSize := uint(prominentcolor.DefaultSize)
	bgmasks := prominentcolor.GetDefaultMasks()

	for i := 0; i < len(bitarr); i++ {
		res, err := prominentcolor.KmeansWithAll(k, img, bitarr[i], resizeSize, bgmasks)
		if err != nil {
			log.Println(err)
			continue
		}
		buff.WriteString("<h3>" + prefix + bitInfo(bitarr[i]) + "<h3>")
		buff.WriteString(outputColorRange(res))
	}

	return buff.String()
}

func bitInfo(bits int) string { // WILL BE DELETED
	list := make([]string, 0, 4)
	// random seed or Kmeans++
	if prominentcolor.IsBitSet(bits, prominentcolor.ArgumentSeedRandom) {
		list = append(list, "Random seed")
	} else {
		list = append(list, "Kmeans++")
	}
	// Mean or median
	if prominentcolor.IsBitSet(bits, prominentcolor.ArgumentAverageMean) {
		list = append(list, "Mean")
	} else {
		list = append(list, "Median")
	}
	// LAB or RGB
	if prominentcolor.IsBitSet(bits, prominentcolor.ArgumentLAB) {
		list = append(list, "LAB")
	} else {
		list = append(list, "RGB")
	}
	// Cropping or not
	if prominentcolor.IsBitSet(bits, prominentcolor.ArgumentNoCropping) {
		list = append(list, "No cropping")
	} else {
		list = append(list, "Cropping center")
	}
	// Done
	return strings.Join(list, ", ")
}

func createImageColorSummary(imagePaths []string) { // WILL BE DELETED
	// Prepare
	outputDirectory := "./"

	var buff strings.Builder
	buff.WriteString("<html><body><h1>Colors listed in order of dominance: hex color followed by number of entries</h1><table border=\"1\">")

	for _, file := range imagePaths {
		// Define the differents sets of params
		methods := []int{
			prominentcolor.ArgumentAverageMean | prominentcolor.ArgumentNoCropping,
			prominentcolor.ArgumentNoCropping,
			prominentcolor.ArgumentDefault,
		}

		// Load the image
		img, err := loadImage(file)
		if err != nil {
			log.Printf("Error loading image %s\n", file)
			log.Println(err)
			continue
		}
		// Process & html output
		buff.WriteString("<tr><td><img src=\"" + file + "\" width=\"http.StatusOK\" border=\"1\"></td><td>")
		buff.WriteString(processBatch(3, methods, img))
		buff.WriteString("</td></tr>")
	}

	// Finalize the html output
	buff.WriteString("</table></body><html>")

	// And write it to the disk
	if err := os.WriteFile(outputDirectory+"output.html", []byte(buff.String()), 0644); err != nil {
		panic(err)
	}
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
		createImageColorSummary(imagePaths)
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
