package workerpool

import (
	"time"
)

// WorkerPool represents a Cloud Build Worker Pool.
type WorkerPool struct {
	Name                 string                `json:"name"` // projects/project-id/locations/location-id/workerPools/workerpool-name
	Uid                  string                `json:"uid"`
	CreateTime           time.Time             `json:"createTime"`
	UpdateTime           time.Time             `json:"updateTime"`
	DeleteTime           time.Time             `json:"deleteTime"`
	State                string                `json:"state"` // PENDING, ACTIVE, DELETE_REQUESTED, DELETING, SUSPENDED
	DisplayName          string                `json:"displayName"`
	Annotations          map[string]string     `json:"annotations"`
	Labels               map[string]string     `json:"labels"`
	Etag                 string                `json:"etag"`
	Project              string                `json:"project"` // The ID of the project in which this worker pool is created.
	Region               string                `json:"region"` // The region of the worker pool.
	WorkerConfig         *WorkerConfig         `json:"workerConfig"`
	NetworkConfig        *NetworkConfig        `json:"networkConfig"`
	PrivatePoolVpcConfig *PrivatePoolVpcConfig `json:"privatePoolVpcConfig"`
	HostIp               string                `json:"hostIp"`
	PublicIp             string                `json:"publicIp"`
}

// WorkerConfig describes the configuration of the workers in a worker pool.
type WorkerConfig struct {
	MachineType  string `json:"machineType"`
	DiskSizeGb   int32  `json:"diskSizeGb"`
	NoExternalIp bool   `json:"noExternalIp"`
	IpCidrRange  string `json:"ipCidrRange"`
}

// NetworkConfig describes the network configuration for a worker pool.
type NetworkConfig struct {
	PeeredNetwork string `json:"peeredNetwork"`
	EgressOption  string `json:"egressOption"` // EgressOptionUnspecified, PrivateEndpoint, NoExternalIP
}

// PrivatePoolVpcConfig describes the VPC configuration for a private worker pool.
type PrivatePoolVpcConfig struct {
	EgressOption string `json:"egressOption"` // PRIVATE_POOL_EGRESS_OPTION_UNSPECIFIED, NO_PUBLIC_EGRESS, PUBLIC_EGRESS
	Subnetwork   string `json:"subnetwork"`
}