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
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/destination"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/test"
	"github.com/conduitio/conduit-commons/opencdc"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/uuid"
)

type driver struct {
	sdk.ConfigurableAcceptanceTestDriver
}

func (d driver) ReadFromDestination(_ *testing.T, records []opencdc.Record) []opencdc.Record {
	return records
}

func (d driver) GenerateRecord(_ *testing.T, operation opencdc.Operation) opencdc.Record {
	id := gofakeit.Int32()

	return opencdc.Record{
		Position:  nil,
		Operation: operation,
		Metadata:  nil,
		Payload: opencdc.Change{
			After: opencdc.RawData([]byte(
				fmt.Sprintf(`"id":%d,"name":"%s"`, id, gofakeit.FirstName()),
			)),
		},
	}
}

//nolint:paralleltest // we don't need the paralleltest here
func TestAcceptance(t *testing.T) {
	cfg := map[string]string{
		destination.ConfigUrls: test.TestURL,
	}

	sdk.AcceptanceTest(t, driver{
		ConfigurableAcceptanceTestDriver: sdk.ConfigurableAcceptanceTestDriver{
			Config: sdk.ConfigurableAcceptanceTestDriverConfig{
				Connector:         Connector,
				SourceConfig:      cfg,
				DestinationConfig: cfg,
				BeforeTest: func(t *testing.T) {
					t.Helper()
					subject := t.Name() + uuid.New().String()
					cfg[destination.ConfigSubject] = subject
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
	})
}
