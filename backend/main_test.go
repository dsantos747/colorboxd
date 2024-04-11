package colorboxd

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/input"
	"github.com/go-rod/rod/lib/launcher"
)

var testToken string
var testUserId string = "67W7X"
var testListId string = "tqtA2"

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

func TestGetAccessTokenAndUser(t *testing.T) {
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

	if member.ID != testUserId {
		t.Errorf("Somehow retrieved the incorrect member id: %v", member.ID)
	}

	testToken = accessTokenResponse.AccessToken
}

// Helper function to log in to letterboxd account and return an auth code.
// This works because I have already authorised colorboxd for my account.
func getAuthCode() (*string, error) {
	err := LoadEnv()
	if err != nil {
		fmt.Printf("Could not load environment variables from .env file: %v\n", err)
		return nil, err
	}

	fmt.Println(os.Getenv("LBOXD_AUTH_URL"))

	// Setup Browser
	launcher := launcher.New()
	launcher.Leakless(false) // Required since antivirus prevents test run if leakless binary is present
	u := launcher.MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	defer browser.MustClose()

	// Navigate and login
	page := browser.MustPage(os.Getenv("LBOXD_AUTH_URL")).MustWaitStable()
	page.MustElement("#field-username").MustInput(os.Getenv("MY_LBOXD_USER"))
	page.MustElement("#field-password").MustInput(os.Getenv("MY_LBOXD_PASS")).MustType(input.Enter)

	// Wait for the redirect
	redirectURL := ""
	for !strings.Contains(redirectURL, "https://colorboxd.com/user?code=") {
		// Get the current URL
		currentURL := page.MustInfo().URL
		if strings.Contains(currentURL, "https://colorboxd.com/user?code=") {
			redirectURL = currentURL
			break
		}
		time.Sleep(250 * time.Millisecond)
	}

	codeIndex := 32 // This is used to remove everything in the url, prior to the code
	authCode := redirectURL[codeIndex:]

	fmt.Println(authCode)

	browser.MustClose()
	launcher.Kill()

	return &authCode, nil
}

func TestGetLists(t *testing.T) {
	userLists, err := getUserLists(testToken, testUserId)
	if err != nil {
		t.Errorf("could not retrieve lists from Letterboxd API: %v", err)
	}

	for _, l := range *userLists {
		if l.ID == testListId {
			return
		}
	}
	t.Errorf("could not find test list ID in returned lists; expected LID: %v", testListId)
}

// getUserLists test
// To do this, first need to generate an access token
