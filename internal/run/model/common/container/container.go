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
	"github.com/JulienBreux/run-cli/internal/run/model/common/env"
	"github.com/JulienBreux/run-cli/internal/run/model/common/resources"
)

// Container represents a single container that is starting and running in a revision.
type Container struct {
	Name                  string               `json:"name,omitempty"`
	Image                 string               `json:"image"`
	Command               []string             `json:"command,omitempty"`
	Args                  []string             `json:"args,omitempty"`
	Env                   []*env.EnvVar        `json:"env,omitempty"`
	Resources             *resources.Resources `json:"resources,omitempty"`
	VolumeMounts          []*VolumeMount       `json:"volumeMounts,omitempty"`
	Ports                 []*Port              `json:"ports,omitempty"`
	LivenessProbe         *Probe               `json:"livenessProbe,omitempty"`
	StartupProbe          *Probe               `json:"startupProbe,omitempty"`
	WorkingDirectory      string               `json:"workingDirectory,omitempty"`
	GRPCtimePeriodSeconds int64                `json:"grpcTimePeriodSeconds,omitempty"`
	DependsOn             []string             `json:"dependsOn,omitempty"`
}
