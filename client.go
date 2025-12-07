// client.go
package smartme

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL = "https://api.smart-me.com/"
	defaultTimeout = 10 * time.Second
)

// Client is the API client for the smart-me API.
type Client struct {
	httpClient *http.Client
	baseURL    *url.URL
	username   string
	password   string
}

// NewClient creates a new instance of the smart-me API client.
// Authentication is done via Basic Auth with username and password.
func NewClient(username, password string, opts ...Option) (*Client, error) {
	if username == "" {
		return nil, fmt.Errorf("username must not be empty")
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		baseURL:  baseURL,
		username: username,
		password: password,
	}

	// Apply functional options
	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: defaultTimeout,
		}
	}

	return c, nil
}

// newRequest creates a new HTTP request with the necessary headers.
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	rel, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	// Create the full URL
	fullURL := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, method, fullURL.String(), body)
	if err != nil {
		return nil, err
	}

	// Set Basic Authentication
	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// do executes the request and decodes the response into the provided struct.
func (c *Client) do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		// Catch context errors (e.g., timeout)
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		// Implement more robust error handling here
		return resp, fmt.Errorf("API error: %s (status code: %d)", resp.Status, resp.StatusCode)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return resp, fmt.Errorf("error decoding response: %w", err)
		}
	}

	return resp, nil
}

// GetDevices retrieves the list of all devices.
// Corresponds to the API call: GET /api/Devices
func (c *Client) GetDevices(ctx context.Context) ([]Device, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "api/Devices", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var devices []Device
	_, err = c.do(req, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// GetValues retrieves the last values of a specific device.
// Corresponds to the API call: GET /api/Values/{id}
func (c *Client) GetValues(ctx context.Context, deviceID string) (*DeviceValues, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("deviceID must not be empty")
	}

	path := fmt.Sprintf("api/Values/%s", deviceID)
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var deviceValues DeviceValues
	_, err = c.do(req, &deviceValues)
	if err != nil {
		return nil, err
	}

	return &deviceValues, nil
}

// GetValuesInPast retrieves the first value found before a given date for a specific device.
// Corresponds to the API call: GET /api/ValuesInPast/{id}?date={date}
func (c *Client) GetValuesInPast(ctx context.Context, deviceID string, date time.Time) (*Value, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("deviceID must not be empty")
	}

	path := fmt.Sprintf("api/ValuesInPast/%s?date=%s", deviceID, date.Format(time.RFC3339))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var value Value
	_, err = c.do(req, &value)
	if err != nil {
		return nil, err
	}

	return &value, nil
}

// GetValuesInPastMultiple retrieves multiple values of a device within a given time range.
// Note: This call might require a professional license for the smart-me account.
// Corresponds to the API call: GET /api/ValuesInPastMultiple/{id}?startDate={startDate}&endDate={endDate}
func (c *Client) GetValuesInPastMultiple(ctx context.Context, deviceID string, startDate, endDate time.Time) ([]Value, error) {
	if deviceID == "" {
		return nil, fmt.Errorf("deviceID must not be empty")
	}

	path := fmt.Sprintf("api/ValuesInPastMultiple/%s?startDate=%s&endDate=%s", deviceID, startDate.Format(time.RFC3339), endDate.Format(time.RFC3339))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	var values []Value
	_, err = c.do(req, &values)
	if err != nil {
		return nil, err
	}

	return values, nil
}
