package colorboxd

import "testing"

func TestLoadImage(t *testing.T) {
	imageSource = &ImgSources.local

	_, err := LoadImage("../images/test/00ff00.png")
	if err != nil {
		t.Errorf("Load valid image 00ff00.png shouldn't error, had error: %v\n", err)
	}
	_, err = LoadImage("./images/test/invalid.txt")
	if err == nil {
		t.Errorf("Expected error loading invalid.txt; no error occured")
	}
}

func TestGetRedImageInfo(t *testing.T) {
	imageSource = &ImgSources.local
	imagePath := "../images/test/ff0000.png"
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

	if redImageInfo.ImageInfo.Hex != expectedHex {
		t.Errorf("Expected hex %v, computed %v.", expectedHex, redImageInfo.ImageInfo.Hex)
	}
	if redImageInfo.ImageInfo.Hue != expectedHue {
		t.Errorf("Expected hue %v, computed %v.", expectedHue, redImageInfo.ImageInfo.Hue)
	}
}
