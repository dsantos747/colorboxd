package colorboxd

import (
	"image"

	"github.com/lucasb-eyer/go-colorful"
)

type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	NotBefore    int    `json:"notBefore"`
	Issuer       string `json:"issuer"`
	EncodedToken string `json:"encodedToken"`
}

// Letterboxd Member
type Member struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	GivenName   string `json:"givenName"`
	Username    string `json:"username"`
}

// Letterboxd Lists
type ListSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     int    `json:"version"`
	FilmCount   int    `json:"filmCount"`
	Description string `json:"description"`
}
type ListsResponse struct {
	Cursor string        `json:"cursor"`
	Items  []ListSummary `json:"items"`
}

type AuthUserResponse struct {
	Token          string
	TokenType      string
	TokenExpiresIn int
	UserId         string
	Username       string
	UserGivenName  string
}

// Image information
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
