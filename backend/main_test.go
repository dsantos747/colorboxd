package colorboxd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/joho/godotenv"
)

// Letterboxd member id: 67W7X
// Letterboxd Test List id: tqtA2

func TestLoadImage(t *testing.T) {
	_, err := loadImage("https://www.colorhexa.com/ff0000.png")
	if err != nil {
		t.Errorf("Load valid image ff0000.png shouldn't error, had error: %v\n", err)
	}
	_, err = loadImage("./images/test/invalid.txt")
	if err == nil {
		t.Errorf("Expected error loading invalid.txt; no error occured")
	}
}

func TestGetRedImageInfo(t *testing.T) {
	imagePath := "https://www.colorhexa.com/ff0000.png"
	entry := Entry{ImageInfo: ImageInfo{Path: imagePath}}
	var expectedHue float64 = 0
	expectedHex := "#FF0000"

	redImage, err := loadImage(imagePath)
	if err != nil {
		t.Errorf("Load valid image ff0000.png shouldn't error, had error: %v\n", err)
	}

	redImageInfo, err := getImageInfo(entry, redImage)
	if err != nil {
		t.Errorf("Error getting image info: %v\n", err)
	}

	if redImageInfo.ImageInfo.Colors[0].hex != expectedHex {
		t.Errorf("Expected hex %v, computed %v.", expectedHex, redImageInfo.ImageInfo.Colors[0].hex)
	}
	if redImageInfo.ImageInfo.Colors[0].h != expectedHue {
		t.Errorf("Expected hue %v, computed %v.", expectedHue, redImageInfo.ImageInfo.Colors[0].h)
	}
}

func TestGetAccessToken(t *testing.T) {
	authCode, err := getAuthCode()
	if err != nil {
		t.Errorf("failed to generate auth code: %v", err)
	}

	accessTokenResponse, err := getAccessToken(*authCode)
	if err != nil {
		t.Errorf("could not create valid access token for provided auth code: %v", err)
	}

	if accessTokenResponse.AccessToken == "" {
		t.Errorf("no valid access token present in response: %v", err)
	}

	member, err := getMemberId(accessTokenResponse.AccessToken)
	if err != nil {
		t.Errorf("could not retrieve member ID: %v", err)
	}

	fmt.Println(member.DisplayName)
}

// Helper function to log in to letterboxd account and return an auth code.
// This works because I have already authorised colorboxd for my account.
func getAuthCode() (*string, error) {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
		return nil, err
	}

	launcher := launcher.New()
	launcher.Leakless(false)
	u := launcher.MustLaunch()

	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	authCodeChan := make(chan string)

	go func() {
		_ = browser.EachEvent(func(e *proto.NetworkResponseReceived) {
			if strings.Contains(e.Response.URL, "https://colorboxd.com/user?code=") {
				// Parse the URL to extract the auth code
				parts := strings.Split(e.Response.URL, "=")
				if len(parts) > 1 {
					authCode := parts[1]
					fmt.Println("Received auth code:", authCode)
					authCodeChan <- authCode
				}
			}
		})
	}()

	// Create a new page
	page := browser.MustPage(os.Getenv("LBOXD_AUTH_URL")).MustWaitStable()

	// Enter login details
	page.MustElement("#field-username").MustInput(os.Getenv("MY_LBOXD_USER"))
	page.MustElement("#field-password").MustInput(os.Getenv("MY_LBOXD_PASS")).MustType(input.Enter)

	page.MustWaitNavigation()

	authCode := <-authCodeChan

	browser.MustClose()
	launcher.Kill()

	return &authCode, nil
}

// getUserLists test
// To do this, first need to generate an access token
