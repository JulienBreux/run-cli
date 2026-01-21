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

package revision

import (
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/container"
	"github.com/JulienBreux/run-cli/internal/run/model/common/volume"
)

// Revision represents a Cloud Run service revision.
type Revision struct {
	Name                 string                 `json:"name"`
	CreateTime           time.Time              `json:"createTime"`
	UpdateTime           time.Time              `json:"updateTime"`
	Service              string                 `json:"service"`
	Containers           []*container.Container `json:"containers"`
	Volumes              []*volume.Volume       `json:"volumes"`
	ExecutionEnvironment string                 `json:"executionEnvironment"`
	EncryptionKey        string                 `json:"encryptionKey"`
	Reconciling          bool                   `json:"reconciling"`
	Conditions           []*condition.Condition `json:"conditions"`
	ObservedGeneration   int64                  `json:"observedGeneration"`
	LogURI               string                 `json:"logUri"`
	Etag                 string                 `json:"etag"`
	// New fields
	MaxInstanceRequestConcurrency int32         `json:"maxInstanceRequestConcurrency"`
	Timeout                       time.Duration `json:"timeout"`
	CpuIdle                       bool          `json:"cpuIdle"`
	StartupCpuBoost               bool          `json:"startupCpuBoost"`
	Accelerator                   string        `json:"accelerator"`
}
