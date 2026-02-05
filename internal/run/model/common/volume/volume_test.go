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

package volume

import (
	"testing"

	"github.com/JulienBreux/run-cli/internal/run/model/common/secret"
	"github.com/stretchr/testify/assert"
)

func TestVolume(t *testing.T) {
	v := Volume{
		Name: "volume",
		Secret: &secret.SecretSource{
			Secret: "secret",
		},
		CloudSQLInstance: &CloudSQLInstanceVolumeSource{
			Instances: []string{"instance"},
		},
		EmptyDir: &EmptyDirVolumeSource{
			Medium:    "Memory",
			SizeLimit: "128Mi",
		},
		GCS: &GCSVolumeSource{
			Bucket: "bucket",
		},
		NFS: &NFSVolumeSource{
			Server: "server",
			Path:   "/path",
		},
	}

	assert.Equal(t, "volume", v.Name)
	assert.NotNil(t, v.Secret)
	assert.Equal(t, "secret", v.Secret.Secret)
	assert.Len(t, v.CloudSQLInstance.Instances, 1)
	assert.Equal(t, "Memory", v.EmptyDir.Medium)
	assert.Equal(t, "bucket", v.GCS.Bucket)
	assert.Equal(t, "server", v.NFS.Server)
}
