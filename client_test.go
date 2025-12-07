// client_test.go
package smartme_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/rolacher/go-smartme-client"
)

// ptr is a helper function to create a pointer to a value of any type.
func ptr[T any](v T) *T {
	return &v
}

// setup sets up a test HTTP server along with a smartme.Client
// configured to communicate with that server.
func setup(t *testing.T) (*smartme.Client, *http.ServeMux, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	// The client is configured to use the mock server's URL.
	client, err := smartme.NewClient("test-user", "test-pass", smartme.WithBaseURL(server.URL+"/"))
	if err != nil {
		t.Fatalf("smartme.NewClient failed: %v", err)
	}

	// A teardown function is returned to close the server.
	teardown := func() {
		server.Close()
	}

	return client, mux, teardown
}

func TestClient_GetDevices_Success(t *testing.T) {
	client, mux, teardown := setup(t)
	defer teardown()

	// Create the expected response from the server
	mockTimeString := "2025-01-01T12:00:00Z"
	expectedDevices := []smartme.Device{
		{
			Id:             ptr("a1b2c3d4-e5f6-7890-1234-567890abcdef"),
			Name:           ptr("Hauptzähler"),
			Serial:         ptr(int64(12345678)),
			ActivePower:    ptr(1500.5),
			CounterReading: ptr(9876.54),
			ValueDate:      ptr(mockTimeString),
		},
	}

	// Configure the mock server to respond to the API path
	mux.HandleFunc("/api/Devices", func(w http.ResponseWriter, r *http.Request) {
		// Check if the method and authentication are correct
		if r.Method != http.MethodGet {
			t.Errorf("Expected request method GET, got %s", r.Method)
		}
		user, pass, ok := r.BasicAuth()
		if !ok || user != "test-user" || pass != "test-pass" {
			t.Errorf("Basic Auth header is missing or wrong. User: %s, Pass: %s", user, pass)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(expectedDevices)
	})

	// Execute the actual test
	devices, err := client.GetDevices(context.Background())
	if err != nil {
		t.Fatalf("client.GetDevices returned an unexpected error: %v", err)
	}

	// Check if the result matches the expectations
	if !reflect.DeepEqual(devices, expectedDevices) {
		t.Errorf("client.GetDevices returned %+v, want %+v", devices, expectedDevices)
	}
}

func TestClient_GetDevices_ServerError(t *testing.T) {
	client, mux, teardown := setup(t)
	defer teardown()

	// Configure the mock server to return an error
	mux.HandleFunc("/api/Devices", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
	})

	// Run the test and expect an error
	_, err := client.GetDevices(context.Background())
	if err == nil {
		t.Fatal("client.GetDevices should have returned an error, but got nil")
	}

	expectedErrorMsg := "API error: 500 Internal Server Error (status code: 500)"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Error message was '%s', want '%s'", err.Error(), expectedErrorMsg)
	}
}
