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

package workerpool

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorkerPool(t *testing.T) {
	now := time.Now()
	wp := WorkerPool{
		Name:        "pool",
		CreateTime:  now,
		State:       "ACTIVE",
		DisplayName: "My Pool",
		Region:      "us-central1",
		WorkerConfig: &WorkerConfig{
			MachineType: "e2-standard-4",
			DiskSizeGb:  100,
		},
	}

	assert.Equal(t, "pool", wp.Name)
	assert.Equal(t, now, wp.CreateTime)
	assert.Equal(t, "ACTIVE", wp.State)
	assert.Equal(t, "My Pool", wp.DisplayName)
	assert.Equal(t, "us-central1", wp.Region)
	assert.Equal(t, "e2-standard-4", wp.WorkerConfig.MachineType)
	assert.Equal(t, int32(100), wp.WorkerConfig.DiskSizeGb)
}
