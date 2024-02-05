package main

import (
	"cmp"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	prominentcolor "github.com/EdlinOrg/prominentcolor"
	"github.com/lucasb-eyer/go-colorful"
)

type ImageColorInfo struct {
	Path  string
	Hex   string
	Color colorful.Color
	Hue   float64
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
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

func outputColorRange(colorRange []prominentcolor.ColorItem) string {
	var buff strings.Builder
	buff.WriteString("<table><tr>")
	for _, color := range colorRange {
		buff.WriteString(fmt.Sprintf("<td style=\"background-color: #%s;width:200px;height:50px;text-align:center;\">#%s %d</td>", color.AsString(), color.AsString(), color.Cnt))
	}
	buff.WriteString("</tr></table>")
	return buff.String()
}

func outputTitle(str string) string {
	return "<h3>" + str + "</h3>"
}

func processBatch(k int, bitarr []int, img image.Image) string {
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
		buff.WriteString(outputTitle(prefix + bitInfo(bitarr[i])))
		buff.WriteString(outputColorRange(res))
	}

	return buff.String()
}

func bitInfo(bits int) string {
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

func getDominantImage(k int, method int, img image.Image) string {
	resizeSize := uint(prominentcolor.DefaultSize)
	bgmasks := prominentcolor.GetDefaultMasks()

	res, err := prominentcolor.KmeansWithAll(k, img, method, resizeSize, bgmasks)
	if err != nil {
		log.Println(err)
	}

	return res[0].AsString()
}

func createImageColorSummary(imagePaths []string) {
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
		buff.WriteString("<tr><td><img src=\"" + file + "\" width=\"200\" border=\"1\"></td><td>")
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

	imagePaths, err := getImagePaths("./images/example")
	if err != nil {
		fmt.Printf("Error reading images from directory: %v\n", err)
		return
	}

	createImageColorSummary(imagePaths)

	var buff strings.Builder

	var imageSlice []ImageColorInfo

	for _, file := range imagePaths {
		var method int = prominentcolor.ArgumentNoCropping
		// Load the image
		img, err := loadImage(file)
		if err != nil {
			log.Printf("Error loading image %s\n", file)
			log.Println(err)
			continue
		}

		hex := "#" + getDominantImage(3, method, img)
		color, _ := colorful.Hex(hex)
		hue, _, _ := color.Hsv()

		imageSlice = append(imageSlice, ImageColorInfo{file, hex, color, hue})

	}

	slices.SortFunc[[]ImageColorInfo](imageSlice, func(a, b ImageColorInfo) int { return cmp.Compare[float64](a.Hue, b.Hue) })

	buff.WriteString("<html><body><h1>Ordered Images</h1><table border=\"1\">")

	for _, img := range imageSlice {
		fmt.Println(img.Hue, img.Path)
		buff.WriteString("<tr><td><img src=\"" + img.Path + "\" width=\"200\" border=\"1\"></td>")
		buff.WriteString(fmt.Sprintf("<td style=\"background-color: %s;width:200px;height:50px;text-align:center;\">Color: %s</td></tr>", img.Hex, img.Hex))

	}

	buff.WriteString("</table></body><html>")
	if err = os.WriteFile("./sortedImages.html", []byte(buff.String()), 0644); err != nil {
		panic(err)
	}

}
