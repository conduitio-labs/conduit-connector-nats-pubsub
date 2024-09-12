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

package pubsub

import (
	"context"
	"fmt"
	"time"

	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

// Iterator is a iterator for Pub/Sub communication model.
// It receives any new message from NATS.
type Iterator struct {
	conn         *nats.Conn
	messages     chan *nats.Msg
	subscription *nats.Subscription
}

// IteratorParams contains incoming params for the NewIterator function.
type IteratorParams struct {
	Conn       *nats.Conn
	BufferSize int
	Subject    string
}

// NewIterator creates new instance of the Iterator.
func NewIterator(params IteratorParams) (*Iterator, error) {
	messages := make(chan *nats.Msg, params.BufferSize)

	subscription, err := params.Conn.ChanSubscribe(params.Subject, messages)
	if err != nil {
		return nil, fmt.Errorf("chan subscribe: %w", err)
	}

	return &Iterator{
		conn:         params.Conn,
		messages:     messages,
		subscription: subscription,
	}, nil
}

// HasNext checks is the iterator has messages.
func (i *Iterator) HasNext() bool {
	return len(i.messages) > 0
}

// Next returns the next record from the underlying messages channel.
func (i *Iterator) Next(ctx context.Context) (opencdc.Record, error) {
	select {
	case msg := <-i.messages:
		return i.messageToRecord(msg)

	case <-ctx.Done():
		return opencdc.Record{}, ctx.Err()
	}
}

// Stop stops the Iterator, unsubscribes from a subject.
func (i *Iterator) Stop() (err error) {
	if i.subscription != nil {
		if err = i.subscription.Unsubscribe(); err != nil {
			return fmt.Errorf("unsubscribe: %w", err)
		}
	}

	close(i.messages)

	if i.conn != nil {
		i.conn.Close()
	}

	return nil
}

// messageToRecord converts a *nats.Msg to a opencdc.Record.
func (i *Iterator) messageToRecord(msg *nats.Msg) (opencdc.Record, error) {
	position, err := i.getPosition()
	if err != nil {
		return opencdc.Record{}, fmt.Errorf("get position: %w", err)
	}

	metadata := make(opencdc.Metadata)
	metadata.SetCreatedAt(time.Now())

	return sdk.Util.Source.NewRecordCreate(position, metadata, nil, opencdc.RawData(msg.Data)), nil
}

// getPosition returns the current iterator position.
func (i *Iterator) getPosition() (opencdc.Position, error) {
	uuidBytes, err := uuid.New().MarshalBinary()
	if err != nil {
		return nil, fmt.Errorf("marshal uuid: %w", err)
	}

	return opencdc.Position(uuidBytes), nil
}
