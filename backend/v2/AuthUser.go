package colorboxd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

// AuthUser is the handler for autorisation of the users letterboxd account to the colorboxd app.
func AuthUser(w http.ResponseWriter, r *http.Request) {
	var err error

	// // Read env variables
	// err = LoadEnv()
	// if err != nil {
	// 	fmt.Printf("Could not load environment variables from .env file: %v\n", err)
	// 	return
	// }

	// // Set necessary headers for CORS and cache policy
	// w.Header().Set("Access-Control-Allow-Origin", os.Getenv("BASE_URL"))
	// w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Cache-Control", "private, max-age=3595") // Expire time of token (-5s for safety)

	// Read authCode from query url - return error if not present
	authCode := r.URL.Query().Get("authCode")
	if authCode == "" {
		ReturnError(w, "Missing or empty 'authCode' query parameter", http.StatusBadRequest)
		return
	}

	// Get Access Token
	accessTokenResponse, err := getAccessToken(authCode)
	if err != nil {
		ReturnError(w, fmt.Errorf("could not create valid access token: %w", err).Error(), http.StatusInternalServerError)
		return
	}
	if accessTokenResponse.AccessToken == "" {
		ReturnError(w, fmt.Errorf("could not generate access token from provided auth code: %w", err).Error(), http.StatusBadRequest)
		return
	}

	member, err := getMemberId(accessTokenResponse.AccessToken)
	if err != nil {
		ReturnError(w, fmt.Errorf("could not retrieve member ID: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	response := AuthUserResponse{
		Token:          accessTokenResponse.AccessToken,
		TokenType:      accessTokenResponse.TokenType,
		TokenRefresh:   accessTokenResponse.RefreshToken,
		TokenExpiresIn: accessTokenResponse.ExpiresIn,
		UserId:         member.ID,
		Username:       member.Username,
		UserGivenName:  member.GivenName,
	}

	// Return response to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getAccessToken(authCode string) (*AccessTokenResponse, error) {
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

func getMemberId(token string) (*Member, error) {
	method := "GET"
	endpoint := fmt.Sprintf("%s/me", os.Getenv("LBOXD_BASEURL"))
	headers := map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)}

	response, err := MakeHTTPRequest(method, endpoint, nil, headers)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var responseData map[string]json.RawMessage
	if err = json.NewDecoder(response.Body).Decode(&responseData); err != nil {
		return nil, err
	}
	var member Member
	if err = json.Unmarshal(responseData["member"], &member); err != nil {
		return nil, err
	}

	return &member, nil
}
