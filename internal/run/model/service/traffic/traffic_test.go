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

package traffic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraffic(t *testing.T) {
	tt := TrafficTarget{
		Type:     TrafficTargetAllocationTypeLatest,
		Revision: "rev1",
		Percent:  100,
		Tag:      "latest",
	}

	assert.Equal(t, TrafficTargetAllocationTypeLatest, tt.Type)
	assert.Equal(t, "rev1", tt.Revision)
	assert.Equal(t, int32(100), tt.Percent)
	assert.Equal(t, "latest", tt.Tag)

	tts := TrafficTargetStatus{
		Revision: "rev1",
		Percent:  100,
		URI:      "https://rev1.service.run.app",
	}

	assert.Equal(t, "rev1", tts.Revision)
	assert.Equal(t, int32(100), tts.Percent)
	assert.Equal(t, "https://rev1.service.run.app", tts.URI)
}
