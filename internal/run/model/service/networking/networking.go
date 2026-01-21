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

package networking

// Networking represents the networking settings of a service.
type Networking struct {
	Ingress            string     `json:"ingress"`
	DefaultUriDisabled bool       `json:"defaultUriDisabled"`
	IapEnabled         bool       `json:"iapEnabled"`
	VpcAccess          *VpcAccess `json:"vpcAccess,omitempty"`
}

// VpcAccess represents the VPC Access settings.
type VpcAccess struct {
	Connector string `json:"connector,omitempty"`
	Egress    string `json:"egress,omitempty"`
}
