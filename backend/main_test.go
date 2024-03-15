package colorboxd

import "testing"

func TestLoadImage(t *testing.T) {
	_, err := LoadImage("https://www.colorhexa.com/ff0000.png")
	if err != nil {
		t.Errorf("Load valid image ff0000.png shouldn't error, had error: %v\n", err)
	}
	_, err = LoadImage("./images/test/invalid.txt")
	if err == nil {
		t.Errorf("Expected error loading invalid.txt; no error occured")
	}
}

func TestGetRedImageInfo(t *testing.T) {
	imagePath := "https://www.colorhexa.com/ff0000.png"
	entry := Entry{ImageInfo: ImageInfo{Path: imagePath}}
	var expectedHue float64 = 0
	expectedHex := "#FF0000"

	redImage, err := LoadImage(imagePath)
	if err != nil {
		t.Errorf("Load valid image 00ff00.png shouldn't error, had error: %v\n", err)
	}

	redImageInfo, err := getImageInfo(entry, redImage)
	if err != nil {
		t.Errorf("Error getting image info: %v\n", err)
		return
	}

	if redImageInfo.ImageInfo.Colors[0].hex != expectedHex {
		t.Errorf("Expected hex %v, computed %v.", expectedHex, redImageInfo.ImageInfo.Colors[0].hex)
	}
	if redImageInfo.ImageInfo.Colors[0].hsl.h != expectedHue {
		t.Errorf("Expected hue %v, computed %v.", expectedHue, redImageInfo.ImageInfo.Colors[0].hsl.h)
	}
}
