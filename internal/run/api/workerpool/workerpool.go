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
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
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

	// resp, err := c.ListWorkerPools(ctx, req)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "Unauthenticated") || strings.Contains(err.Error(), "PermissionDenied") {
	// 		return nil, fmt.Errorf("authentication failed: %w. Tip: Ensure your 'gcloud auth application-default login' is valid and has permissions", err)
	// 	}
	// 	return nil, fmt.Errorf("list worker pools failed (parent=%s): %w", req.Parent, err)
	// }

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

		// Map fields
		workerPools = append(workerPools, model.WorkerPool{
			DisplayName: name,
			Name:        resp.Name,
			State:       "...", ///resp.State.String(),
			UpdateTime:  resp.UpdateTime.AsTime(),
			Region:      region,
			Labels:      resp.Annotations,
		})
	}

	return workerPools, nil

	// var workerPools []model.WorkerPool
	// for _, wp := range resp.WorkerPools {
	// 	// Determine Worker Pool Name (Last part of resource name)
	// 	nameParts := strings.Split(wp.Name, "/")
	// 	name := nameParts[len(nameParts)-1]

	// 	// Map fields
	// 	// Note: Detailed config fields are commented out due to proto version mismatch issues.
	// 	// TODO: Map WorkerConfig, NetworkConfig, PrivatePoolVpcConfig when correct proto is available.

	// 	workerPools = append(workerPools, model.WorkerPool{
	// 		Name:        name, // Use simple name
	// 		DisplayName: wp.DisplayName,
	// 		State:       wp.State.String(),
	// 		UpdateTime:  wp.UpdateTime.AsTime(),
	// 		Region:      region,
	// 		Labels:      wp.Annotations,
	// 	})
	// }

	// return workerPools, nil
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
