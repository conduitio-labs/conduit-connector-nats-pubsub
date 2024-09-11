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

package source

import (
	"context"
	"fmt"
	"strings"

	"github.com/conduitio-labs/conduit-connector-nats-pubsub/source/pubsub"
	"github.com/conduitio/conduit-commons/config"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/nats-io/nats.go"
)

// Iterator defines an iterator interface.
type Iterator interface {
	HasNext() bool
	Next(ctx context.Context) (opencdc.Record, error)
	Stop() error
}

// Source operates source logic.
type Source struct {
	sdk.UnimplementedSource

	config   Config
	iterator Iterator
	errC     chan error
}

// NewSource creates new instance of the Source.
func NewSource() sdk.Source {
	return sdk.SourceWithMiddleware(&Source{}, sdk.DefaultSourceMiddleware()...)
}

// Parameters returns a map of named config.Parameters that describe how to configure the Source.
func (s *Source) Parameters() config.Parameters {
	return s.config.Parameters()
}

// Configure parses and initializes the config.
func (s *Source) Configure(ctx context.Context, cfg config.Config) error {
	err := sdk.Util.ParseConfig(ctx, cfg, &s.config, NewSource().Parameters())
	if err != nil {
		return err //nolint:wrapcheck // we don't need to wrap the error here
	}

	connName := s.config.GetConnectionName()
	sdk.Logger(ctx).Info().Str("connectionName", connName).Msg("configured connection name")

	return nil
}

// Open opens a connection to NATS and initializes iterators.
func (s *Source) Open(context.Context, opencdc.Position) error {
	s.errC = make(chan error, 1)

	opts, err := s.config.ConnectionOptions()
	if err != nil {
		return fmt.Errorf("get connection options: %w", err)
	}

	conn, err := nats.Connect(strings.Join(s.config.URLs, ","), opts...)
	if err != nil {
		return fmt.Errorf("connect to NATS: %w", err)
	}

	// register an error handler for async errors,
	// the Source listens to them within the Read method and propagates the error if it occurs.
	conn.SetErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
		s.errC <- err
	})

	s.iterator, err = pubsub.NewIterator(pubsub.IteratorParams{
		Conn:       conn,
		BufferSize: s.config.BufferSize,
		Subject:    s.config.Subject,
	})
	if err != nil {
		return fmt.Errorf("init pubsub iterator: %w", err)
	}

	return nil
}

// Read fetches a record from an iterator.
// If there's no record will return sdk.ErrBackoffRetry.
// If the Source's errC is not empty will return the underlying error.
func (s *Source) Read(ctx context.Context) (opencdc.Record, error) {
	select {
	case err := <-s.errC:
		return opencdc.Record{}, fmt.Errorf("got an async error: %w", err)

	default:
		if !s.iterator.HasNext() {
			return opencdc.Record{}, sdk.ErrBackoffRetry
		}

		record, err := s.iterator.Next(ctx)
		if err != nil {
			return opencdc.Record{}, fmt.Errorf("read next record: %w", err)
		}

		return record, nil
	}
}

// Teardown closes connections, stops iterator.
func (s *Source) Teardown(context.Context) error {
	if s.iterator != nil {
		if err := s.iterator.Stop(); err != nil {
			return fmt.Errorf("stop iterator: %w", err)
		}
	}

	return nil
}
