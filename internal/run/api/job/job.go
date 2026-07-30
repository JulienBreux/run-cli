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

package job

import (
	"context"
	"sync"

	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	model "github.com/JulienBreux/run-cli/internal/run/model/job"
)

var apiClient Client = &GCPClient{}

// List returns a list of jobs for the given project and region.
// If region is api_region.ALL, it lists jobs from all supported Cloud Run regions.
func List(project, region string) ([]model.Job, error) {
	if region == api_region.ALL {
		return listAllRegions(project)
	}

	ctx := context.Background()
	pbJobs, err := apiClient.ListJobs(ctx, project, region)
	if err != nil {
		return nil, err
	}

	var jobs []model.Job
	for _, resp := range pbJobs {
		jobs = append(jobs, mapJob(resp, region))
	}

	return jobs, nil
}

func mapJob(resp *runpb.Job, region string) model.Job {
	// Map LatestCreatedExecution
	var latestExecution *model.ExecutionReference
	if resp.LatestCreatedExecution != nil {
		latestExecution = &model.ExecutionReference{
			Name:       resp.LatestCreatedExecution.Name,
			CreateTime: resp.LatestCreatedExecution.CreateTime.AsTime(),
		}
	}

	// Map TerminalCondition
	var terminalCondition *condition.Condition
	if resp.TerminalCondition != nil {
		terminalCondition = &condition.Condition{
			State:              resp.TerminalCondition.State.String(),
			Message:            resp.TerminalCondition.Message,
			LastTransitionTime: resp.TerminalCondition.LastTransitionTime.AsTime(),
		}
	}

	return model.Job{
		Name:                   resp.Name,
		LatestCreatedExecution: latestExecution,
		TerminalCondition:      terminalCondition,
		Creator:                resp.Creator,
		Region:                 region,
	}
}

func listAllRegions(project string) ([]model.Job, error) {
	regions := api_region.List()
	results := make([][]model.Job, len(regions))
	var wg sync.WaitGroup

	for i, region := range regions {
		wg.Add(1)
		go func(idx int, r string) {
			defer wg.Done()
			// Call List recursively for each region
			// We ignore errors here to allow partial success (e.g. if one region is down or disabled)
			if j, err := List(project, r); err == nil {
				results[idx] = j
			}
		}(i, region)
	}

	wg.Wait()

	// Calculate exact total capacity needed to allocate final slice exactly once
	// to completely eliminate lock contention (sync.Mutex) and multiple slice reallocation overheads.
	total := 0
	for _, j := range results {
		total += len(j)
	}

	jobs := make([]model.Job, 0, total)
	for _, j := range results {
		jobs = append(jobs, j...)
	}

	return jobs, nil
}

// Execute executes a Cloud Run job.
func Execute(project, region, jobName string) (*runpb.Execution, error) {
	ctx := context.Background()
	// Name format: projects/{project}/locations/{region}/jobs/{job}
	// The client's RunJob usually takes the full name resource string or the method handles it.
	// My GCPClient.RunJob takes just 'name'.
	// I should construct the full name here before passing it.
	fullName := "projects/" + project + "/locations/" + region + "/jobs/" + jobName
	return apiClient.RunJob(ctx, fullName)
}
