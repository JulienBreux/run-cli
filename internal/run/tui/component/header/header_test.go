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

package header_test

import (
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	testInfo := info.Info{
		Project: "test-project",
		Region:  "us-central1",
		User:    "test-user",
	}

	h := header.New(testInfo)

	assert.NotNil(t, h)
	assert.Equal(t, 3, h.GetItemCount())
}

func TestUpdateInfo(t *testing.T) {
	// Initialize the component first to set the global infoView variable
	initialInfo := info.Info{
		Project: "p1",
		Region:  "r1",
		User:    "u1",
	}
	_ = header.New(initialInfo)

	// Now update it
	newInfo := info.Info{
		Project: "p2",
		Region:  "r2",
		User:    "u2",
	}

	// This function modifies the global infoView.
	// Since we can't inspect the text content easily without drawing, we just ensure it doesn't panic.
	assert.NotPanics(t, func() {
		header.UpdateInfo(newInfo)
	})

	// Note: Testing side effects on global variables is brittle in parallel tests,
	// but acceptable here given the legacy code structure.
}
