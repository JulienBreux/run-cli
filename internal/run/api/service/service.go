package service

import (
	"context"
	"fmt"
	"strings"
	"sync"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	model "github.com/JulienBreux/run-cli/internal/run/model/service"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var allRegions = []string{
	"asia-east1", "asia-northeast1", "asia-northeast2", "asia-northeast3",
	"asia-south1", "asia-southeast1", "asia-southeast2", "australia-southeast1",
	"europe-central2", "europe-north1", "europe-west1", "europe-west2",
	"europe-west3", "europe-west4", "europe-west6", "northamerica-northeast1",
	"southamerica-east1", "us-central1", "us-east1", "us-east4",
	"us-west1", "us-west2", "us-west3", "us-west4",
}

// List returns a list of services for the given project and region.
// If region is "-", it lists services from all supported Cloud Run regions.
func List(project, region string) ([]model.Service, error) {
	if region == "-" {
		return listAllRegions(project)
	}

	ctx := context.Background()

	// Explicitly find default credentials
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w. Tip: Try running 'gcloud auth application-default login' to authenticate the Go client", err)
	}

	c, err := run.NewServicesClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer c.Close()

	req := &runpb.ListServicesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var services []model.Service
	it := c.ListServices(ctx, req)
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

		// Determine Last Modifier
		lastModifier := resp.LastModifier
		if lastModifier == "" {
			lastModifier = resp.Creator
		}

		services = append(services, model.Service{
			Name:         name,
			URI:          resp.Uri,
			LastModifier: lastModifier,
			UpdateTime:   resp.UpdateTime.AsTime(),
			Region:       region, // Add region here
		})
	}

	return services, nil
}

func listAllRegions(project string) ([]model.Service, error) {
	var (
		mu       sync.Mutex
		services []model.Service
		wg       sync.WaitGroup
	)

	for _, region := range api_region.List() {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			// Call List recursively for each region
			// We ignore errors here to allow partial success (e.g. if one region is down or disabled)
			if s, err := List(project, r); err == nil {
				mu.Lock()
				services = append(services, s...)
				mu.Unlock()
			}
		}(region)
	}

	wg.Wait()
	return services, nil
}
