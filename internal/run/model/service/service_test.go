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

package service

import (
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	now := time.Now()
	s := Service{
		Name:       "service",
		URI:        "https://service.run.app",
		CreateTime: now,
		UpdateTime: now,
		Region:     "us-central1",
		Project:    "project",
		Proxy: &ProxyStatus{
			Enabled: true,
			Port:    8080,
			URL:     "http://localhost:8080",
		},
		TerminalCondition: &condition.Condition{
			Type:  "Ready",
			State: "True",
		},
	}

	assert.Equal(t, "service", s.Name)
	assert.Equal(t, "https://service.run.app", s.URI)
	assert.Equal(t, now, s.CreateTime)
	assert.Equal(t, "us-central1", s.Region)
	assert.True(t, s.Proxy.Enabled)
	assert.Equal(t, "Ready", s.TerminalCondition.Type)
}
