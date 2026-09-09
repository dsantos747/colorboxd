package colorboxd

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/ratelimit"
)

// LoadEnv attempts to load an env var "ENVIRONMENT". If successful, no further action.
// If not successful, load all envs with godotenv instead
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

// ReturnError sends a http error back to the ResponseWriter w
func ReturnError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "text/plain")
	http.Error(w, message, statusCode)
}

var HTTPclient = &http.Client{
	Timeout: time.Second * 30,
}

// lboxdLimiter throttles all outgoing requests to the Letterboxd API so that a high
// redis cache-miss rate (e.g. when processing a list of posters) doesn't hammer it.
// Rate is requests/sec and can be tuned via the LBOXD_RATE_LIMIT env var.
var lboxdLimiter = ratelimit.New(lboxdRateLimit())

func lboxdRateLimit() int {
	if v := os.Getenv("LBOXD_RATE_LIMIT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return 200
}

const (
	maxRetries     = 3
	retryBaseDelay = 200 * time.Millisecond
	retryMaxDelay  = 5 * time.Second
)

// Makes an HTTP request of the required method to the specified endpoint. Includes rate-limiting and retry/backoff.
func MakeHTTPRequest(method, endpoint string, body io.Reader, headers map[string]string) (*http.Response, error) {
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, err
		}
	}

	var lastErr error
	var retryAfter time.Duration
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			delay := retryAfter
			if delay <= 0 {
				delay = retryBackoff(attempt)
			}
			time.Sleep(delay)
			retryAfter = 0
		}

		lboxdLimiter.Take()

		var reqBody io.Reader
		if bodyBytes != nil {
			reqBody = bytes.NewReader(bodyBytes)
		}

		req, err := http.NewRequest(method, endpoint, reqBody)
		if err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}

		response, err := HTTPclient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		if response.StatusCode == http.StatusTooManyRequests || response.StatusCode >= 500 {
			lastErr = fmt.Errorf("%s", response.Status)
			retryAfter = parseRetryAfter(response.Header.Get("Retry-After"))
			response.Body.Close()
			continue
		}

		if response.StatusCode >= 400 {
			err = fmt.Errorf("%s", response.Status)
			response.Body.Close()
			return nil, err
		}

		return response, nil
	}

	return nil, fmt.Errorf("request to %s failed after %d attempts: %w", endpoint, maxRetries+1, lastErr)
}

// retryBackoff returns an exponential backoff delay (with jitter) for the given attempt number.
func retryBackoff(attempt int) time.Duration {
	delay := min(retryBaseDelay*time.Duration(int64(1)<<uint(attempt-1)), retryMaxDelay)

	return delay/2 + time.Duration(rand.Int63n(int64(delay)/2+1))
}

// parseRetryAfter parses a Retry-After header (either delay-seconds or HTTP-date form).
func parseRetryAfter(header string) time.Duration {
	if header == "" {
		return 0
	}
	if secs, err := strconv.Atoi(header); err == nil {
		return time.Duration(secs) * time.Second
	}
	if t, err := http.ParseTime(header); err == nil {
		return time.Until(t)
	}
	return 0
}
