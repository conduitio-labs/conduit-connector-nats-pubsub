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

// ConnectionOptions returns connection options based on the provided config.
func (c Config) ConnectionOptions() ([]nats.Option, error) {
	var opts []nats.Option

	if c.ConnectionName != "" {
		opts = append(opts, nats.Name(c.ConnectionName))
	}

	if c.NKeyPath != "" {
		opt, err := nats.NkeyOptionFromSeed(c.NKeyPath)
		if err != nil {
			return nil, fmt.Errorf("load NKey pair: %w", err)
		}

		opts = append(opts, opt)
	}

	if c.CredentialsFilePath != "" {
		opts = append(opts, nats.UserCredentials(c.CredentialsFilePath))
	}

	if c.TLS.ClientCertPath != "" && c.TLS.ClientPrivateKeyPath != "" {
		opts = append(opts, nats.ClientCert(
			c.TLS.ClientCertPath,
			c.TLS.ClientPrivateKeyPath,
		))
	}

	if c.TLS.RootCACertPath != "" {
		opts = append(opts, nats.RootCAs(c.TLS.RootCACertPath))
	}

	opts = append(opts, nats.MaxReconnects(c.MaxReconnects))
	opts = append(opts, nats.ReconnectWait(c.ReconnectWait))

	return opts, nil
}
