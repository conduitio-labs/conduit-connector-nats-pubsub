// Copyright © 2022 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package test

import (
	"fmt"

	"github.com/nats-io/nats.go"
)

// TestURL is a URL of a test NATS server.
var (
	TestURL             = "nats://127.0.0.1:4222"
	TestURLWithPassword = "nats://admin:supersecret@127.0.0.1:4223" //nolint:gosec // test creds
	TestURLWithToken    = "nats://s3cr3t@127.0.0.1:4224"            //nolint:gosec // test creds
	TestURLWithNKey     = "nats://127.0.0.1:4225"
)

// GetTestConnection returns a connection to a test NATS server.
func GetTestConnection(url string) (*nats.Conn, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, fmt.Errorf("connect to NATS server: %w", err)
	}

	return conn, nil
}
