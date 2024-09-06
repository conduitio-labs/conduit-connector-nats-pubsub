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

package destination

import (
	"context"
	"fmt"
	"strings"

	"github.com/conduitio-labs/conduit-connector-nats-pubsub/common"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/destination/pubsub"
	"github.com/conduitio/conduit-commons/config"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/nats-io/nats.go"
)

// Writer defines a writer interface needed for the Destination.
type Writer interface {
	Write(record opencdc.Record) error
	Close() error
}

// Destination NATS Connector sends records to a NATS subject.
type Destination struct {
	sdk.UnimplementedDestination

	config Config
	writer Writer
}

type Config struct {
	common.Config
}

// NewDestination creates new instance of the Destination.
func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

// Parameters returns a map of named config.Parameters that describe how to configure the Destination.
func (d *Destination) Parameters() config.Parameters {
	return map[string]config.Parameter{
		common.KeyURLs: {
			Default:     "",
			Description: "The connection URLs pointed to NATS instances.",
			Validations: []config.Validation{config.ValidationRequired{}},
		},
		common.KeySubject: {
			Default:     "",
			Description: "A name of a subject to which the connector should write.",
			Validations: []config.Validation{config.ValidationRequired{}},
		},
		common.KeyConnectionName: {
			Default:     "conduit-connection-<uuid>",
			Description: "Optional connection name which will come in handy when it comes to monitoring.",
		},
		common.KeyNKeyPath: {
			Default:     "",
			Description: "A path pointed to a NKey pair.",
		},
		common.KeyCredentialsFilePath: {
			Default:     "",
			Description: "A path pointed to a credentials file.",
		},
		common.KeyTLSClientCertPath: {
			Default: "",
			Description: "A path pointed to a TLS client certificate, must be present " +
				"if tls.clientPrivateKeyPath field is also present.",
		},
		common.KeyTLSClientPrivateKeyPath: {
			Default: "",
			Description: "A path pointed to a TLS client private key, must be present " +
				"if tls.clientCertPath field is also present.",
		},
		common.KeyTLSRootCACertPath: {
			Default:     "",
			Description: "A path pointed to a TLS root certificate, provide if you want to verify server’s identity.",
		},
		common.KeyMaxReconnects: {
			Default: "5",
			Description: "Sets the number of reconnect attempts " +
				"that will be tried before giving up. If negative, " +
				"then it will never give up trying to reconnect.",
		},
		common.KeyReconnectWait: {
			Default: "5s",
			Description: "Sets the time to backoff after attempting a reconnect " +
				"to a server that we were already connected to previously.",
		},
	}
}

// Configure parses and initializes the config.
func (d *Destination) Configure(_ context.Context, cfg config.Config) error {
	commonCfg, err := common.Parse(cfg)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	d.config.Config = commonCfg

	return nil
}

// Open makes sure everything is prepared to receive records.
func (d *Destination) Open(context.Context) error {
	opts, err := common.GetConnectionOptions(d.config.Config)
	if err != nil {
		return fmt.Errorf("get connection options: %s", err)
	}

	conn, err := nats.Connect(strings.Join(d.config.URLs, ","), opts...)
	if err != nil {
		return fmt.Errorf("connect to NATS: %w", err)
	}

	d.writer, err = pubsub.NewWriter(pubsub.WriterParams{
		Conn:    conn,
		Subject: d.config.Subject,
	})
	if err != nil {
		return fmt.Errorf("init pubsub writer: %w", err)
	}

	return nil
}

// Write writes a record into a Destination.
func (d *Destination) Write(_ context.Context, records []opencdc.Record) (int, error) {
	for i, record := range records {
		err := d.writer.Write(record)
		if err != nil {
			return i, fmt.Errorf("write: %w", err)
		}
	}

	return len(records), nil
}

// Teardown gracefully closes connections.
func (d *Destination) Teardown(context.Context) error {
	if d.writer != nil {
		return d.writer.Close()
	}

	return nil
}
