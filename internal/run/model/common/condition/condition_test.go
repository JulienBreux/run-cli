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

package condition

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCondition(t *testing.T) {
	now := time.Now()
	c := Condition{
		Type:               "Ready",
		State:              "True",
		Message:            "Service is ready",
		LastTransitionTime: now,
		Severity:           "Info",
		Reason:             "Completed",
	}

	assert.Equal(t, "Ready", c.Type)
	assert.Equal(t, "True", c.State)
	assert.Equal(t, "Service is ready", c.Message)
	assert.Equal(t, now, c.LastTransitionTime)
	assert.Equal(t, "Info", c.Severity)
	assert.Equal(t, "Completed", c.Reason)
}
