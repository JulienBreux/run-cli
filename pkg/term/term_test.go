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

package term

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatTitle(t *testing.T) {
	tests := []struct {
		name     string
		parts    []string
		expected string
	}{
		{
			name:     "no parts",
			parts:    nil,
			expected: "Run",
		},
		{
			name:     "empty parts",
			parts:    []string{},
			expected: "Run",
		},
		{
			name:     "single part",
			parts:    []string{"Services"},
			expected: "Run | Services",
		},
		{
			name:     "multiple parts",
			parts:    []string{"Services", "frontend"},
			expected: "Run | Services | frontend",
		},
		{
			name:     "with empty string in parts",
			parts:    []string{"Services", "", "frontend"},
			expected: "Run | Services | frontend",
		},
		{
			name:     "already starts with Run",
			parts:    []string{"Run", "Services"},
			expected: "Run | Services",
		},
		{
			name:     "only Run in parts",
			parts:    []string{"Run"},
			expected: "Run",
		},
		{
			name:     "modal with three levels",
			parts:    []string{"Services", "frontend", "Scale"},
			expected: "Run | Services | frontend | Scale",
		},
		{
			name:     "loading",
			parts:    []string{"Loading..."},
			expected: "Run | Loading...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, FormatTitle(tt.parts...))
		})
	}
}

func TestSetTitleAndReset(t *testing.T) {
	var buf bytes.Buffer
	Output = &buf

	SetTitle("Run | Services")
	assert.Equal(t, "\033]0;Run | Services\007", buf.String())

	buf.Reset()
	ResetTitle()
	assert.Equal(t, "\033]0;\007", buf.String())

	buf.Reset()
	SetFormattedTitle("Jobs", "backup")
	assert.Equal(t, "\033]0;Run | Jobs | backup\007", buf.String())

	// Test nil output doesn't panic
	Output = nil
	assert.NotPanics(t, func() {
		SetTitle("Test")
	})
}
