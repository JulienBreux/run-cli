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
}
