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

package execution

import (
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/stretchr/testify/assert"
)

func TestExecution(t *testing.T) {
	now := time.Now()
	e := Execution{
		Name:           "execution",
		Job:            "job",
		CreateTime:     now,
		StartTime:      now,
		CompletionTime: now,
		TaskCount:      10,
		SucceededCount: 10,
		Region:         "us-central1",
		TerminalCondition: &condition.Condition{
			Type:  "Ready",
			State: "True",
		},
	}

	assert.Equal(t, "execution", e.Name)
	assert.Equal(t, "job", e.Job)
	assert.Equal(t, now, e.CreateTime)
	assert.Equal(t, int32(10), e.TaskCount)
	assert.Equal(t, int32(10), e.SucceededCount)
	assert.Equal(t, "us-central1", e.Region)
	assert.Equal(t, "Ready", e.TerminalCondition.Type)
}
