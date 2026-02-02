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

package traffic

const (
	TrafficTargetAllocationTypeLatest = "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST"
)

// TrafficTarget represents a traffic target allocation.
type TrafficTarget struct {
	Type     string `json:"type,omitempty"`
	Revision string `json:"revision,omitempty"`
	Percent  int32  `json:"percent"`
	Tag      string `json:"tag,omitempty"`
}

// TrafficTargetStatus represents the actual traffic allocated to a revision.
type TrafficTargetStatus struct {
	Type     string `json:"type,omitempty"`
	Revision string `json:"revision,omitempty"`
	Percent  int32  `json:"percent,omitempty"`
	Tag      string `json:"tag,omitempty"`
	URI      string `json:"uri,omitempty"`
}
