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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDomainMapping(t *testing.T) {
	now := time.Now()
	dm := DomainMapping{
		Name:      "example.com",
		RouteName: "service",
		Region:    "us-central1",
		Project:   "project",
		Records: []ResourceRecord{
			{Type: "CNAME", Name: "example.com", RRData: "ghs.googlehosted.com."},
		},
		CreateTime: now,
	}

	assert.Equal(t, "example.com", dm.Name)
	assert.Equal(t, "service", dm.RouteName)
	assert.Equal(t, "us-central1", dm.Region)
	assert.Len(t, dm.Records, 1)
	assert.Equal(t, "CNAME", dm.Records[0].Type)
}
