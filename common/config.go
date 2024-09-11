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

package common

import (
	"time"

	"github.com/google/uuid"
)

// Config contains configurable values
// shared between source and destination NATS PubSub connector.
type Config struct {
	// The connection URLs pointed to NATS instances.
	URLs []string `json:"urls" validate:"required"`
	// The name of a subject which the connector should use to read/write records.
	Subject string `json:"subject" validate:"required"`
	// Optional connection name (can come in handy when it comes to monitoring).
	ConnectionName string `json:"connectionName"`
	// A path pointed to a NKey pair.
	NKeyPath string `json:"nkeyPath"`
	// A path pointed to a credentials file.
	CredentialsFilePath string `json:"credentialsFilePath"`
	// Sets the number of reconnect attempts that will be tried before giving up.
	// If negative, it will never give up trying to reconnect.
	MaxReconnects int `json:"maxReconnects" default:"5"`
	// Sets the time to backoff after attempting a reconnect to a server that we
	// were already connected to previously.
	ReconnectWait time.Duration `json:"reconnectWait" default:"5s"`

	TLS TLSConfig `json:"tls"`
}

type TLSConfig struct {
	// A path pointed to a TLS client certificate, must be present if
	// tls.clientPrivateKeyPath field is also present.
	ClientCertPath string `json:"clientCertPath"`
	// A path pointed to a TLS client private key, must be present if
	// tls.clientCertPath field is also present.
	ClientPrivateKeyPath string `json:"clientPrivateKeyPath"`
	// A path pointed to a TLS root certificate, provide it if you want to verify
	// the server's identity.
	RootCACertPath string `json:"rootCACertPath"`
}

func (c *Config) GetConnectionName() string {
	if c.ConnectionName == "" {
		c.ConnectionName = GenerateConnectionName()
	}
	return c.ConnectionName
}

const (
	// DefaultConnectionNamePrefix is the default connection name prefix.
	DefaultConnectionNamePrefix = "conduit-connection-"
)

// GenerateConnectionName generates a random connection name.
// The connection name will be made up of the default connection name and a random UUID.
func GenerateConnectionName() string {
	return DefaultConnectionNamePrefix + uuid.New().String()
}
