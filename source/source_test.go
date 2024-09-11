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
	"testing"

	"github.com/conduitio/conduit-commons/config"
	"github.com/matryer/is"
)

func TestSource_Configure(t *testing.T) {
	tests := []struct {
		name    string
		cfg     config.Config
		wantErr bool
	}{
		{
			name: "success, correct config",
			cfg: map[string]string{
				ConfigUrls:    "nats://127.0.0.1:4222",
				ConfigSubject: "foo",
			},
			wantErr: false,
		},
		{
			name:    "fail, empty config",
			cfg:     map[string]string{},
			wantErr: true,
		},
		{
			name: "fail, invalid config",
			cfg: map[string]string{
				ConfigUrls: "nats://127.0.0.1:4222",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			is := is.New(t)
			s := NewSource()
			err := s.Configure(context.Background(), tt.cfg)

			if tt.wantErr {
				is.True(err != nil)
				return
			}
			is.NoErr(err)
		})
	}
}
