package workerpool

import (
	"context"
	"fmt"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Interfaces for mocking
type WorkerPoolsClientWrapper interface {
	ListWorkerPools(ctx context.Context, req *runpb.ListWorkerPoolsRequest, opts ...gax.CallOption) WorkerPoolIteratorWrapper
	GetWorkerPool(ctx context.Context, req *runpb.GetWorkerPoolRequest, opts ...gax.CallOption) (*runpb.WorkerPool, error)
	UpdateWorkerPool(ctx context.Context, req *runpb.UpdateWorkerPoolRequest, opts ...gax.CallOption) (UpdateWorkerPoolOperationWrapper, error)
	Close() error
}

type WorkerPoolIteratorWrapper interface {
	Next() (*runpb.WorkerPool, error)
}

type UpdateWorkerPoolOperationWrapper interface {
	Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.WorkerPool, error)
}

// Variables for dependency injection
var createWorkerPoolsClient = func(ctx context.Context, opts ...option.ClientOption) (WorkerPoolsClientWrapper, error) {
	c, err := run.NewWorkerPoolsClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPWorkerPoolsClientWrapper{client: c}, nil
}

// Real implementations
type GCPWorkerPoolsClientWrapper struct {
	client *run.WorkerPoolsClient
}

func (w *GCPWorkerPoolsClientWrapper) ListWorkerPools(ctx context.Context, req *runpb.ListWorkerPoolsRequest, opts ...gax.CallOption) WorkerPoolIteratorWrapper {
	return &GCPWorkerPoolIteratorWrapper{it: w.client.ListWorkerPools(ctx, req, opts...)}
}

func (w *GCPWorkerPoolsClientWrapper) GetWorkerPool(ctx context.Context, req *runpb.GetWorkerPoolRequest, opts ...gax.CallOption) (*runpb.WorkerPool, error) {
	return w.client.GetWorkerPool(ctx, req, opts...)
}

func (w *GCPWorkerPoolsClientWrapper) UpdateWorkerPool(ctx context.Context, req *runpb.UpdateWorkerPoolRequest, opts ...gax.CallOption) (UpdateWorkerPoolOperationWrapper, error) {
	op, err := w.client.UpdateWorkerPool(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPUpdateWorkerPoolOperationWrapper{op: op}, nil
}

func (w *GCPWorkerPoolsClientWrapper) Close() error {
	return w.client.Close()
}

type GCPWorkerPoolIteratorWrapper struct {
	it *run.WorkerPoolIterator
}

func (w *GCPWorkerPoolIteratorWrapper) Next() (*runpb.WorkerPool, error) {
	return w.it.Next()
}

type GCPUpdateWorkerPoolOperationWrapper struct {
	op *run.UpdateWorkerPoolOperation
}

func (w *GCPUpdateWorkerPoolOperationWrapper) Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.WorkerPool, error) {
	return w.op.Wait(ctx, opts...)
}

// Client defines the interface for Cloud Run WorkerPool operations.
type Client interface {
	ListWorkerPools(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error)
	GetWorkerPool(ctx context.Context, name string) (*runpb.WorkerPool, error)
	UpdateWorkerPool(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error)
}

var _ Client = (*GCPClient)(nil)

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct{}

// ListWorkerPools lists worker pools for a project and region.
func (c *GCPClient) ListWorkerPools(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
	creds, err := client.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	cClient, err := createWorkerPoolsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cClient.Close()
	}()

	req := &runpb.ListWorkerPoolsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var workerPools []*runpb.WorkerPool
	it := cClient.ListWorkerPools(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, client.WrapError(err)
		}
		workerPools = append(workerPools, resp)
	}

	return workerPools, nil
}

// GetWorkerPool gets a worker pool.
func (c *GCPClient) GetWorkerPool(ctx context.Context, name string) (*runpb.WorkerPool, error) {
	creds, err := client.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	cClient, err := createWorkerPoolsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cClient.Close()
	}()

	return cClient.GetWorkerPool(ctx, &runpb.GetWorkerPoolRequest{Name: name})
}

// UpdateWorkerPool updates a worker pool.
func (c *GCPClient) UpdateWorkerPool(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error) {
	creds, err := client.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	cClient, err := createWorkerPoolsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = cClient.Close()
	}()

	op, err := cClient.UpdateWorkerPool(ctx, &runpb.UpdateWorkerPoolRequest{WorkerPool: workerPool})
	if err != nil {
		return nil, err
	}

	return op.Wait(ctx)
}
