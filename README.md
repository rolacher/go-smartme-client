# Go Client for the smart-me API

[![Go Report Card](https://goreportcard.com/badge/github.com/rolacher/go-smartme-client)](https://goreportcard.com/report/github.com/rolacher/go-smartme-client)
[![Go Reference](https://pkg.go.dev/badge/github.com/rolacher/go-smartme-client.svg)](https://pkg.go.dev/github.com/rolacher/go-smartme-client)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`go-smartme-client` is a Go client library for the smart-me.com API. It provides a convenient way to interact with your smart-me devices, retrieve current values, and access historical data.

## Features

*   Typed Go structs for all major API resources (`Device`, `DeviceValues`, etc.).
*   Full support for `context.Context` for request cancellation and deadlines.
*   Clean, idiomatic Go API design.
*   Configurable HTTP client for custom timeouts or transport layers.
*   Includes unit tests with mocks and optional integration tests against the live API.

## Installation

To add the library to your project, run:

```sh
go get github.com/rolacher/go-smartme-client
```

## Usage

Here is a basic example of how to create a client and retrieve a list of your devices.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/rolacher/go-smartme-client"
)

func main() {
	// It's recommended to use environment variables for credentials.
	username := os.Getenv("SMARTME_USERNAME")
	password := os.Getenv("SMARTME_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("SMARTME_USERNAME and SMARTME_PASSWORD environment variables must be set")
	}

	// 1. Create a new client
	client, err := smartme.NewClient(username, password)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// 2. Get all devices
	devices, err := client.GetDevices(context.Background())
	if err != nil {
		log.Fatalf("Failed to get devices: %v", err)
	}

	fmt.Printf("Found %d devices.\n", len(devices))

	// 3. Print details of the first device
	if len(devices) > 0 {
		firstDevice := devices
		if firstDevice.Name != nil && firstDevice.Id != nil {
			fmt.Printf("First device: %s (ID: %s)\n", *firstDevice.Name, *firstDevice.Id)

			// 4. Get current values for the first device
			deviceValues, err := client.GetValues(context.Background(), *firstDevice.Id)
			if err != nil {
				log.Fatalf("Failed to get values for device %s: %v", *firstDevice.Name, err)
			}

			fmt.Printf("Current values for %s (at %s):\n", *firstDevice.Name, deviceValues.Date.Local())
			for _, v := range deviceValues.Values {
				// Example: Print only the active power
				if v.Obis == "1-0:1.7.0*255" {
					fmt.Printf("  - Active Power: %.2f W\n", v.Value)
				}
			}
		}
	}
}
```

## Testing

The library includes both unit and integration tests.

### Unit Tests

The unit tests use a mock server and do not require a network connection.

```sh
go test -v .
```

### Integration Tests

The integration tests run against the live smart-me API and require credentials. Create a file `~/.smartme-client-config.json` with your username and password:

```json
{
  "username": "YOUR_USERNAME",
  "password": "YOUR_PASSWORD"
}
```

Then, run the tests using the `integration` build tag:

```sh
go test -v -tags=integration .
```

## License

This project is licensed under the MIT License. See the LICENSE file for details.