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
	"testing"

	"github.com/conduitio-labs/conduit-connector-nats-pubsub/config"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/test"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/matryer/is"
)

func TestDestination_Open(t *testing.T) {
	t.Parallel()

	t.Run("success, pubsub", func(t *testing.T) {
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
	})

	t.Run("fail, url is pointed to a non-existent server", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)

		destination := NewDestination()

		err := destination.Configure(context.Background(), map[string]string{
			config.KeyURLs:    "nats://localhost:6666",
			config.KeySubject: "foo_destination",
		})
		is.NoErr(err)

		err = destination.Open(context.Background())
		is.True(err != nil)

		err = destination.Teardown(context.Background())
		is.NoErr(err)
	})
}

func TestDestination_Write(t *testing.T) {
	t.Parallel()

	t.Run("success, pubsub sync write, 1 message", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)

		destination := NewDestination()

		err := destination.Configure(context.Background(), map[string]string{
			config.KeyURLs:    test.TestURL,
			config.KeySubject: "foo_destination_write_pubsub",
		})
		is.NoErr(err)

		err = destination.Open(context.Background())
		is.NoErr(err)

		err = destination.Write(context.Background(), sdk.Record{
			Payload: sdk.RawData([]byte("hello")),
		})
		is.NoErr(err)

		err = destination.Teardown(context.Background())
		is.NoErr(err)
	})
}

func TestDestination_WriteAsync(t *testing.T) {
	t.Parallel()

	t.Run("fail, try to pubsub async write, expected Unimplemented error", func(t *testing.T) {
		t.Parallel()

		is := is.New(t)

		destination := NewDestination()

		err := destination.Configure(context.Background(), map[string]string{
			config.KeyURLs:    test.TestURL,
			config.KeySubject: "foo_destination_write_async_pubsub",
		})
		is.NoErr(err)

		err = destination.Open(context.Background())
		is.NoErr(err)

		err = destination.WriteAsync(context.Background(), sdk.Record{
			Payload: sdk.RawData([]byte("hello")),
		}, func(err error) error {
			return err
		})
		is.Equal(err, sdk.ErrUnimplemented)

		err = destination.Teardown(context.Background())
		is.NoErr(err)
	})
}
