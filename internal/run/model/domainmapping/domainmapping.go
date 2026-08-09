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

package domainmapping

import (
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
)

// DomainMapping represents a Cloud Run domain mapping.
type DomainMapping struct {
	Name       string                 `json:"name"`
	RouteName  string                 `json:"routeName"`
	Region     string                 `json:"region"`
	Project    string                 `json:"project"`
	Creator    string                 `json:"creator"`
	Records    []ResourceRecord       `json:"records"`
	CreateTime time.Time              `json:"createTime"`
	UpdateTime time.Time              `json:"updateTime"`
	Conditions []*condition.Condition `json:"conditions,omitempty"`
}

// ResourceRecord represents a DNS resource record.
type ResourceRecord struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	RRData string `json:"rrData"`
}
