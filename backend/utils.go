package colorboxd

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if os.Getenv("ENVIRONMENT") == "" {
		err := godotenv.Load()
		if err != nil {
			fmt.Printf("Could not load environment variables from .env file: %v\n", err)
			return err
		}
	}
	return nil
}

func ReturnError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, message, statusCode)
}

var HTTPclient = &http.Client{
	Timeout: time.Second * 30,
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
	response, err := HTTPclient.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		err = fmt.Errorf("%s", response.Status)
		return nil, err
	}

	return response, nil
}
