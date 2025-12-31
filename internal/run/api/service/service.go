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
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/service/scaling"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// List returns a list of services for the given project and region.
// If region is api_region.ALL, it lists services from all supported Cloud Run regions.
func List(project, region string) ([]model.Service, error) {
	if region == api_region.ALL {
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
	defer func() {
		// TODO: Improve error management.
		_ = c.Close()
	}()

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

		services = append(services, model.Service{
			Name:         name,
			URI:          resp.Uri,
			LastModifier: lastModifier,
			UpdateTime:   resp.UpdateTime.AsTime(),
			Region:       region,
			Scaling:      &s,
			Project:      project,
		})
	}

	return services, nil
}

// UpdateScaling updates the scaling settings for a service.
func UpdateScaling(ctx context.Context, project, region, serviceName string, min, max, manual int) (*model.Service, error) {
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	c, err := run.NewServicesClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		// TODO: Improve error management.
		_ = c.Close()
	}()

	// First, get the latest version of the service
	service, err := c.GetService(ctx, &runpb.GetServiceRequest{
		Name: fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, serviceName),
	})
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

		manualInstanceCount := int32(manual)
		service.Scaling.ManualInstanceCount = &manualInstanceCount
	} else {
		service.Scaling.ScalingMode = runpb.ServiceScaling_AUTOMATIC
		service.Scaling.MinInstanceCount = int32(min)
		service.Scaling.MaxInstanceCount = int32(max)
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
	op, err := c.UpdateService(ctx, &runpb.UpdateServiceRequest{
		Service: service,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	resp, err := op.Wait(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for service update: %w", err)
	}

	nameParts := strings.Split(resp.Name, "/")
	name := nameParts[len(nameParts)-1]

	s := model_scaling.Scaling{}
	if resp.Scaling != nil {
		switch resp.Scaling.ScalingMode {
		case runpb.ServiceScaling_AUTOMATIC:
			s.ScalingMode = "AUTOMATIC"
			s.MinInstances = resp.Scaling.MinInstanceCount
			s.MaxInstances = resp.Scaling.MaxInstanceCount
		case runpb.ServiceScaling_MANUAL:
			s.ScalingMode = "MANUAL"
			if resp.Scaling.ManualInstanceCount != nil {
				s.ManualInstanceCount = *resp.Scaling.ManualInstanceCount
			}
		}
	}

	return &model.Service{
		Name:         name,
		URI:          resp.Uri,
		LastModifier: resp.LastModifier,
		UpdateTime:   resp.UpdateTime.AsTime(),
		Region:       region,
		Scaling:      &s,
		Project:      project,
	}, nil
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
