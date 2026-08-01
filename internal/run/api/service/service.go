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

package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	model "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_networking "github.com/JulienBreux/run-cli/internal/run/model/service/networking"
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/service/scaling"
	model_security "github.com/JulienBreux/run-cli/internal/run/model/service/security"
	model_traffic "github.com/JulienBreux/run-cli/internal/run/model/service/traffic"
)

var apiClient Client = &GCPClient{}

// List returns a list of services for the given project and region.
// If region is api_region.ALL, it lists services from all supported Cloud Run regions.
func List(project, region string) ([]model.Service, error) {
	if region == api_region.ALL {
		return listAllRegions(project)
	}

	ctx := context.Background()
	pbServices, err := apiClient.ListServices(ctx, project, region)
	if err != nil {
		return nil, err
	}

	var services []model.Service
	for _, resp := range pbServices {
		services = append(services, mapService(resp, project, region))
	}

	return services, nil
}

func mapService(resp *runpb.Service, project, region string) model.Service {
	// Determine Service Name (Last part of resource name)
	nameParts := strings.Split(resp.Name, "/")
	name := nameParts[len(nameParts)-1]

	// Determine Last Modifier
	lastModifier := resp.LastModifier
	if lastModifier == "" {
		lastModifier = resp.Creator
	}

	var trafficStatuses []*model_traffic.TrafficTargetStatus
	for _, ts := range resp.TrafficStatuses {
		revName := ts.Revision
		if strings.Contains(revName, "/") {
			parts := strings.Split(revName, "/")
			revName = parts[len(parts)-1]
		}
		trafficStatuses = append(trafficStatuses, &model_traffic.TrafficTargetStatus{
			Type:     ts.Type.String(),
			Revision: revName,
			Percent:  ts.Percent,
			Tag:      ts.Tag,
			URI:      ts.Uri,
		})
	}

	s := model_scaling.Scaling{
		ScalingMode: "AUTOMATIC",
	}
	if resp.Scaling != nil {
		switch resp.Scaling.ScalingMode {
		case runpb.ServiceScaling_AUTOMATIC:
			s.ScalingMode = "AUTOMATIC"
			s.MinInstances = resp.Scaling.MinInstanceCount
		case runpb.ServiceScaling_MANUAL:
			s.ScalingMode = "MANUAL"
			if resp.Scaling.ManualInstanceCount != nil {
				s.ManualInstanceCount = *resp.Scaling.ManualInstanceCount
			}
		}
		s.MinInstances = resp.Scaling.MinInstanceCount
		s.MaxInstances = resp.Scaling.MaxInstanceCount
	}

	latestReadyRevision := resp.LatestReadyRevision
	if strings.Contains(latestReadyRevision, "/") {
		parts := strings.Split(latestReadyRevision, "/")
		latestReadyRevision = parts[len(parts)-1]
	}

	latestCreatedRevision := resp.LatestCreatedRevision
	if strings.Contains(latestCreatedRevision, "/") {
		parts := strings.Split(latestCreatedRevision, "/")
		latestCreatedRevision = parts[len(parts)-1]
	}

	n := model_networking.Networking{
		Ingress:            resp.Ingress.String(),
		DefaultUriDisabled: resp.DefaultUriDisabled,
		IapEnabled:         resp.IapEnabled,
	}
	if resp.Template != nil && resp.Template.VpcAccess != nil {
		n.VpcAccess = &model_networking.VpcAccess{
			Connector: resp.Template.VpcAccess.Connector,
			Egress:    resp.Template.VpcAccess.Egress.String(),
		}
	}

	sec := model_security.Security{
		InvokerIAMDisabled: resp.InvokerIamDisabled,
	}
	if resp.BinaryAuthorization != nil {
		sec.BinaryAuthorization = resp.BinaryAuthorization.GetPolicy()
		if resp.BinaryAuthorization.GetUseDefault() {
			sec.BinaryAuthorization = "default"
		}
		sec.BreakglassJustification = resp.BinaryAuthorization.BreakglassJustification
	}
	if resp.Template != nil {
		sec.ServiceAccount = resp.Template.ServiceAccount
		sec.EncryptionKey = resp.Template.EncryptionKey
	}

	return model.Service{
		Name:                  name,
		URI:                   resp.Uri,
		LastModifier:          lastModifier,
		UpdateTime:            resp.UpdateTime.AsTime(),
		Region:                region,
		Scaling:               &s,
		Project:               project,
		TrafficStatuses:       trafficStatuses,
		LatestReadyRevision:   latestReadyRevision,
		LatestCreatedRevision: latestCreatedRevision,
		Networking:            &n,
		Security:              &sec,
	}
}

