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

package workerpool

import (
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/workerpool/scaling"
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
	LastModifier         string                `json:"lastModifier,omitempty"`
	Etag                 string                `json:"etag"`
	Project              string                `json:"project"` // The ID of the project in which this worker pool is created.
	Region               string                `json:"region"`  // The region of the worker pool.
	WorkerConfig         *WorkerConfig         `json:"workerConfig"`
	NetworkConfig        *NetworkConfig        `json:"networkConfig"`
	PrivatePoolVpcConfig *PrivatePoolVpcConfig `json:"privatePoolVpcConfig"`
	HostIp               string                `json:"hostIp"`
	PublicIp             string                `json:"publicIp"`
	Scaling              *scaling.Scaling      `json:"scaling,omitempty"`
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
