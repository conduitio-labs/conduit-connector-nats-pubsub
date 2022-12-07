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
	"testing"
	"time"

	"github.com/conduitio-labs/conduit-connector-nats-pubsub/config"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/test"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
	"github.com/nats-io/nats.go"
)

func TestDestination_OpenSuccess(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	destination := NewDestination()

	err := destination.Configure(context.Background(), map[string]string{
		config.KeyURLs:    test.TestURL,
		config.KeySubject: "foo_destination",
	})
	is.NoErr(err)

	err = destination.Open(context.Background())
	is.NoErr(err)

	err = destination.Teardown(context.Background())
	is.NoErr(err)
}

func TestDestination_OpenFail(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	destination := NewDestination()

	err := destination.Configure(context.Background(), map[string]string{
		config.KeyURLs:          "nats://localhost:6666",
		config.KeySubject:       "foo_destination",
		config.KeyMaxReconnects: "0",
		config.KeyReconnectWait: "2s",
	})
	is.NoErr(err)

	err = destination.Open(context.Background())
	is.True(err != nil)

	err = destination.Teardown(context.Background())
	is.NoErr(err)
}

func TestDestination_WriteOneMessage(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	subject := "foo_destination_write_one_pubsub"

	testConn, err := nats.Connect(test.TestURL)
	is.NoErr(err)

	t.Cleanup(func() {
		testConn.Flush()
		testConn.Close()
	})

	subscription, err := testConn.SubscribeSync(subject)
	is.NoErr(err)

	destination := NewDestination()

	err = destination.Configure(context.Background(), map[string]string{
		config.KeyURLs:    test.TestURL,
		config.KeySubject: subject,
	})
	is.NoErr(err)

	err = destination.Open(context.Background())
	is.NoErr(err)

	var count int
	count, err = destination.Write(context.Background(), []sdk.Record{
		{
			Operation: sdk.OperationCreate,
			Payload: sdk.Change{
				After: sdk.RawData([]byte("hello")),
			},
		},
	})
	is.NoErr(err)

	msg, err := subscription.NextMsg(time.Second * 2)
	is.NoErr(err)

	is.Equal(count, 1)
	is.Equal(msg.Data, []byte("hello"))

	err = destination.Teardown(context.Background())
	is.NoErr(err)
}

func TestDestination_WriteManyMessages(t *testing.T) {
	t.Parallel()

	is := is.New(t)

	subject := "foo_destination_write_many_pubsub"

	testConn, err := nats.Connect(test.TestURL)
	is.NoErr(err)

	t.Cleanup(func() {
		testConn.Flush()
		testConn.Close()
	})

	subscription, err := testConn.SubscribeSync(subject)
	is.NoErr(err)

	destination := NewDestination()

	err = destination.Configure(context.Background(), map[string]string{
		config.KeyURLs:    test.TestURL,
		config.KeySubject: subject,
	})
	is.NoErr(err)

	err = destination.Open(context.Background())
	is.NoErr(err)

	records := make([]sdk.Record, 1000)
	for i := 0; i < 1000; i++ {
		records[i] = sdk.Record{
			Operation: sdk.OperationCreate,
			Payload: sdk.Change{
				After: sdk.RawData([]byte(fmt.Sprintf("message #%d", i))),
			},
		}
	}

	var count int
	count, err = destination.Write(context.Background(), records)
	is.NoErr(err)
	is.Equal(count, 1000)

	messages := make([]*nats.Msg, 0, 1000)
	for i := 0; i < 1000; i++ {
		message, err := subscription.NextMsg(time.Second * 2)
		is.NoErr(err)

		messages = append(messages, message)
	}

	for i, message := range messages {
		is.Equal(message.Data, []byte(fmt.Sprintf("message #%d", i)))
	}

	err = destination.Teardown(context.Background())
	is.NoErr(err)
}
