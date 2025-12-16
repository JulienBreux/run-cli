package workerpool

import (
	"context"
	"fmt"
	"strings"
	"sync"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1"
	"cloud.google.com/go/cloudbuild/apiv1/v2/cloudbuildpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	model "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

// List returns a list of worker pools for the given project and region.
// If region is "-", it lists worker pools from all supported Cloud Run regions.
func List(project, region string) ([]model.WorkerPool, error) {
	if region == "-" {
		return listAllRegions(project)
	}

	ctx := context.Background()

	// Explicitly find default credentials
	creds, err := google.FindDefaultCredentials(ctx, cloudbuild.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w. Tip: Try running 'gcloud auth application-default login' to authenticate the Go client", err)
	}

	c, err := cloudbuild.NewClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer c.Close()

	req := &cloudbuildpb.ListWorkerPoolsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var workerPools []model.WorkerPool
	resp, err := c.ListWorkerPools(ctx, req)
	if err != nil {
		if strings.Contains(err.Error(), "Unauthenticated") || strings.Contains(err.Error(), "PermissionDenied") {
			return nil, fmt.Errorf("authentication failed: %w. Tip: Ensure your 'gcloud auth application-default login' is valid and has permissions", err)
		}
		return nil, err
	}

	for _, wp := range resp.WorkerPools {
		// Map fields (simplified for now due to proto version mismatch)
		// TODO: Map WorkerConfig, NetworkConfig, PrivatePoolVpcConfig correctly based on cloudbuildpb version

		workerPools = append(workerPools, model.WorkerPool{
			Name:        wp.Name,
			DisplayName: wp.DisplayName,
			State:       wp.State.String(),
			UpdateTime:  wp.UpdateTime.AsTime(),
			Region:      region,
			Labels:      wp.Annotations,
		})
	}

	return workerPools, nil
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
