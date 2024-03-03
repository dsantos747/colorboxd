package colorboxd

import (
	"io"
	"net/http"
)

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
	// defer response.Body.Close()

	return response, nil
}
