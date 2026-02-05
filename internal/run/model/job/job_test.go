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

package job

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJob(t *testing.T) {
	now := time.Now()
	j := Job{
		Name:       "job",
		UID:        "uid",
		CreateTime: now,
		Region:     "us-central1",
		Template: &ExecutionTemplate{
			Parallelism: 1,
			TaskCount:   10,
		},
	}

	assert.Equal(t, "job", j.Name)
	assert.Equal(t, "uid", j.UID)
	assert.Equal(t, now, j.CreateTime)
	assert.Equal(t, "us-central1", j.Region)
	assert.Equal(t, int32(1), j.Template.Parallelism)
	assert.Equal(t, int32(10), j.Template.TaskCount)
}

func TestBinaryAuthorization(t *testing.T) {
	ba := BinaryAuthorization{
		UseDefault: true,
		Policy:     "policy",
	}
	assert.True(t, ba.UseDefault)
	assert.Equal(t, "policy", ba.Policy)
}
