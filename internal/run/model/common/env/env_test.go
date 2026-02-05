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

package env

import (
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/model/common/secret"
	"github.com/stretchr/testify/assert"
)

func TestEnvVar(t *testing.T) {
	e := EnvVar{
		Name:  "ENV_VAR",
		Value: "value",
		Source: &EnvVarSource{
			SecretKeyRef: &secret.SecretKeySelector{
				Secret: "secret",
				Key:    "key",
			},
		},
	}

	assert.Equal(t, "ENV_VAR", e.Name)
	assert.Equal(t, "value", e.Value)
	assert.NotNil(t, e.Source)
	assert.Equal(t, "secret", e.Source.SecretKeyRef.Secret)
	assert.Equal(t, "key", e.Source.SecretKeyRef.Key)
}
