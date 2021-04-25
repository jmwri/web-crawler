package internal

import (
	"errors"
	"io"
	"net/http"
	"time"
)

// LoaderFunc loads requested page
type LoaderFunc func(p string) (io.ReadCloser, error)

// NewHTTPGetLoader returns a new HTTPGetLoader
func NewHTTPGetLoader(client *http.Client) HTTPGetLoader {
	return HTTPGetLoader{
		client: client,
	}
}

// HTTPGetLoader loads pages with an HTTP GET request
type HTTPGetLoader struct {
	client *http.Client
}

// Load the requested page with an HTTP GET request
func (l *HTTPGetLoader) Load(p string) (io.ReadCloser, error) {
	res, err := l.client.Get(p)
	if err != nil {
		return nil, err
	}
	// TODO: Check content type
	if res.StatusCode < 200 || res.StatusCode > 299 {
		return res.Body, errors.New("failed to load page")
	}
	return res.Body, err
}

// LoaderWithRetry wraps a LoaderFunc with a retry
func LoaderWithRetry(l LoaderFunc, b BackoffFunc, attempts int) LoaderFunc {
	return func(p string) (io.ReadCloser, error) {
		var attempt int
		for {
			attempt += 1
			time.Sleep(b(attempt))
			res, err := l(p)
			if err != nil {
				if attempt >= attempts {
					return nil, err
				}
				continue
			}
			return res, err
		}
	}
}
