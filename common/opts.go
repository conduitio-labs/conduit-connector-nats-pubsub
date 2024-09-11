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
	"fmt"

	"github.com/nats-io/nats.go"
)

// GetConnectionOptions returns connection options based on the provided config.
func GetConnectionOptions(config Config) ([]nats.Option, error) {
	var opts []nats.Option

	if config.ConnectionName != "" {
		opts = append(opts, nats.Name(config.ConnectionName))
	}

	if config.NKeyPath != "" {
		opt, err := nats.NkeyOptionFromSeed(config.NKeyPath)
		if err != nil {
			return nil, fmt.Errorf("load NKey pair: %w", err)
		}

		opts = append(opts, opt)
	}

	if config.CredentialsFilePath != "" {
		opts = append(opts, nats.UserCredentials(config.CredentialsFilePath))
	}

	if config.TLS.ClientCertPath != "" && config.TLS.ClientPrivateKeyPath != "" {
		opts = append(opts, nats.ClientCert(
			config.TLS.ClientCertPath,
			config.TLS.ClientPrivateKeyPath,
		))
	}

	if config.TLS.RootCACertPath != "" {
		opts = append(opts, nats.RootCAs(config.TLS.RootCACertPath))
	}

	opts = append(opts, nats.MaxReconnects(config.MaxReconnects))
	opts = append(opts, nats.ReconnectWait(config.ReconnectWait))

	return opts, nil
}
