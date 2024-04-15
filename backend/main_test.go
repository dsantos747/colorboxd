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
var testListEntries *[]Entry

func TestLoadImage(t *testing.T) {
	_, err := loadImage("https://www.colorhexa.com/ff0000.png")
	if err != nil {
		t.Errorf("Load valid image ff0000.png shouldn't error, had error: %v\n", err)
	}
	_, err = loadImage("./.gitignore")
	if err == nil {
		t.Errorf("Expected error loading .gitignore as image; no error occurred")
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
	err := LoadEnv()
	if err != nil {
		t.Errorf("Could not load environment variables from .env file: %v\n", err)
		return
	}

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
	tStart := time.Now()
	for !strings.Contains(redirectURL, "https://colorboxd.com/user?code=") {
		currentURL := page.MustInfo().URL
		if strings.Contains(currentURL, "https://colorboxd.com/user?code=") {
			redirectURL = currentURL
			break
		}
		if time.Since(tStart) > 10*time.Second {
			return nil, fmt.Errorf("Timed out waiting for redirect")
		}
		time.Sleep(250 * time.Millisecond)
	}

	codeIndex := 32 // This is used to remove everything in the url, prior to the code
	authCode := redirectURL[codeIndex:]

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

// Fetch the entries of the test list. If the amount of entries doesn't match
// the expected amount, fail test.
func TestGetListEntries(t *testing.T) {
	var err error
	testListEntries, err = getListEntries(testToken, testListId)
	if err != nil {
		t.Errorf("failed to retrieve entries from list: %v", err)
	}
	if len(*testListEntries) != 44 {
		t.Errorf("retrieved list entries qty didn't match expected entries qty, expected 44, found %v", len(*testListEntries))
	}
}

func TestProcessListImages(t *testing.T) {
	entriesWithImageInfo, err := processListImages(testListEntries)
	if err != nil {
		t.Errorf("failed to process posters for list entries: %v", err)
		return
	}
	for _, entry := range *entriesWithImageInfo {
		if len(entry.ImageInfo.Colors) == 0 {
			t.Errorf("0 colors found for poster for film: %s", entry.Name)
			return
		}
	}
}
