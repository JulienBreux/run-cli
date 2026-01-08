package job

import (
	"context"
	"fmt"
	"strings"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/googleapis/gax-go/v2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Interfaces for mocking
type JobsClientWrapper interface {
	ListJobs(ctx context.Context, req *runpb.ListJobsRequest, opts ...gax.CallOption) JobIteratorWrapper
	RunJob(ctx context.Context, req *runpb.RunJobRequest, opts ...gax.CallOption) (RunJobOperationWrapper, error)
	Close() error
}

type JobIteratorWrapper interface {
	Next() (*runpb.Job, error)
}

type RunJobOperationWrapper interface {
	Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Execution, error)
}

// Variables for dependency injection
var findDefaultCredentials = google.FindDefaultCredentials
var createJobsClient = func(ctx context.Context, opts ...option.ClientOption) (JobsClientWrapper, error) {
	c, err := run.NewJobsClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPJobsClientWrapper{client: c}, nil
}

// Real implementations
type GCPJobsClientWrapper struct {
	client *run.JobsClient
}

func (w *GCPJobsClientWrapper) ListJobs(ctx context.Context, req *runpb.ListJobsRequest, opts ...gax.CallOption) JobIteratorWrapper {
	return &GCPJobIteratorWrapper{it: w.client.ListJobs(ctx, req, opts...)}
}

func (w *GCPJobsClientWrapper) RunJob(ctx context.Context, req *runpb.RunJobRequest, opts ...gax.CallOption) (RunJobOperationWrapper, error) {
	op, err := w.client.RunJob(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPRunJobOperationWrapper{op: op}, nil
}

func (w *GCPJobsClientWrapper) Close() error {
	return w.client.Close()
}

type GCPJobIteratorWrapper struct {
	it *run.JobIterator
}

func (w *GCPJobIteratorWrapper) Next() (*runpb.Job, error) {
	return w.it.Next()
}

type GCPRunJobOperationWrapper struct {
	op *run.RunJobOperation
}

func (w *GCPRunJobOperationWrapper) Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Execution, error) {
	return w.op.Wait(ctx, opts...)
}

// Client defines the interface for Cloud Run Job operations.
type Client interface {
	ListJobs(ctx context.Context, project, region string) ([]*runpb.Job, error)
	RunJob(ctx context.Context, name string) (*runpb.Execution, error)
}

var _ Client = (*GCPClient)(nil)

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct{}

// ListJobs lists jobs for a project and region.
func (c *GCPClient) ListJobs(ctx context.Context, project, region string) ([]*runpb.Job, error) {
	creds, err := findDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := createJobsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()

	req := &runpb.ListJobsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var jobs []*runpb.Job
	it := client.ListJobs(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			if strings.Contains(err.Error(), "Unauthenticated") || strings.Contains(err.Error(), "PermissionDenied") {
				return nil, fmt.Errorf("authentication failed: %w", err)
			}
			return nil, err
		}
		jobs = append(jobs, resp)
	}

	return jobs, nil
}

// RunJob runs a job.
func (c *GCPClient) RunJob(ctx context.Context, name string) (*runpb.Execution, error) {
	creds, err := findDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	client, err := createJobsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = client.Close()
	}()

	op, err := client.RunJob(ctx, &runpb.RunJobRequest{Name: name})
	if err != nil {
		return nil, err
	}

	return op.Wait(ctx)
}
