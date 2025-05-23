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
	"reflect"
	"testing"
	"time"

	"github.com/conduitio/conduit-commons/opencdc"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func TestPubSubIterator_HasNext(t *testing.T) {
	type fields struct {
		messages chan *nats.Msg
	}

	tests := []struct {
		name     string
		fields   fields
		fillFunc func(chan *nats.Msg)
		want     bool
	}{
		{
			name: "true, one message",
			fields: fields{
				messages: make(chan *nats.Msg, 1),
			},
			fillFunc: func(c chan *nats.Msg) {
				c <- &nats.Msg{
					Subject: "foo",
					Data:    []byte(`"name": "bob"`),
				}
			},
			want: true,
		},
		{
			name: "false, no messages",
			fields: fields{
				messages: make(chan *nats.Msg),
			},
			fillFunc: nil,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				messages: tt.fields.messages,
			}

			if tt.fillFunc != nil {
				tt.fillFunc(tt.fields.messages)
			}

			if got := i.HasNext(); got != tt.want {
				t.Errorf("PubSubIterator.HasNext() = %v, want %v", got, tt.want)

				return
			}
		})
	}
}

func TestPubSubIterator_Next(t *testing.T) {
	type fields struct {
		messages chan *nats.Msg
	}

	tests := []struct {
		name     string
		fields   fields
		timeout  time.Duration
		fillFunc func(chan *nats.Msg)
		want     opencdc.Record
		wantErr  bool
	}{
		{
			name: "success, one message",
			fields: fields{
				messages: make(chan *nats.Msg, 1),
			},
			fillFunc: func(c chan *nats.Msg) {
				c <- &nats.Msg{
					Subject: "foo",
					Data:    []byte(`"name": "bob"`),
				}
			},
			want: opencdc.Record{
				Operation: opencdc.OperationCreate,
				Payload: opencdc.Change{
					After: opencdc.RawData([]byte(`"name": "bob"`)),
				},
			},
			wantErr: false,
		},
		{
			name: "success, no messages, skip",
			fields: fields{
				messages: make(chan *nats.Msg, 1),
			},
			fillFunc: nil,
			want:     opencdc.Record{},
			wantErr:  false,
		},
		{
			name: "success, no messages, context deadline",
			fields: fields{
				messages: make(chan *nats.Msg, 1),
			},
			timeout:  20 * time.Millisecond,
			fillFunc: nil,
			want:     opencdc.Record{},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{
				messages: tt.fields.messages,
			}

			if tt.fillFunc != nil {
				tt.fillFunc(tt.fields.messages)
			}

			if !i.HasNext() {
				return
			}

			for i.HasNext() {
				var ctx context.Context
				var cancel context.CancelFunc

				if tt.timeout != 0 {
					ctx, cancel = context.WithTimeout(context.Background(), tt.timeout)
				} else {
					ctx = context.Background()
				}

				got, err := i.Next(ctx)

				if cancel != nil {
					cancel()
				}

				if (err != nil) != tt.wantErr {
					t.Errorf("PubSubIterator.Next() error = %v, wantErr %v", err, tt.wantErr)

					return
				}

				// we don't care about these fields
				tt.want.Metadata = got.Metadata
				tt.want.Position = got.Position

				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("PubSubIterator.Next() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestPubSubIterator_messageToRecord(t *testing.T) {
	type args struct {
		msg *nats.Msg
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    opencdc.Record
	}{
		{
			name: "success",
			args: args{
				msg: &nats.Msg{
					Subject: "foo",
					Data:    []byte("sample"),
				},
			},
			wantErr: false,
			want: opencdc.Record{
				Operation: opencdc.OperationCreate,
				Payload: opencdc.Change{
					After: opencdc.RawData([]byte("sample")),
				},
			},
		},
		{
			name: "success, nil data",
			args: args{
				msg: &nats.Msg{
					Subject: "foo",
				},
			},
			wantErr: false,
			want: opencdc.Record{
				Operation: opencdc.OperationCreate,
				Payload: opencdc.Change{
					After: opencdc.RawData(nil),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Iterator{}

			got, err := i.messageToRecord(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("PubSubIterator.messageToRecord() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			// we don't care about time
			tt.want.Metadata = got.Metadata

			// check if the position is a valid UUID
			_, err = uuid.FromBytes(got.Position)
			if err != nil {
				t.Errorf("uuid.ParseBytes() = %v", err)

				return
			}

			tt.want.Position = got.Position

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PubSubIterator.messageToRecord() = %v, want %v", got, tt.want)
			}
		})
	}
}
