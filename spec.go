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

package nats

import (
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/config"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/source"
	sdk "github.com/conduitio/conduit-connector-sdk"
)

type Spec struct{}

// Specification returns the Plugin's Specification.
func Specification() sdk.Specification {
	return sdk.Specification{
		Name:    "nats",
		Summary: "A NATS PubSub source and destination plugin for Conduit, written in Go.",
		Description: "The NATS PubSub connector is one of Conduit plugins. " +
			"It provides both, a source and a destination NATS PubSub connector.",
		Version: "v0.1.0",
		Author:  "Meroxa, Inc.",
		SourceParams: map[string]sdk.Parameter{
			config.KeyURLs: {
				Default:     "",
				Required:    true,
				Description: "The connection URLs pointed to NATS instances.",
			},
			config.KeySubject: {
				Default:     "",
				Required:    true,
				Description: "A name of a subject from which or to which the connector should read/write.",
			},
			config.KeyConnectionName: {
				Default:     "conduit-connection-<uuid>",
				Required:    false,
				Description: "Optional connection name which will come in handy when it comes to monitoring.",
			},
			config.KeyNKeyPath: {
				Default:     "",
				Required:    false,
				Description: "A path pointed to a NKey pair.",
			},
			config.KeyCredentialsFilePath: {
				Default:     "",
				Required:    false,
				Description: "A path pointed to a credentials file.",
			},
			config.KeyTLSClientCertPath: {
				Default:  "",
				Required: false,
				//nolint:lll // long description
				Description: "A path pointed to a TLS client certificate, must be present if tls.clientPrivateKeyPath field is also present.",
			},
			config.KeyTLSClientPrivateKeyPath: {
				Default:  "",
				Required: false,
				//nolint:lll // long description
				Description: "A path pointed to a TLS client private key, must be present if tls.clientCertPath field is also present.",
			},
			config.KeyTLSRootCACertPath: {
				Default:     "",
				Required:    false,
				Description: "A path pointed to a TLS root certificate, provide if you want to verify server’s identity.",
			},
			config.KeyMaxReconnects: {
				Default:  "5",
				Required: false,
				Description: "Sets the number of reconnect attempts " +
					"that will be tried before giving up. If negative, " +
					"then it will never give up trying to reconnect.",
			},
			config.KeyReconnectWait: {
				Default:  "5s",
				Required: false,
				Description: "Sets the time to backoff after attempting a reconnect " +
					"to a server that we were already connected to previously.",
			},
			source.ConfigKeyBufferSize: {
				Default:     "1024",
				Required:    false,
				Description: "A buffer size for consumed messages.",
			},
		},
		DestinationParams: map[string]sdk.Parameter{
			config.KeyURLs: {
				Default:     "",
				Required:    true,
				Description: "The connection URLs pointed to NATS instances.",
			},
			config.KeySubject: {
				Default:     "",
				Required:    true,
				Description: "A name of a subject from which or to which the connector should read/write.",
			},
			config.KeyConnectionName: {
				Default:     "conduit-connection-<uuid>",
				Required:    false,
				Description: "Optional connection name which will come in handy when it comes to monitoring.",
			},
			config.KeyNKeyPath: {
				Default:     "",
				Required:    false,
				Description: "A path pointed to a NKey pair.",
			},
			config.KeyCredentialsFilePath: {
				Default:     "",
				Required:    false,
				Description: "A path pointed to a credentials file.",
			},
			config.KeyTLSClientCertPath: {
				Default:  "",
				Required: false,
				//nolint:lll // long description
				Description: "A path pointed to a TLS client certificate, must be present if tls.clientPrivateKeyPath field is also present.",
			},
			config.KeyTLSClientPrivateKeyPath: {
				Default:  "",
				Required: false,
				//nolint:lll // long description
				Description: "A path pointed to a TLS client private key, must be present if tls.clientCertPath field is also present.",
			},
			config.KeyTLSRootCACertPath: {
				Default:     "",
				Required:    false,
				Description: "A path pointed to a TLS root certificate, provide if you want to verify server’s identity.",
			},
			config.KeyMaxReconnects: {
				Default:  "5",
				Required: false,
				Description: "Sets the number of reconnect attempts " +
					"that will be tried before giving up. If negative, " +
					"then it will never give up trying to reconnect.",
			},
			config.KeyReconnectWait: {
				Default:  "5s",
				Required: false,
				Description: "Sets the time to backoff after attempting a reconnect " +
					"to a server that we were already connected to previously.",
			},
		},
	}
}
