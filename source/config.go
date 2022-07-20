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
	"fmt"

	"strconv"

	"github.com/conduitio-labs/conduit-connector-nats-pubsub/config"
	"github.com/conduitio-labs/conduit-connector-nats-pubsub/validator"
)

// defaultBufferSize is a default buffer size for consumed messages.
// It must be set to avoid the problem with slow consumers.
// See details about slow consumers here https://docs.nats.io/using-nats/developer/connecting/events/slow.
const defaultBufferSize = 1024

// ConfigKeyBufferSize is a config name for a buffer size.
const ConfigKeyBufferSize = "bufferSize"

// Config holds source specific configurable values.
type Config struct {
	config.Config

	BufferSize int `key:"bufferSize" validate:"omitempty,min=64"`
}

// Parse maps the incoming map to the Config and validates it.
func Parse(cfg map[string]string) (Config, error) {
	common, err := config.Parse(cfg)
	if err != nil {
		return Config{}, fmt.Errorf("parse common config: %w", err)
	}

	sourceConfig := Config{
		Config: common,
	}

	if err := sourceConfig.parseBufferSize(cfg[ConfigKeyBufferSize]); err != nil {
		return Config{}, fmt.Errorf("parse buffer size: %w", err)
	}

	sourceConfig.setDefaults()

	if err := validator.Validate(&sourceConfig); err != nil {
		return Config{}, fmt.Errorf("validate source config: %w", err)
	}

	return sourceConfig, nil
}

// parseBufferSize parses the bufferSize string and
// if it's not empty set cfg.BufferSize to its integer representation.
func (c *Config) parseBufferSize(bufferSizeStr string) error {
	if bufferSizeStr != "" {
		bufferSize, err := strconv.Atoi(bufferSizeStr)
		if err != nil {
			return fmt.Errorf("\"%s\" must be an integer", ConfigKeyBufferSize)
		}

		c.BufferSize = bufferSize
	}

	return nil
}

// setDefaults set default values for empty fields.
func (c *Config) setDefaults() {
	if c.BufferSize == 0 {
		c.BufferSize = defaultBufferSize
	}
}
