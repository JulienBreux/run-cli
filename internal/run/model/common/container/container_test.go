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

package container

import (
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/model/common/env"
	"github.com/JulienBreux/run-cli/internal/run/model/common/resources"
	"github.com/stretchr/testify/assert"
)

func TestContainer(t *testing.T) {
	c := Container{
		Name:  "container",
		Image: "image",
		Env: []*env.EnvVar{
			{Name: "VAR", Value: "VAL"},
		},
		Resources: &resources.Resources{
			Limits: map[string]string{"cpu": "1"},
		},
		Ports: []*Port{
			{Name: "http", ContainerPort: 8080},
		},
		LivenessProbe: &Probe{
			HTTPGet: &HTTPGetAction{
				Path: "/health",
				Port: 8080,
				HTTPHeaders: []*HTTPHeader{
					{Name: "X-Header", Value: "Value"},
				},
			},
		},
		Exec: &ExecAction{
			Command: []string{"ls"},
		},
		TCPSocket: &TCPSocketAction{
			Port: 8080,
		},
		VolumeMounts: []*VolumeMount{
			{Name: "vol", MountPath: "/mnt"},
		},
	}

	assert.Equal(t, "container", c.Name)
	assert.Equal(t, "image", c.Image)
	assert.Len(t, c.Env, 1)
	assert.NotNil(t, c.Resources)
	assert.Len(t, c.Ports, 1)
	assert.NotNil(t, c.LivenessProbe.HTTPGet)
	assert.Equal(t, "/health", c.LivenessProbe.HTTPGet.Path)
	assert.Len(t, c.LivenessProbe.HTTPGet.HTTPHeaders, 1)
	assert.NotNil(t, c.Exec)
	assert.NotNil(t, c.TCPSocket)
	assert.Len(t, c.VolumeMounts, 1)
}

func TestProbe(t *testing.T) {
	p := Probe{
		InitialDelaySeconds: 1,
		TimeoutSeconds:      2,
		PeriodSeconds:       3,
		SuccessThreshold:    4,
		FailureThreshold:    5,
	}

	assert.Equal(t, int32(1), p.InitialDelaySeconds)
	assert.Equal(t, int32(2), p.TimeoutSeconds)
	assert.Equal(t, int32(3), p.PeriodSeconds)
	assert.Equal(t, int32(4), p.SuccessThreshold)
	assert.Equal(t, int32(5), p.FailureThreshold)
}
