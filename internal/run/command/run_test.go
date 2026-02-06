/*
Copyright 2026 Julien Breux

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package command

import (
	"bytes"
	"errors"
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/config"
	"github.com/stretchr/testify/assert"
)

func TestRunE(t *testing.T) {
	// Save original functions
	origConfigLoad := configLoad
	origAppRun := appRun
	defer func() {
		configLoad = origConfigLoad
		appRun = origAppRun
	}()

	t.Run("Success", func(t *testing.T) {
		configLoad = func() (*config.Config, error) {
			return &config.Config{}, nil
		}
		appRun = func(cfg *config.Config) error {
			return nil
		}

		cmd := New(&bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{})
		err := cmd.RunE(cmd, []string{})
		assert.NoError(t, err)
	})

	t.Run("ConfigLoadError", func(t *testing.T) {
		configLoad = func() (*config.Config, error) {
			return nil, errors.New("config error")
		}

		cmd := New(&bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{})
		err := cmd.RunE(cmd, []string{})
		assert.Error(t, err)
		assert.Equal(t, "config error", err.Error())
	})

	t.Run("AppRunError", func(t *testing.T) {
		configLoad = func() (*config.Config, error) {
			return &config.Config{}, nil
		}
		appRun = func(cfg *config.Config) error {
			return errors.New("app error")
		}

		cmd := New(&bytes.Buffer{}, &bytes.Buffer{}, &bytes.Buffer{})
		err := cmd.RunE(cmd, []string{})
		assert.Error(t, err)
		assert.Equal(t, "app error", err.Error())
	})
}
