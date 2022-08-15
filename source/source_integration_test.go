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
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/conduitio-labs/conduit-connector-nats-pubsub/config"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/test"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/nats-io/nats.go"
)

func TestSource_Open(t *testing.T) {
	t.Parallel()

	source := NewSource()
	err := source.Configure(context.Background(), map[string]string{
		config.KeyURLs:           test.TestURL,
		config.KeySubject:        "foo",
		config.KeyConnectionName: "super_connection",
	})
	if err != nil {
		t.Fatalf("configure source: %v", err)

		return
	}

	err = source.Open(context.Background(), sdk.Position(nil))
	if err != nil {
		t.Fatalf("open source: %v", err)

		return
	}

	err = source.Teardown(context.Background())
	if err != nil {
		t.Fatalf("teardown source: %v", err)

		return
	}
}

func TestSource_OpenWithPassword(t *testing.T) {
	t.Parallel()

	source := NewSource()
	err := source.Configure(context.Background(), map[string]string{
		config.KeyURLs:           test.TestURLWithPassword,
		config.KeySubject:        "foo",
		config.KeyConnectionName: "super_connection",
	})
	if err != nil {
		t.Fatalf("configure source: %v", err)

		return
	}

	err = source.Open(context.Background(), sdk.Position(nil))
	if err != nil {
		t.Fatalf("open source: %v", err)

		return
	}

	err = source.Teardown(context.Background())
	if err != nil {
		t.Fatalf("teardown source: %v", err)

		return
	}
}

func TestSource_OpenWithToken(t *testing.T) {
	t.Parallel()

	source := NewSource()
	err := source.Configure(context.Background(), map[string]string{
		config.KeyURLs:           test.TestURLWithToken,
		config.KeySubject:        "foo",
		config.KeyConnectionName: "super_connection",
	})
	if err != nil {
		t.Fatalf("configure source: %v", err)

		return
	}

	err = source.Open(context.Background(), sdk.Position(nil))
	if err != nil {
		t.Fatalf("open source: %v", err)

		return
	}

	err = source.Teardown(context.Background())
	if err != nil {
		t.Fatalf("teardown source: %v", err)

		return
	}
}

func TestSource_OpenWithNKey(t *testing.T) {
	t.Parallel()

	source := NewSource()
	err := source.Configure(context.Background(), map[string]string{
		config.KeyURLs:           test.TestURLWithNKey,
		config.KeySubject:        "foo",
		config.KeyConnectionName: "super_connection",
		config.KeyNKeyPath:       "../test/fixtures/test_nkey_seed.txt",
	})
	if err != nil {
		t.Fatalf("configure source: %v", err)

		return
	}

	err = source.Open(context.Background(), sdk.Position(nil))
	if err != nil {
		t.Fatalf("open source: %v", err)

		return
	}

	err = source.Teardown(context.Background())
	if err != nil {
		t.Fatalf("teardown source: %v", err)

		return
	}
}

