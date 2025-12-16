package job

import (
	"context"
	"fmt"
	"strings"
	"sync"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	model "github.com/JulienBreux/run-cli/internal/run/model/job"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// List returns a list of jobs for the given project and region.
// If region is "-", it lists jobs from all supported Cloud Run regions.
func List(project, region string) ([]model.Job, error) {
	if region == "-" {
		return listAllRegions(project)
	}

	ctx := context.Background()

	// Explicitly find default credentials
	creds, err := google.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w. Tip: Try running 'gcloud auth application-default login' to authenticate the Go client", err)
	}

	c, err := run.NewJobsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer c.Close()

	req := &runpb.ListJobsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var jobs []model.Job
	it := c.ListJobs(ctx, req)
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

		jobs = append(jobs, model.Job{
			Name:                   resp.Name,
			LatestCreatedExecution: latestExecution,
			TerminalCondition:      terminalCondition,
			Creator:                resp.Creator,
			Region:                 region, // Assuming single region listing for now
		})
	}

	return jobs, nil
}

func listAllRegions(project string) ([]model.Job, error) {
	var (
		mu   sync.Mutex
		jobs []model.Job
		wg   sync.WaitGroup
	)

	for _, region := range api_region.List() {
		wg.Add(1)
		go func(r string) {
			defer wg.Done()
			// Call List recursively for each region
			// We ignore errors here to allow partial success (e.g. if one region is down or disabled)
			if j, err := List(project, r); err == nil {
				mu.Lock()
				jobs = append(jobs, j...)
				mu.Unlock()
			}
		}(region)
	}

	wg.Wait()
	return jobs, nil
}
