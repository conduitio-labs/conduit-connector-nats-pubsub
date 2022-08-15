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
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/config"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/destination/pubsub"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/nats-io/nats.go"
)

// Writer defines a writer interface needed for the Destination.
type Writer interface {
	Write(ctx context.Context, record sdk.Record) error
	Close(ctx context.Context) error
}

// Destination NATS Connector sends records to a NATS subject.
type Destination struct {
	sdk.UnimplementedDestination

	config config.Config
	writer Writer
}

// NewDestination creates new instance of the Destination.
func NewDestination() sdk.Destination {
	return &Destination{}
}

// Configure parses and initializes the config.
func (d *Destination) Configure(ctx context.Context, cfg map[string]string) error {
	config, err := config.Parse(cfg)
	if err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	d.config = config

	return nil
}

// Open makes sure everything is prepared to receive records.
func (d *Destination) Open(ctx context.Context) error {
	opts, err := common.GetConnectionOptions(d.config)
	if err != nil {
		return fmt.Errorf("get connection options: %s", err)
	}

	conn, err := nats.Connect(strings.Join(d.config.URLs, ","), opts...)
	if err != nil {
		return fmt.Errorf("connect to NATS: %w", err)
	}

	d.writer, err = pubsub.NewWriter(ctx, pubsub.WriterParams{
		Conn:    conn,
		Subject: d.config.Subject,
	})
	if err != nil {
		return fmt.Errorf("init pubsub writer: %w", err)
	}

	return nil
}

// Write writes a record into a Destination.
func (d *Destination) Write(ctx context.Context, record sdk.Record) error {
	return d.writer.Write(ctx, record)
}

// Teardown gracefully closes connections.
func (d *Destination) Teardown(ctx context.Context) error {
	if d.writer != nil {
		return d.writer.Close(ctx)
	}

	return nil
}
