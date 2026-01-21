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
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/service/networking"
	"github.com/JulienBreux/run-cli/internal/run/model/service/scaling"
	"github.com/JulienBreux/run-cli/internal/run/model/service/security"
	"github.com/JulienBreux/run-cli/internal/run/model/service/traffic"
)

// Service represents a Cloud Run service.
type Service struct {
	Name                  string                         `json:"name"`
	Description           string                         `json:"description,omitempty"`
	URI                   string                         `json:"uri"`
	CreateTime            time.Time                      `json:"createTime"`
	UpdateTime            time.Time                      `json:"updateTime"`
	DeleteTime            time.Time                      `json:"deleteTime"`
	ExpireTime            time.Time                      `json:"expireTime"`
	Creator               string                         `json:"creator,omitempty"`
	LastModifier          string                         `json:"lastModifier,omitempty"`
	Reconciling           bool                           `json:"reconciling"`
	TrafficStatuses       []*traffic.TrafficTargetStatus `json:"trafficStatuses,omitempty"`
	LatestReadyRevision   string                         `json:"latestReadyRevision,omitempty"`
	LatestCreatedRevision string                         `json:"latestCreatedRevision,omitempty"`
	TerminalCondition     *condition.Condition           `json:"terminalCondition,omitempty"`
	Conditions            []*condition.Condition         `json:"conditions,omitempty"`
	Scaling               *scaling.Scaling               `json:"scaling,omitempty"`
	Networking            *networking.Networking         `json:"networking,omitempty"`
	Security              *security.Security             `json:"security,omitempty"` // New field
	Etag                  string                         `json:"etag,omitempty"`
	Region                string                         `json:"region"` // New field
	Project               string                         `json:"project"`
	Proxy                 *ProxyStatus                   `json:"proxy,omitempty"`
}

// ProxyStatus represents the status of a local proxy for the service.
type ProxyStatus struct {
	Enabled bool
	Port    int
	URL     string
}
