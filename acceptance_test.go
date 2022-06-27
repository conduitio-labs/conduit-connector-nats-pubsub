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
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/config"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/test"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
	"go.uber.org/goleak"
)

type driver struct {
	sdk.ConfigurableAcceptanceTestDriver

	cache []sdk.Record
}

func (d driver) WriteToSource(t *testing.T, records []sdk.Record) []sdk.Record {
	out := d.ConfigurableAcceptanceTestDriver.WriteToSource(t, records)

	d.cache = append(d.cache, out...)

	return out
}

func (d driver) ReadFromDestination(t *testing.T, records []sdk.Record) []sdk.Record {
	if len(d.cache) != len(records) {
		return records
	}

	out := make([]sdk.Record, 0, cap(records))
	for i := 0; i < cap(records); i++ {
		out = append(out, d.cache[i])
		d.cache = d.cache[1:]
	}

	return out
}

func (d driver) GenerateRecord(t *testing.T) sdk.Record {
	id := gofakeit.Int32()

	return sdk.Record{
		Position: nil,
		Metadata: nil,
		Payload: sdk.RawData([]byte(
			fmt.Sprintf(`"id":%d,"name":"%s"`, id, gofakeit.FirstName()),
		)),
	}
}

//nolint:paralleltest // we don't need the paralleltest here
func TestAcceptance(t *testing.T) {
	cfg := map[string]string{
		config.ConfigKeyURLs: test.TestURL,
	}

	sdk.AcceptanceTest(t, driver{
		ConfigurableAcceptanceTestDriver: sdk.ConfigurableAcceptanceTestDriver{
			Config: sdk.ConfigurableAcceptanceTestDriverConfig{
				Connector:         Connector,
				SourceConfig:      cfg,
				DestinationConfig: cfg,
				BeforeTest: func(t *testing.T) {
					subject := t.Name() + uuid.New().String()
					cfg[config.ConfigKeySubject] = subject
				},
				GoleakOptions: []goleak.Option{
					// nats.go spawns a separate goroutine to process flush requests
					goleak.IgnoreTopFunction("github.com/nats-io/nats%2ego.(*Conn).flusher"),
					goleak.IgnoreTopFunction("sync.runtime_notifyListWait"),
					goleak.IgnoreTopFunction("internal/poll.runtime_pollWait"),
				},
				Skip: []string{
					// NATS PubSub doesn't handle position
					"TestSource_Open_ResumeAtPositionSnapshot",
					"TestSource_Open_ResumeAtPositionCDC",
					// NATS PubSub doesn't persist messages,
					// so if there's no readers messages are deleted immediately.
					"TestSource_Read_Success",
				},
			},
		},
		cache: make([]sdk.Record, 0),
	})
}
