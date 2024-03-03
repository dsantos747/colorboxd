package colorboxd

import (
	"image"

	"github.com/lucasb-eyer/go-colorful"
)

// The response format from Letterboxd auth/token endpoint
type AccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	NotBefore    int    `json:"notBefore"`
	Issuer       string `json:"issuer"`
	EncodedToken string `json:"encodedToken"`
}

// The response format from Letterboxd /me endpoint
type Member struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	GivenName   string `json:"givenName"`
	Username    string `json:"username"`
}

// HTTPAuthUser responds to the client with this format
type AuthUserResponse struct {
	Token          string
	TokenType      string
	TokenExpiresIn int
	TokenRefresh   string
	UserId         string
	Username       string
	UserGivenName  string
}

// The response format from Letterboxd /lists endpoint
type ListsResponse struct {
	Cursor string        `json:"cursor"` // Keep this cursor around - need it to handle cases where user has over 100 lists (king)
	Items  []ListSummary `json:"items"`
}

// Summary information about a list. In our case, one of the user's lists
type ListSummary struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     int    `json:"version"`
	FilmCount   int    `json:"filmCount"`
	Description string `json:"description"`
}

// The (partial) response format from Letterboxd list/{id}/entries endpoint
type ListEntriesResponse struct {
	EntryID string `json:"entryId"`
	Film    film   `json:"film"`
}
type film struct {
	Adult              bool     `json:"adult"`
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Poster             coverImg `json:"poster"`
	AdultPoster        coverImg `json:"adultPoster"`
	PosterCustomisable bool     `json:"posterCustomisable"`
	ReleaseYear        int      `json:"releaseYear"`
}
type coverImg struct {
	Sizes []imgSize `json:"sizes"`
}
type imgSize struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
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

// An images path and colour information
type ImageInfo struct {
	Path  string
	Hex   string
	Color colorful.Color
	Hue   float64
}

// Holds a loaded image, alongside all information about that film entry
type Image struct {
	img  image.Image
	info Entry
}

// This is the format of the request body for HTTPWriteList
type WriteListRequest struct {
	AccessToken string          `json:"accessToken"`
	List        ListWithEntries `json:"list"` // This being ListWithEntries (rather than any) is what is causing the error
	Offset      int             `json:"offset"`
}
type ListWithEntries struct {
	ListSummary
	Entries []Entry `json:"entries"`
}

// This is the required format for making a PATCH request to letterboxd list/{id} endpoint.
//
// Note: This struct only includes parameters we are interested in controlling/modifying
type ListUpdateRequest struct {
	Version int               `json:"version"`
	Entries []listUpdateEntry `json:"entries"`
}
type listUpdateEntry struct {
	Action      string `json:"action"`
	Position    int    `json:"position"`
	NewPosition int    `json:"newPosition"`
}

// This is the response format from a PATCH request to the letterboxd list/{id} endpoint.
type ListUpdateResponse struct {
	Messages []ListUpdateMessage `json:"messages"`
}
type ListUpdateMessage struct {
	Type  string `json:"type"`
	Code  string `json:"code"`
	Title string `json:"title"`
}