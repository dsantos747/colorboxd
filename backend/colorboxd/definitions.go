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
	TokenRefresh   string
	UserId         string
	Username       string
	UserGivenName  string
}

// This is the output format of GetListEntries
type Entry struct {
	EntryID            string `json:"entryId"`
	FilmID             string `json:"filmId"`
	Name               string `json:"name"`
	ReleaseYear        int    `json:"releaseYear"`
	Adult              bool   `json:"adult"`
	PosterCustomisable bool   `json:"posterCustomisable"`
	PosterURL          string `json:"posterUrl"`
	AdultPosterURL     string `json:"adultPosterUrl"`
	ImageInfo          ImageInfo
}

// //// These types are associated with the list/{id}/entries endpoint response
type GetListEntriesResponse struct {
	Items []struct {
		EntryID string `json:"entryId"`
		Film    Film   `json:"film"`
	} `json:"items"`
}

type ListEntriesResponse struct {
	EntryID string `json:"entryId"`
	Film    Film   `json:"film"`
}

type Film struct {
	Adult              bool     `json:"adult"`
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Poster             CoverImg `json:"poster"`
	AdultPoster        CoverImg `json:"adultPoster"`
	PosterCustomisable bool     `json:"posterCustomisable"`
	ReleaseYear        int      `json:"releaseYear"`
}

type CoverImg struct {
	Sizes []Size `json:"sizes"`
}

type Size struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

// ////

type SortListResponse struct {
	Entries  any
	SortList any // Can do better than "any"
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
	info Entry
}
