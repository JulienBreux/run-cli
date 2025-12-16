package volume

import (
	"github.com/JulienBreux/run-cli/internal/run/model/common/secret"
)

// Volume represents a volume that is accessible to containers.
type Volume struct {
	Name             string                        `json:"name,omitempty"`
	Secret           *secret.SecretSource          `json:"secret,omitempty"`
	CloudSQLInstance *CloudSQLInstanceVolumeSource `json:"cloudSqlInstance,omitempty"`
	EmptyDir         *EmptyDirVolumeSource         `json:"emptyDir,omitempty"`
	GCS              *GCSVolumeSource              `json:"gcs,omitempty"`
	NFS              *NFSVolumeSource              `json:"nfs,omitempty"`
}

// CloudSQLInstanceVolumeSource represents a volume backed by a Cloud SQL instance.
type CloudSQLInstanceVolumeSource struct {
	Instances []string `json:"instances,omitempty"`
}

// EmptyDirVolumeSource represents an empty directory for a container.
type EmptyDirVolumeSource struct {
	Medium    string `json:"medium,omitempty"`
	SizeLimit string `json:"sizeLimit,omitempty"` // Represented as resource.Quantity in Kubernetes
}

// GCSVolumeSource represents a volume backed by a Google Cloud Storage bucket.
type GCSVolumeSource struct {
	Bucket    string `json:"bucket,omitempty"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
	MountPath string `json:"mountPath,omitempty"`
}

// NFSVolumeSource represents a volume backed by a NFS share.
type NFSVolumeSource struct {
	Server   string `json:"server,omitempty"`
	Path     string `json:"path,omitempty"`
	ReadOnly bool   `json:"readOnly,omitempty"`
}
