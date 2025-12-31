package workerpool

import (
	"context"
	"fmt"
	"strings"
	"sync"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	model "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/workerpool/scaling"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// List returns a list of worker pools for the given project and region.
// If region is api_region.ALL, it lists worker pools from all supported Cloud Run regions.
func List(project, region string) ([]model.WorkerPool, error) {
	if region == api_region.ALL {
		return listAllRegions(project)
	}

	ctx := context.Background()

	// Explicitly find default credentials
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w. Tip: Try running 'gcloud auth application-default login' to authenticate the Go client", err)
	}

	c, err := run.NewWorkerPoolsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = c.Close()
	}()

	req := &runpb.ListWorkerPoolsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	// Fetch worker pools
	var workerPools []model.WorkerPool
	it := c.ListWorkerPools(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), "Unauthenticated") || strings.Contains(err.Error(), "PermissionDenied") {
				return nil, fmt.Errorf("authentication failed: %w. Tip: Ensure your 'gcloud auth application-default login' is valid and has permissions", err)
			}
			return nil, err
		}

		// Determine Service Name (Last part of resource name)
		nameParts := strings.Split(resp.Name, "/")
		name := nameParts[len(nameParts)-1]

		s := model_scaling.Scaling{}
		if resp.Scaling != nil {
			s.ManualInstanceCount = *resp.Scaling.ManualInstanceCount

		}

		// Map fields
		workerPools = append(workerPools, model.WorkerPool{
			DisplayName:  name,
			Name:         resp.Name,
			State:        "...", ///resp.State.String(),
			UpdateTime:   resp.UpdateTime.AsTime(),
			LastModifier: resp.LastModifier,
			Region:       region,
			Project:      project,
			Scaling:      &s,
			Labels:       resp.Labels,
		})
	}

	return workerPools, nil
}

// UpdateScaling updates the scaling settings for a worker pool.
func UpdateScaling(ctx context.Context, project, region, workerPoolName string, instanceCount int) (*model.WorkerPool, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	c, err := run.NewWorkerPoolsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = c.Close()
	}()

	// Get the worker pool
	fullPoolName := fmt.Sprintf("projects/%s/locations/%s/workerPools/%s", project, region, workerPoolName)
	workerPool, err := c.GetWorkerPool(ctx, &runpb.GetWorkerPoolRequest{
		Name: fullPoolName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get worker pool: %w", err)
	}

	// Update Scaling
	if workerPool.Scaling == nil {
		workerPool.Scaling = &runpb.WorkerPoolScaling{}
	}
	manualCount := int32(instanceCount)
	workerPool.Scaling.ManualInstanceCount = &manualCount

	// Clean up output-only fields
	workerPool.Uid = ""
	workerPool.CreateTime = nil
	workerPool.UpdateTime = nil
	workerPool.DeleteTime = nil
	// workerPool.State is not accessible/exported
	// Keep Etag for concurrency control

	// Update
	op, err := c.UpdateWorkerPool(ctx, &runpb.UpdateWorkerPoolRequest{
		WorkerPool: workerPool,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update worker pool: %w", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for worker pool update: %w", err)
	}

	// Map response to model
	nameParts := strings.Split(resp.Name, "/")
	name := nameParts[len(nameParts)-1]

	s := model_scaling.Scaling{}
	if resp.Scaling != nil && resp.Scaling.ManualInstanceCount != nil {
		s.ManualInstanceCount = *resp.Scaling.ManualInstanceCount
	}

	return &model.WorkerPool{
		DisplayName:  name,
		Name:         resp.Name,
		State:        "...", // resp.State.String(),
		UpdateTime:   resp.UpdateTime.AsTime(),
		LastModifier: resp.LastModifier,
		Region:       region,
		Scaling:      &s,
		Labels:       resp.Labels,
	}, nil
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