// UpdateScaling updates the scaling settings for a service.
func UpdateScaling(ctx context.Context, project, region, serviceName string, min, max, manual int32) (*model.Service, error) {
	fullServiceName := fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, serviceName)
	
	service, err := apiClient.GetService(ctx, fullServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	// Modify the scaling settings
	if service.Scaling == nil {
		service.Scaling = &runpb.ServiceScaling{}
	}

	if manual > 0 {
		service.Scaling.ScalingMode = runpb.ServiceScaling_MANUAL
		service.Scaling.MinInstanceCount = 0
		service.Scaling.MaxInstanceCount = 0

		manualInstanceCount := manual
		service.Scaling.ManualInstanceCount = &manualInstanceCount
	} else {
		service.Scaling.ScalingMode = runpb.ServiceScaling_AUTOMATIC
		service.Scaling.MinInstanceCount = min
		service.Scaling.MaxInstanceCount = max
		service.Scaling.ManualInstanceCount = nil
	}

	// Clean up output-only fields before update
	service.Uid = ""
	service.Generation = 0
	service.CreateTime = nil
	service.UpdateTime = nil
	service.DeleteTime = nil
	service.ExpireTime = nil
	service.Creator = ""
	service.LastModifier = ""
	service.Reconciling = false
	service.ObservedGeneration = 0
	service.TerminalCondition = nil
	service.Conditions = nil
	// Keep Etag for concurrency control

	// Update the service
	resp, err := apiClient.UpdateService(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	s := mapService(resp, project, region)
	return &s, nil
}

// UpdateAuthentication updates the authentication settings for a service.
func UpdateAuthentication(ctx context.Context, project, region, serviceName string, allowUnauthenticated bool) (*model.Service, error) {
	fullServiceName := fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, serviceName)

	service, err := apiClient.GetService(ctx, fullServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	// Modify the authentication settings
	service.InvokerIamDisabled = allowUnauthenticated

	// Clean up output-only fields before update
	service.Uid = ""
	service.Generation = 0
	service.CreateTime = nil
	service.UpdateTime = nil
	service.DeleteTime = nil
	service.ExpireTime = nil
	service.Creator = ""
	service.LastModifier = ""
	service.Reconciling = false
	service.ObservedGeneration = 0
	service.TerminalCondition = nil
	service.Conditions = nil
	// Keep Etag for concurrency control

	// Update the service
	resp, err := apiClient.UpdateService(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	s := mapService(resp, project, region)
	return &s, nil
}

// UpdateTraffic updates the traffic split for a service.
func UpdateTraffic(ctx context.Context, project, region, serviceName string, targets []model_traffic.TrafficTarget) (*model.Service, error) {
	fullServiceName := fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, serviceName)

	service, err := apiClient.GetService(ctx, fullServiceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	var pbTargets []*runpb.TrafficTarget
	for _, t := range targets {
		pbTarget := &runpb.TrafficTarget{
			Percent:  t.Percent,
			Revision: t.Revision,
			Tag:      t.Tag,
		}
		if t.Type != "" {
			if val, ok := runpb.TrafficTargetAllocationType_value[t.Type]; ok {
				pbTarget.Type = runpb.TrafficTargetAllocationType(val)
			}
		} else if t.Revision == "" {
			pbTarget.Type = runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST
		} else {
			pbTarget.Type = runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_REVISION
		}
		pbTargets = append(pbTargets, pbTarget)
	}

	service.Traffic = pbTargets

	// Clean up output-only fields before update
	service.Uid = ""
	service.Generation = 0
	service.CreateTime = nil
	service.UpdateTime = nil
	service.DeleteTime = nil
	service.ExpireTime = nil
	service.Creator = ""
	service.LastModifier = ""
	service.Reconciling = false
	service.ObservedGeneration = 0
	service.TerminalCondition = nil
	service.Conditions = nil

	// Update the service
	resp, err := apiClient.UpdateService(ctx, service)
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	s := mapService(resp, project, region)
	return &s, nil
}

func listAllRegions(project string) ([]model.Service, error) {
	regions := api_region.List()
	results := make([][]model.Service, len(regions))
	var wg sync.WaitGroup

	for i, region := range regions {
		wg.Add(1)
		go func(idx int, r string) {
			defer wg.Done()
			// Call List recursively for each region
			// We ignore errors here to allow partial success (e.g. if one region is down or disabled)
			if s, err := List(project, r); err == nil {
				results[idx] = s
			}
		}(i, region)
	}

	wg.Wait()

	// Pre-calculate exact capacity needed to eliminate dynamic slice growth
	var totalSize int
	for _, s := range results {
		totalSize += len(s)
	}

	// Lock-free flat pre-allocated slice
	services := make([]model.Service, 0, totalSize)
	for _, s := range results {
		services = append(services, s...)
	}

	return services, nil
}
