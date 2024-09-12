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

// NewDestination creates new instance of the Destination.
func NewDestination() sdk.Destination {
	return sdk.DestinationWithMiddleware(&Destination{}, sdk.DefaultDestinationMiddleware()...)
}

// Parameters returns a map of named config.Parameters that describe how to configure the Destination.
func (d *Destination) Parameters() config.Parameters {
	return d.config.Parameters()
}

// Configure parses and initializes the config.
func (d *Destination) Configure(ctx context.Context, cfg config.Config) error {
	err := sdk.Util.ParseConfig(ctx, cfg, &d.config, NewDestination().Parameters())
	if err != nil {
		return err //nolint:wrapcheck // we don't need to wrap the error here
	}

	connName := d.config.GetConnectionName()
	sdk.Logger(ctx).Info().Str("connectionName", connName).Msg("configured connection name")

	return nil
}

// Open makes sure everything is prepared to receive records.
func (d *Destination) Open(context.Context) error {
	opts, err := d.config.ConnectionOptions()
	if err != nil {
		return fmt.Errorf("get connection options: %w", err)
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
		err := d.writer.Close()
		if err != nil {
			return fmt.Errorf("failed to close writer: %w", err)
		}
	}

	return nil
}
