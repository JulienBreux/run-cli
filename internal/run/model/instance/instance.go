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
	"strings"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/container"
	"github.com/JulienBreux/run-cli/internal/run/model/common/volume"
)

// Instance represents a Google Cloud Run Instance.
type Instance struct {
	Name                       string                 `json:"name"`
	Description                string                 `json:"description,omitempty"`
	UID                        string                 `json:"uid"`
	Generation                 int64                  `json:"generation"`
	Labels                     map[string]string      `json:"labels,omitempty"`
	Annotations                map[string]string      `json:"annotations,omitempty"`
	CreateTime                 time.Time              `json:"createTime"`
	UpdateTime                 time.Time              `json:"updateTime"`
	DeleteTime                 time.Time              `json:"deleteTime,omitempty"`
	ExpireTime                 time.Time              `json:"expireTime,omitempty"`
	Creator                    string                 `json:"creator,omitempty"`
	LastModifier               string                 `json:"lastModifier,omitempty"`
	Client                     string                 `json:"client,omitempty"`
	ClientVersion              string                 `json:"clientVersion,omitempty"`
	LaunchStage                string                 `json:"launchStage,omitempty"`
	ServiceAccount             string                 `json:"serviceAccount,omitempty"`
	Containers                 []*container.Container `json:"containers,omitempty"`
	Volumes                    []*volume.Volume       `json:"volumes,omitempty"`
	EncryptionKey              string                 `json:"encryptionKey,omitempty"`
	RestartPolicy              string                 `json:"restartPolicy,omitempty"`
	ObservedGeneration         int64                  `json:"observedGeneration,omitempty"`
	LogURI                     string                 `json:"logUri,omitempty"`
	TerminalCondition          *condition.Condition   `json:"terminalCondition,omitempty"`
	Conditions                 []*condition.Condition `json:"conditions,omitempty"`
	ContainerStatuses          []*ContainerStatus     `json:"containerStatuses,omitempty"`
	SatisfiesPZS               bool                   `json:"satisfiesPzs,omitempty"`
	URLs                       []string               `json:"urls,omitempty"`
	Reconciling                bool                   `json:"reconciling,omitempty"`
	ETag                       string                 `json:"etag,omitempty"`
	GpuZonalRedundancyDisabled bool                   `json:"gpuZonalRedundancyDisabled,omitempty"`
	InvokerIamDisabled         bool                   `json:"invokerIamDisabled,omitempty"`
	IapEnabled                 bool                   `json:"iapEnabled,omitempty"`
	Region                     string                 `json:"region"`
	Project                    string                 `json:"project"`
}

// ContainerStatus holds container status information for an Instance.
type ContainerStatus struct {
	Name        string `json:"name"`
	ImageDigest string `json:"imageDigest"`
}

// ID returns the short instance ID from the resource name.
func (i *Instance) ID() string {
	if i.Name == "" {
		return ""
	}
	parts := strings.Split(i.Name, "/")
	return parts[len(parts)-1]
}

// Status returns a human-readable status for the instance.
func (i *Instance) Status() string {
	if i.Reconciling {
		return "Reconciling"
	}
	if i.TerminalCondition != nil {
		switch i.TerminalCondition.State {
		case "CONDITION_SUCCEEDED":
			return "Ready"
		case "CONDITION_FAILED":
			return "Failed"
		case "CONDITION_PENDING":
			return "Pending"
		}
	}
	for _, c := range i.Conditions {
		if c.Type == "Ready" {
			switch c.State {
			case "CONDITION_SUCCEEDED":
				return "Ready"
			case "CONDITION_FAILED":
				return "Failed"
			case "CONDITION_PENDING":
				return "Pending"
			}
		}
	}
	return "Unknown"
}

// FormattedURLs returns URLs joined as a comma-separated string.
func (i *Instance) FormattedURLs() string {
	return strings.Join(i.URLs, ", ")
}
