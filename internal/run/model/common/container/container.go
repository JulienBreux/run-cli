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
