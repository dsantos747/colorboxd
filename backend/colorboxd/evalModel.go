package colorboxd

import (
	"fmt"
	"image"
	"log"
	"os"
	"strings"

	prominentcolor "github.com/EdlinOrg/prominentcolor"
)

func CreateImageColorSummary(imagePaths []string) {
	outputDirectory := "./"

	var buff strings.Builder
	buff.WriteString("<html><body><h1>Colors listed in order of dominance: hex color followed by number of entries</h1><table border=\"1\">")

	for _, file := range imagePaths {
		methods := []int{
			prominentcolor.ArgumentAverageMean | prominentcolor.ArgumentNoCropping,
			prominentcolor.ArgumentNoCropping,
			prominentcolor.ArgumentDefault,
		}

		img, err := LoadImage(file)
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
	if err := os.WriteFile(outputDirectory+"output.html", []byte(buff.String()), 0644); err != nil {
		panic(err)
	}
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
		buff.WriteString("<h3>" + prefix + bitInfo(bitarr[i]) + "<h3>")
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
