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

package instance

import (
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/container"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/stretchr/testify/assert"
)

func TestInstance(t *testing.T) {
	now := time.Now()
	inst := Instance{
		Name:          "projects/test-project/locations/us-central1/instances/my-instance-123",
		Description:   "Test instance",
		UID:           "test-uid-123",
		Generation:    1,
		CreateTime:    now,
		UpdateTime:    now,
		Creator:       "user@example.com",
		Region:        "us-central1",
		Project:       "test-project",
		RestartPolicy: "RESTART_POLICY_ALWAYS",
		Containers: []*container.Container{
			{
				Name:  "web",
				Image: "gcr.io/test/web:latest",
			},
		},
		ContainerStatuses: []*ContainerStatus{
			{
				Name:        "web",
				ImageDigest: "sha256:abcdef1234567890",
			},
		},
		URLs: []string{
			"https://my-instance-123-uc.a.run.app",
		},
		TerminalCondition: &condition.Condition{
			Type:    "Ready",
			State:   "CONDITION_SUCCEEDED",
			Message: "Instance is ready",
		},
	}

	assert.Equal(t, "my-instance-123", inst.ID())
	assert.Equal(t, "Ready", inst.Status())
	assert.Equal(t, "Test instance", inst.Description)
	assert.Equal(t, "test-uid-123", inst.UID)
	assert.Equal(t, int64(1), inst.Generation)
	assert.Equal(t, now, inst.CreateTime)
	assert.Equal(t, "us-central1", inst.Region)
	assert.Equal(t, "test-project", inst.Project)
	assert.Equal(t, "RESTART_POLICY_ALWAYS", inst.RestartPolicy)
	assert.Equal(t, 1, len(inst.Containers))
	assert.Equal(t, "web", inst.Containers[0].Name)
	assert.Equal(t, 1, len(inst.ContainerStatuses))
	assert.Equal(t, "sha256:abcdef1234567890", inst.ContainerStatuses[0].ImageDigest)
	assert.Equal(t, "https://my-instance-123-uc.a.run.app", inst.FormattedURLs())
}

func TestInstanceStatus(t *testing.T) {
	tests := []struct {
		name     string
		instance Instance
		expected string
	}{
		{
			name: "Reconciling",
			instance: Instance{
				Reconciling: true,
			},
			expected: "Reconciling",
		},
		{
			name: "Terminal Condition Succeeded",
			instance: Instance{
				TerminalCondition: &condition.Condition{
					State: "CONDITION_SUCCEEDED",
				},
			},
			expected: "Ready",
		},
		{
			name: "Terminal Condition Failed",
			instance: Instance{
				TerminalCondition: &condition.Condition{
					State: "CONDITION_FAILED",
				},
			},
			expected: "Failed",
		},
		{
			name: "Terminal Condition Pending",
			instance: Instance{
				TerminalCondition: &condition.Condition{
					State: "CONDITION_PENDING",
				},
			},
			expected: "Pending",
		},
		{
			name: "Condition Ready Succeeded",
			instance: Instance{
				Conditions: []*condition.Condition{
					{
						Type:  "Ready",
						State: "CONDITION_SUCCEEDED",
					},
				},
			},
			expected: "Ready",
		},
		{
			name: "Condition Ready Failed",
			instance: Instance{
				Conditions: []*condition.Condition{
					{
						Type:  "Ready",
						State: "CONDITION_FAILED",
					},
				},
			},
			expected: "Failed",
		},
		{
			name:     "Unknown",
			instance: Instance{},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.instance.Status())
		})
	}
}

func TestInstanceID(t *testing.T) {
	inst1 := Instance{Name: "projects/p/locations/r/instances/inst-1"}
	assert.Equal(t, "inst-1", inst1.ID())

	inst2 := Instance{Name: "simple-name"}
	assert.Equal(t, "simple-name", inst2.ID())

	inst3 := Instance{Name: ""}
	assert.Equal(t, "", inst3.ID())
}

func TestInstanceProxy(t *testing.T) {
	inst := Instance{
		Name: "test-instance",
		Proxy: &model_service.ProxyStatus{
			Enabled: true,
			Port:    8080,
			URL:     "http://127.0.0.1:8080",
		},
	}
	assert.NotNil(t, inst.Proxy)
	assert.True(t, inst.Proxy.Enabled)
	assert.Equal(t, 8080, inst.Proxy.Port)
	assert.Equal(t, "http://127.0.0.1:8080", inst.Proxy.URL)
}
