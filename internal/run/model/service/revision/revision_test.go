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

package revision

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRevision_NewFields(t *testing.T) {
	r := Revision{
		Name:         "rev-1",
		TrafficShare: 50,
		Author:       "user@example.com",
	}

	assert.Equal(t, 50, r.TrafficShare)
	assert.Equal(t, "user@example.com", r.Author)
}
