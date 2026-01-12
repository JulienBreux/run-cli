package execution

import (
	"context"
	"fmt"
	"strings"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	model "github.com/JulienBreux/run-cli/internal/run/model/job/execution"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var apiClient Client = &GCPClient{}

// List returns a list of executions for the given project, region and job.
func List(project, region, jobName string) ([]model.Execution, error) {
	ctx := context.Background()
	pbExecutions, err := apiClient.ListExecutions(ctx, project, region, jobName)
	if err != nil {
		return nil, err
	}

	var executions []model.Execution
	for _, resp := range pbExecutions {
		executions = append(executions, mapExecution(resp, region))
	}

	return executions, nil
}

func mapExecution(resp *runpb.Execution, region string) model.Execution {
	var terminalCondition *condition.Condition
	// Cloud Run v2 API usually puts conditions in Conditions list.
	// We can try to find the "Completed" or "Succeeded" condition or check terminal condition logic if exposed directly.
	// runpb.Execution doesn't have a direct TerminalCondition field like Job, but it has Conditions.
	// We usually look for "Completed" type.
	
	// Helper to find latest relevant condition or just map all of them.
	var conditions []*condition.Condition
	for _, c := range resp.Conditions {
		cond := &condition.Condition{
			Type:               c.Type,
			State:              c.State.String(),
			Message:            c.Message,
			LastTransitionTime: c.LastTransitionTime.AsTime(),
		}
		conditions = append(conditions, cond)
		
		// Heuristic: If condition is "Completed", treat as terminal status for summary
		if c.Type == "Completed" {
			terminalCondition = cond
		}
	}

	return model.Execution{
		Name:              resp.Name,
		Job:               resp.Job,
		CreateTime:        resp.CreateTime.AsTime(),
		StartTime:         resp.StartTime.AsTime(),
		CompletionTime:    resp.CompletionTime.AsTime(),
		DeleteTime:        resp.DeleteTime.AsTime(),
		ExpireTime:        resp.ExpireTime.AsTime(),
		TaskCount:         resp.TaskCount,
		SucceededCount:    resp.SucceededCount,
		FailedCount:       resp.FailedCount,
		RunningCount:      resp.RunningCount,
		CancelledCount:    resp.CancelledCount,
		RetriedCount:      resp.RetriedCount,
		LogURI:            resp.LogUri,
		Region:            region,
		Conditions:        conditions,
		TerminalCondition: terminalCondition,
	}
}

// Client defines the interface for Cloud Run Execution operations.
type Client interface {
	ListExecutions(ctx context.Context, project, region, jobName string) ([]*runpb.Execution, error)
}

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct{}

// ListExecutions lists executions for a project, region and job.
func (c *GCPClient) ListExecutions(ctx context.Context, project, region, jobName string) ([]*runpb.Execution, error) {
	creds, err := client.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	cClient, err := createExecutionsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cClient.Close()
	}()

	// Filter by job name
	// The parent is the location. We filter by label or just iterate and filter?
	// API doesn't support generic filtering in ListExecutionsRequest like v1 did?
	// v2 ListExecutionsRequest has no 'Filter' field?
	// Let's check proto.
	// runpb.ListExecutionsRequest: Parent, PageSize, PageToken, ShowDeleted.
	// It does NOT have a Filter field in standard v2 protobufs I see online.
	// Wait, standard Cloud Run API v2 usually returns ALL executions in the location if we don't filter.
	// BUT, executions are children of Jobs? No, they are children of Location.
	// The resource name is projects/*/locations/*/jobs/*/executions/ -> v1
	// v2: projects/*/locations/*/jobs/*/executions
	
	// Wait, `runpb.ListExecutionsRequest` expects Parent = `projects/{project}/locations/{location}/jobs/{job}` OR `projects/{project}/locations/{location}`.
	// If we can pass the job as parent, we get filtered list!
	
	parent := jobName
	if !strings.HasPrefix(jobName, "projects/") {
		parent = fmt.Sprintf("projects/%s/locations/%s/jobs/%s", project, region, jobName)
	}

	req := &runpb.ListExecutionsRequest{
		Parent: parent,
	}

	var executions []*runpb.Execution
	it := cClient.ListExecutions(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			// If passing Job as parent is not supported (it might strictly require Location),
			// we might get an error.
			// Documentation says: "The location and project to list resources on. Format: projects/{project}/locations/{location}. If a parent Job is specified, the format is projects/{project}/locations/{location}/jobs/{job}."
			// So it SHOULD work.
			return nil, client.WrapError(err)
		}
		executions = append(executions, resp)
	}

	return executions, nil
}