func TestSource_ReadPubSubSuccessOneMessage(t *testing.T) {
	t.Parallel()

	subject := "foo_one"

	source, err := createTestPubSub(t, map[string]string{
		config.KeyURLs:    test.TestURL,
		config.KeySubject: subject,
	})
	if err != nil {
		t.Fatalf("create test pubsub: %v", err)

		return
	}

	t.Cleanup(func() {
		if err := source.Teardown(context.Background()); err != nil {
			t.Fatalf("teardown source: %v", err)
		}
	})

	testConn, err := test.GetTestConnection(test.TestURL)
	if err != nil {
		t.Fatalf("get test connection: %v", err)

		return
	}

	err = testConn.Publish(subject, []byte(`{"level": "info"}`))
	if err != nil {
		t.Fatalf("publish message: %v", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var record sdk.Record
	for {
		record, err = source.Read(ctx)
		if err != nil {
			if errors.Is(err, sdk.ErrBackoffRetry) {
				continue
			}
			t.Fatalf("read message: %v", err)

			return
		}

		break
	}

	if !reflect.DeepEqual(record.Payload.Bytes(), []byte(`{"level": "info"}`)) {
		t.Fatalf("Source.Read = %v, want %v", record.Payload.Bytes(), []byte(`{"level": "info"}`))

		return
	}
}

func TestSource_ReadPubSubSuccessManyMessage(t *testing.T) {
	t.Parallel()

	subject := "foo_many"

	source, err := createTestPubSub(t, map[string]string{
		config.KeyURLs:    test.TestURL,
		config.KeySubject: subject,
	})
	if err != nil {
		t.Fatalf("create test pubsub: %v", err)

		return
	}

	t.Cleanup(func() {
		if err := source.Teardown(context.Background()); err != nil {
			t.Fatalf("teardown source: %v", err)
		}
	})

	testConn, err := test.GetTestConnection(test.TestURL)
	if err != nil {
		t.Fatalf("get test connection: %v", err)

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	records := make([]sdk.Record, 0)
	for i := 0; i < 128; i++ {
		err = testConn.Publish(subject, []byte(`{"level": "info"}`))
		if err != nil {
			t.Fatalf("publish message: %v", err)

			return
		}

		record, err := source.Read(ctx)
		if err != nil {
			if errors.Is(err, sdk.ErrBackoffRetry) {
				i--

				continue
			}
			t.Fatalf("read message: %v", err)

			return
		}

		records = append(records, record)
	}

	if len(records) != 128 {
		t.Fatalf("len(records) = %d, expected = %d", len(records), 128)

		return
	}
}

func TestSource_ReadPubSubSuccessNoMessagesBackoff(t *testing.T) {
	t.Parallel()

	subject := "no_messages"

	source, err := createTestPubSub(t, map[string]string{
		config.KeyURLs:    test.TestURL,
		config.KeySubject: subject,
	})
	if err != nil {
		t.Fatalf("create test pubsub: %v", err)

		return
	}

	t.Cleanup(func() {
		if err := source.Teardown(context.Background()); err != nil {
			t.Fatalf("teardown source: %v", err)
		}
	})

	_, err = source.Read(context.Background())
	if err == nil {
		t.Fatal("Source.Read expected backoff retry error, got nil")

		return
	}

	if err != nil && !errors.Is(err, sdk.ErrBackoffRetry) {
		t.Fatalf("read message: %v", err)

		return
	}
}

func TestSource_ReadPubSubManyMessagesSlowConsumerErr(t *testing.T) {
	t.Parallel()

	subject := "slow_consumers_subj"

	source, err := createTestPubSub(t, map[string]string{
		config.KeyURLs:      test.TestURL,
		config.KeySubject:   subject,
		ConfigKeyBufferSize: "64",
	})
	if err != nil {
		t.Fatalf("create test pubsub: %v", err)

		return
	}

	t.Cleanup(func() {
		if err := source.Teardown(context.Background()); err != nil {
			t.Fatalf("teardown source: %v", err)
		}
	})

	testConn, err := test.GetTestConnection(test.TestURL)
	if err != nil {
		t.Fatalf("get test connection: %v", err)

		return
	}

	for i := 0; i < 1_000_000; i++ {
		err = testConn.Publish(subject, []byte(`{"level": "info"}`))
		if err != nil {
			t.Fatalf("publish test mesage: %v", err)

			return
		}

		_, err := source.Read(context.Background())
		if err != nil {
			if errors.Is(err, sdk.ErrBackoffRetry) {
				continue
			}

			if !errors.Is(errors.Unwrap(err), nats.ErrSlowConsumer) {
				t.Fatalf("Source.Read expected slow consumer error, got %v", err)

				return
			}

			return
		}

		continue
	}

	t.Fatalf("Source.Read didn't get the expected slow consumer error")
}

func createTestPubSub(t *testing.T, cfg map[string]string) (sdk.Source, error) {
	source := NewSource()

	err := source.Configure(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("configure source: %w", err)
	}

	err = source.Open(context.Background(), sdk.Position(nil))
	if err != nil {
		return nil, fmt.Errorf("open source: %w", err)
	}

	return source, nil
}
