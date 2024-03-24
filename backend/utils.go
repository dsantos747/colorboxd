package colorboxd

import (
	"fmt"
	"io"
	"net/http"
)

func ReturnError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, message, statusCode)
}

// Makes an HTTP request of the required method to the specified endpoint.
// If response code is >= 400 , returns an error with response.Status
func MakeHTTPRequest(method, endpoint string, body io.Reader, headers map[string]string) (*http.Response, error) {
	// Prepare the request
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Perform the request
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		err = fmt.Errorf("%s", response.Status)
		return nil, err
	}

	return response, nil
}
