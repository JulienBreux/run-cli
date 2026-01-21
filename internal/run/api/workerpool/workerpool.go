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
	"context"
	"fmt"
	"strings"
	"sync"

	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	model "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/workerpool/scaling"
)

var apiClient Client = &GCPClient{}

// List returns a list of worker pools for the given project and region.
// If region is api_region.ALL, it lists worker pools from all supported Cloud Run regions.
func List(project, region string) ([]model.WorkerPool, error) {
	if region == api_region.ALL {
		return listAllRegions(project)
	}

	ctx := context.Background()
	pbWorkerPools, err := apiClient.ListWorkerPools(ctx, project, region)
	if err != nil {
		return nil, err
	}

	// Fetch worker pools
	var workerPools []model.WorkerPool
	for _, resp := range pbWorkerPools {
		workerPools = append(workerPools, mapWorkerPool(resp, project, region))
	}

	return workerPools, nil
}

func mapWorkerPool(resp *runpb.WorkerPool, project, region string) model.WorkerPool {
	// Determine Service Name (Last part of resource name)
	nameParts := strings.Split(resp.Name, "/")
	name := nameParts[len(nameParts)-1]

	s := model_scaling.Scaling{}
	if resp.Scaling != nil {
		if resp.Scaling.ManualInstanceCount != nil {
			s.ManualInstanceCount = *resp.Scaling.ManualInstanceCount
		}
	}

	// Map fields
	return model.WorkerPool{
		DisplayName:  name,
		Name:         resp.Name,
		State:        "...", ///resp.State.String(),
		UpdateTime:   resp.UpdateTime.AsTime(),
		LastModifier: resp.LastModifier,
		Region:       region,
		Project:      project,
		Scaling:      &s,
		Labels:       resp.Labels,
	}
}

// UpdateScaling updates the scaling settings for a worker pool.
func UpdateScaling(ctx context.Context, project, region, workerPoolName string, instanceCount int32) (*model.WorkerPool, error) {
	fullPoolName := fmt.Sprintf("projects/%s/locations/%s/workerPools/%s", project, region, workerPoolName)
	workerPool, err := apiClient.GetWorkerPool(ctx, fullPoolName)
	if err != nil {
		return nil, fmt.Errorf("failed to get worker pool: %w", err)
	}

	// Update Scaling
	if workerPool.Scaling == nil {
		workerPool.Scaling = &runpb.WorkerPoolScaling{}
	}
	workerPool.Scaling.ManualInstanceCount = &instanceCount

	// Clean up output-only fields
	workerPool.Uid = ""
	workerPool.CreateTime = nil
	workerPool.UpdateTime = nil
	workerPool.DeleteTime = nil
	// workerPool.State is not accessible/exported
	// Keep Etag for concurrency control

	// Update
	resp, err := apiClient.UpdateWorkerPool(ctx, workerPool)
	if err != nil {
		return nil, fmt.Errorf("failed to update worker pool: %w", err)
	}

	wp := mapWorkerPool(resp, project, region)
	return &wp, nil
}

func listAllRegions(project string) ([]model.WorkerPool, error) {
	var (
		mu          sync.Mutex
		workerPools []model.WorkerPool
		wg          sync.WaitGroup
	)

	for _, region := range api_region.List() {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			// Call List recursively for each region
			// We ignore errors here to allow partial success (e.g. if one region is down or disabled)
			if wp, err := List(project, r); err == nil {
				mu.Lock()
				workerPools = append(workerPools, wp...)
				mu.Unlock()
			}
		}(region)
	}

	wg.Wait()
	return workerPools, nil
}

