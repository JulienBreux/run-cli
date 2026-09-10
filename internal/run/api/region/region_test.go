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

package region

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, "all", ALL)
}

func TestList(t *testing.T) {
	regions := List()
	assert.NotEmpty(t, regions)
	assert.Contains(t, regions, "us-central1")
	assert.Contains(t, regions, "europe-west1")

	// Ensure sorted
	assert.True(t, sort.StringsAreSorted(regions), "regions list should be sorted alphabetically")

	// Ensure no duplicates and valid format
	seen := make(map[string]bool)
	for _, r := range regions {
		assert.NotEmpty(t, r)
		assert.False(t, seen[r], "duplicate region found: %s", r)
		seen[r] = true
		assert.Regexp(t, "^[a-z0-9-]+$", r, "region name should be lowercase alphanumeric with hyphens")
	}
}
