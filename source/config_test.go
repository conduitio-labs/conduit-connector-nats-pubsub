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
	"strings"
	"testing"
	"time"

	"github.com/conduitio-labs/conduit-connector-nats-pubsub/common"
	sdk "github.com/conduitio/conduit-connector-sdk"
	"github.com/google/go-cmp/cmp"
	"github.com/matryer/is"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		cfg     map[string]string
		want    Config
		wantErr bool
	}{
		{
			name: "success, only required fields provided, many connection URLs",
			cfg: map[string]string{
				ConfigUrls:    "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
				ConfigSubject: "foo",
			},
			want: Config{
				Config: common.Config{
					URLs:          []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"},
					Subject:       "foo",
					MaxReconnects: 5,
					ReconnectWait: time.Second * 5,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "success, only required fields provided, one connection URL",
			cfg: map[string]string{
				ConfigUrls:    "nats://127.0.0.1:1222",
				ConfigSubject: "foo",
			},
			want: Config{
				Config: common.Config{
					URLs:          []string{"nats://127.0.0.1:1222"},
					Subject:       "foo",
					MaxReconnects: 5,
					ReconnectWait: time.Second * 5,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "success, url with token",
			cfg: map[string]string{
				ConfigUrls:    "nats://token:127.0.0.1:1222",
				ConfigSubject: "foo",
			},
			want: Config{
				Config: common.Config{
					URLs:          []string{"nats://token:127.0.0.1:1222"},
					Subject:       "foo",
					MaxReconnects: 5,
					ReconnectWait: time.Second * 5,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "success, url with user/password",
			cfg: map[string]string{
				ConfigUrls:    "nats://admin:admin@127.0.0.1:1222",
				ConfigSubject: "foo",
			},
			want: Config{
				Config: common.Config{
					URLs:          []string{"nats://admin:admin@127.0.0.1:1222"},
					Subject:       "foo",
					MaxReconnects: 5,
					ReconnectWait: time.Second * 5,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "fail, required field (subject) is missing",
			cfg: map[string]string{
				ConfigUrls: "nats://localhost:1222",
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "success, nkey pair",
			cfg: map[string]string{
				ConfigUrls:     "nats://127.0.0.1:1222",
				ConfigSubject:  "foo",
				ConfigNkeyPath: "./config.go",
			},
			want: Config{
				Config: common.Config{
					URLs:          []string{"nats://127.0.0.1:1222"},
					Subject:       "foo",
					NKeyPath:      "./config.go",
					MaxReconnects: 5,
					ReconnectWait: time.Second * 5,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "success, credentials file",
			cfg: map[string]string{
				ConfigUrls:                "nats://127.0.0.1:1222",
				ConfigSubject:             "foo",
				ConfigCredentialsFilePath: "./config.go",
			},
			want: Config{
				Config: common.Config{
					URLs:                []string{"nats://127.0.0.1:1222"},
					Subject:             "foo",
					CredentialsFilePath: "./config.go",
					MaxReconnects:       5,
					ReconnectWait:       time.Second * 5,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "success, custom connection name",
			cfg: map[string]string{
				ConfigUrls:           "nats://127.0.0.1:1222",
				ConfigSubject:        "foo",
				ConfigConnectionName: "my_super_connection",
			},
			want: Config{
				Config: common.Config{
					URLs:           []string{"nats://127.0.0.1:1222"},
					Subject:        "foo",
					ConnectionName: "my_super_connection",
					MaxReconnects:  5,
					ReconnectWait:  time.Second * 5,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "success, custom reconnect options",
			cfg: map[string]string{
				ConfigUrls:          "nats://127.0.0.1:1222",
				ConfigSubject:       "foo",
				ConfigMaxReconnects: "20",
				ConfigReconnectWait: "10s",
			},
			want: Config{
				Config: common.Config{
					URLs:          []string{"nats://127.0.0.1:1222"},
					Subject:       "foo",
					MaxReconnects: 20,
					ReconnectWait: time.Second * 10,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "success, default values",
			cfg: map[string]string{
				ConfigUrls:    "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
				ConfigSubject: "foo",
			},
			want: Config{
				Config: common.Config{
					URLs:          []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"},
					Subject:       "foo",
					MaxReconnects: 5,
					ReconnectWait: time.Second * 5,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "success, set buffer size",
			cfg: map[string]string{
				ConfigUrls:       "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
				ConfigSubject:    "foo",
				ConfigBufferSize: "128",
			},
			want: Config{
				Config: common.Config{
					URLs:          []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"},
					Subject:       "foo",
					MaxReconnects: 5,
					ReconnectWait: time.Second * 5,
				},
				BufferSize: 128,
			},
			wantErr: false,
		},
		{
			name: "success, default buffer size",
			cfg: map[string]string{
				ConfigUrls:    "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
				ConfigSubject: "foo",
			},
			want: Config{
				Config: common.Config{
					URLs:          []string{"nats://127.0.0.1:1222", "nats://127.0.0.1:1223", "nats://127.0.0.1:1224"},
					Subject:       "foo",
					MaxReconnects: 5,
					ReconnectWait: time.Second * 5,
				},
				BufferSize: 1024,
			},
			wantErr: false,
		},
		{
			name: "fail, invalid buffer size",
			cfg: map[string]string{
				ConfigUrls:       "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
				ConfigSubject:    "foo",
				ConfigBufferSize: "8",
			},
			want:    Config{},
			wantErr: true,
		},
		{
			name: "fail, invalid buffer size",
			cfg: map[string]string{
				ConfigUrls:       "nats://127.0.0.1:1222,nats://127.0.0.1:1223,nats://127.0.0.1:1224",
				ConfigSubject:    "foo",
				ConfigBufferSize: "what",
			},
			want:    Config{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			var got Config
			err := sdk.Util.ParseConfig(context.Background(), tt.cfg, &got, NewSource().Parameters())
			got.GetConnectionName() // initialize connection name

			if tt.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)

			if tt.want.ConnectionName == "" {
				if !strings.HasPrefix(got.ConnectionName, common.DefaultConnectionNamePrefix) {
					is.Equal(common.DefaultConnectionNamePrefix, got.ConnectionName) // expected default connection name prefix
				}
				tt.want.ConnectionName = got.ConnectionName
			}

			is.Equal("", cmp.Diff(tt.want, got))
		})
	}
}
