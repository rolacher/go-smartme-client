//go:build integration

// integration_test.go
package smartme_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/rolacher/go-smartme-client"
)

// testConfig holds the credentials for the integration tests.
type testConfig struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const configFileName = ".smartme-client-config.json"

var config testConfig

// init loads the configuration from the user's home directory.
func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		return // Cannot find home directory, tests will be skipped.
	}

	configPath := filepath.Join(home, configFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return // Config file not found, tests will be skipped.
	}
	_ = json.Unmarshal(data, &config)
}

// setupIntegrationTest creates a real client for integration tests.
// It skips the test if credentials are not available.
func setupIntegrationTest(t *testing.T) *smartme.Client {
	if config.Username == "" || config.Password == "" {
		t.Skipf("Skipping integration test: credentials not found in ~/%s", configFileName)
	}

	client, err := smartme.NewClient(config.Username, config.Password)
	if err != nil {
		t.Fatalf("Failed to create client for integration test: %v", err)
	}

	return client
}

// TestIntegration_GetDevices performs a real API call to get devices.
func TestIntegration_GetDevices(t *testing.T) {
	client := setupIntegrationTest(t)

	devices, err := client.GetDevices(context.Background())
	if err != nil {
		t.Fatalf("client.GetDevices() returned an error: %v", err)
	}

	if devices == nil {
		t.Fatal("client.GetDevices() returned a nil slice, expected non-nil")
	}

	// This is a successful test if we get here without errors.
	// We can add a log for more info.
	t.Logf("Successfully retrieved %d devices from the API.", len(devices))

	if len(devices) > 0 {
		firstDevice := devices[0]
		var deviceName, deviceID string
		if firstDevice.Name != nil {
			deviceName = *firstDevice.Name
		}
		if firstDevice.Id != nil {
			deviceID = *firstDevice.Id
		}
		t.Logf("-> First device found: Name='%s', ID='%s'", deviceName, deviceID)
	}
}
