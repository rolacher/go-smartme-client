// options.go
package smartme

import (
	"net/http"
	"net/url"
	"time"
)

// Option is a functional option for configuring the client.
type Option func(*Client)

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets a custom base URL (useful for testing).
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if u, err := url.Parse(baseURL); err == nil {
			c.baseURL = u
		}
	}
}

// WithTimeout sets a custom timeout for the HTTP client.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
	}
}
