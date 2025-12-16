package env

import "github.com/JulienBreux/run-cli/internal/run/model/common/secret"

// EnvVar represents an environment variable present in a container.
type EnvVar struct {
	Name   string        `json:"name"`
	Value  string        `json:"value"`
	Source *EnvVarSource `json:"source"`
}

// EnvVarSource represents a source for the value of an EnvVar.
type EnvVarSource struct {
	SecretKeyRef *secret.SecretKeySelector `json:"secretKeyRef"`
}
