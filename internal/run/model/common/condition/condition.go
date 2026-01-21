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

package condition

import "time"

// Condition represents a resource's status.
type Condition struct {
	Type                string    `json:"type,omitempty"`
	State               string    `json:"state,omitempty"`
	Message             string    `json:"message,omitempty"`
	LastTransitionTime  time.Time `json:"lastTransitionTime"`
	Severity            string    `json:"severity,omitempty"`
	Reason              string    `json:"reason,omitempty"`
	RevisionGeneration  int64     `json:"revisionGeneration,omitempty"`
	ObservedGeneration  int64     `json:"observedGeneration,omitempty"`
	DomainMappingReason string    `json:"domainMappingReason,omitempty"`
}
