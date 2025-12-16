package job

import (
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/container"
	"github.com/JulienBreux/run-cli/internal/run/model/common/volume"
)

// Job represents a Cloud Run job.
type Job struct {
	Name                   string               `json:"name"`
	UID                    string               `json:"uid"`
	Generation             int64                `json:"generation"`
	Labels                 map[string]string    `json:"labels"`
	Annotations            map[string]string    `json:"annotations"`
	CreateTime             time.Time            `json:"createTime"`
	UpdateTime             time.Time            `json:"updateTime"`
	DeleteTime             time.Time            `json:"deleteTime"`
	ExpireTime             time.Time            `json:"expireTime"`
	Creator                string               `json:"creator"`
	LastModifier           string               `json:"lastModifier"`
	Client                 string               `json:"client"`
	ClientVersion          string               `json:"clientVersion"`
	LaunchStage            string               `json:"launchStage"`
	BinaryAuthorization    *BinaryAuthorization `json:"binaryAuthorization"`
	Template               *ExecutionTemplate   `json:"template"`
	ObservedGeneration     int64                `json:"observedGeneration"`
	TerminalCondition      *condition.Condition `json:"terminalCondition"`
	Conditions             []*condition.Condition `json:"conditions"`
	ExecutionCount         int64                `json:"executionCount"`
	LatestCreatedExecution *ExecutionReference  `json:"latestCreatedExecution"`
	Reconciling            bool                 `json:"reconciling"`
	SatisfiesPZS           bool                 `json:"satisfiesPzs"`
	Region                 string               `json:"region"` // New field
}

// ExecutionReference represents a reference to a specific execution.
type ExecutionReference struct {
	Name           string    `json:"name"`
	CreateTime     time.Time `json:"createTime"`
	CompletionTime time.Time `json:"completionTime,omitempty"`
}

// BinaryAuthorization represents the binary authorization configuration.
type BinaryAuthorization struct {
	UseDefault              bool   `json:"useDefault"`
	Policy                  string `json:"policy"`
	BreakglassJustification string `json:"breakglassJustification"`
}

// ExecutionTemplate represents the template used to create executions.
type ExecutionTemplate struct {
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	Parallelism int32             `json:"parallelism"`
	TaskCount   int32             `json:"taskCount"`
	Template    *TaskTemplate     `json:"template"`
}

// TaskTemplate represents the template used to create tasks.
type TaskTemplate struct {
	Containers           []*container.Container `json:"containers"`
	Volumes              []*volume.Volume       `json:"volumes"`
	MaxRetries           int32                  `json:"maxRetries"`
	Timeout              string                 `json:"timeout"`
	ServiceAccount       string                 `json:"serviceAccount"`
	ExecutionEnvironment string                 `json:"executionEnvironment"`
	EncryptionKey        string                 `json:"encryptionKey"`
	VPCAccess            *VPCAccess             `json:"vpcAccess"`
}

// VPCAccess represents the VPC access configuration.
type VPCAccess struct {
	Connector         string              `json:"connector"`
	Egress            string              `json:"egress"`
	NetworkInterfaces []*NetworkInterface `json:"networkInterfaces"`
}

// NetworkInterface represents a network interface.
type NetworkInterface struct {
	Network    string   `json:"network"`
	Subnetwork string   `json:"subnetwork"`
	Tags       []string `json:"tags"`
}