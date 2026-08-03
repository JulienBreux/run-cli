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
	"context"
	"sync"
	"time"

	"google.golang.org/api/run/v1"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	model "github.com/JulienBreux/run-cli/internal/run/model/domainmapping"
	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
)

var apiClient Client = &GCPClient{}

// List returns a list of domain mappings for the given project and region.
func List(project, region string) ([]model.DomainMapping, error) {
	if region == api_region.ALL {
		return listAllRegions(project)
	}

	ctx := context.Background()
	pbDomainMappings, err := apiClient.ListDomainMappings(ctx, project, region)
	if err != nil {
		return nil, err
	}

	var domainMappings []model.DomainMapping
	for _, resp := range pbDomainMappings {
		domainMappings = append(domainMappings, mapDomainMapping(resp, project, region))
	}

	return domainMappings, nil
}

func listAllRegions(project string) ([]model.DomainMapping, error) {
	regions := api_region.List()
	results := make([][]model.DomainMapping, len(regions))
	var wg sync.WaitGroup

	for i, region := range regions {
		wg.Add(1)
		go func(idx int, r string) {
			defer wg.Done()
			if dms, err := List(project, r); err == nil {
				results[idx] = dms
			}
		}(i, region)
	}

	wg.Wait()

	// Pre-allocate final slice with exact total capacity to eliminate reallocation overhead
	total := 0
	for _, dms := range results {
		total += len(dms)
	}

	domainMappings := make([]model.DomainMapping, 0, total)
	for _, dms := range results {
		domainMappings = append(domainMappings, dms...)
	}

	return domainMappings, nil
}

func mapDomainMapping(resp *run.DomainMapping, project, region string) model.DomainMapping {
	var records []model.ResourceRecord
	if resp.Status != nil {
		for _, r := range resp.Status.ResourceRecords {
			records = append(records, model.ResourceRecord{
				Type:   r.Type,
				Name:   r.Name,
				RRData: r.Rrdata,
			})
		}
	}

	var conditions []*condition.Condition
	if resp.Status != nil {
		for _, c := range resp.Status.Conditions {
			conditions = append(conditions, &condition.Condition{
				Type:    c.Type,
				State:   c.Status,
				Message: c.Message,
				Reason:  c.Reason,
			})
		}
	}
	
	routeName := ""
	if resp.Spec != nil {
		routeName = resp.Spec.RouteName
	}
	
	createTime := time.Time{}
	name := ""
	creator := ""
	if resp.Metadata != nil {
		name = resp.Metadata.Name
		if resp.Metadata.CreationTimestamp != "" {
			t, err := time.Parse(time.RFC3339, resp.Metadata.CreationTimestamp)
			if err == nil {
				createTime = t
			}
		}
		if resp.Metadata.Annotations != nil {
			creator = resp.Metadata.Annotations["serving.knative.dev/creator"]
		}
	}

	return model.DomainMapping{
		Name:       name,
		RouteName:  routeName,
		Region:     region,
		Project:    project,
		Creator:    creator,
		Records:    records,
		CreateTime: createTime,
		Conditions: conditions,
	}
}