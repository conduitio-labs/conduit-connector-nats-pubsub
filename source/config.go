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

//go:generate paramgen -output=paramgen.go Config

package source

import (
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/common"
)

// Config holds source specific configurable values.
type Config struct {
	common.Config

	// A buffer size for consumed messages.
	BufferSize int `json:"bufferSize" default:"1024" validate:"gt=63"`
}
