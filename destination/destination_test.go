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

	"github.com/conduitio-labs/conduit-connector-nats-pubsub/common"
)

func TestDestination_Configure(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		cfg map[string]string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success, correct config",
			args: args{
				ctx: context.Background(),
				cfg: map[string]string{
					common.KeyURLs:    "nats://127.0.0.1:4222",
					common.KeySubject: "foo",
				},
			},
			wantErr: false,
		},
		{
			name: "fail, empty config",
			args: args{
				ctx: context.Background(),
				cfg: map[string]string{},
			},
			wantErr: true,
		},
		{
			name: "fail, invalid config",
			args: args{
				ctx: context.Background(),
				cfg: map[string]string{
					common.KeyURLs: "nats://127.0.0.1:4222",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := &Destination{}
			if err := s.Configure(tt.args.ctx, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Source.Configure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
