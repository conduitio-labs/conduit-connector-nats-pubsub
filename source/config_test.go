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
	"reflect"
	"strings"
	"testing"

	"github.com/conduitio-labs/conduit-connector-nats-pubsub/config"
)

func TestParse(t *testing.T) {
	t.Parallel()

	type args struct {
		cfg map[string]string
	}

	tests := []struct {
		name    string
		args    args
		want    Config
		wantErr bool
	}{
		{
			name: "success, default values",
			args: args{
				cfg: map[string]string{
					config.KeyURLs:    "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
					config.KeySubject: "foo",
				},
			},
			want: Config{
				Config: config.Config{
					URLs:          []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"},
					Subject:       "foo",
					MaxReconnects: config.DefaultMaxReconnects,
					ReconnectWait: config.DefaultReconnectWait,
				},
				BufferSize: defaultBufferSize,
			},
			wantErr: false,
		},
		{
			name: "success, set buffer size",
			args: args{
				cfg: map[string]string{
					config.KeyURLs:      "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
					config.KeySubject:   "foo",
					ConfigKeyBufferSize: "128",
				},
			},
			want: Config{
				Config: config.Config{
					URLs:          []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"},
					Subject:       "foo",
					MaxReconnects: config.DefaultMaxReconnects,
					ReconnectWait: config.DefaultReconnectWait,
				},
				BufferSize: 128,
			},
			wantErr: false,
		},
		{
			name: "success, default buffer size",
			args: args{
				cfg: map[string]string{
					config.KeyURLs:    "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
					config.KeySubject: "foo",
				},
			},
			want: Config{
				Config: config.Config{
					URLs:          []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"},
					Subject:       "foo",
					MaxReconnects: config.DefaultMaxReconnects,
					ReconnectWait: config.DefaultReconnectWait,
				},
				BufferSize: defaultBufferSize,
			},
			wantErr: false,
		},
		{
			name: "fail, invalid buffer size",
			args: args{
				cfg: map[string]string{
					config.KeyURLs:      "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
					config.KeySubject:   "foo",
					ConfigKeyBufferSize: "8",
				},
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "fail, invalid buffer size",
			args: args{
				cfg: map[string]string{
					config.KeyURLs:      "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
					config.KeySubject:   "foo",
					ConfigKeyBufferSize: "what",
				},
			},
			want:    Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if strings.HasPrefix(got.ConnectionName, config.DefaultConnectionNamePrefix) {
				tt.want.ConnectionName = got.ConnectionName
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}
