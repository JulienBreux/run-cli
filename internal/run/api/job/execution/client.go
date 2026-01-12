package execution

import (
	"context"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/option"
)

// Interfaces for mocking
type ExecutionsClientWrapper interface {
	ListExecutions(ctx context.Context, req *runpb.ListExecutionsRequest, opts ...gax.CallOption) ExecutionIteratorWrapper
	Close() error
}

type ExecutionIteratorWrapper interface {
	Next() (*runpb.Execution, error)
}

// Variables for dependency injection
var createExecutionsClient = func(ctx context.Context, opts ...option.ClientOption) (ExecutionsClientWrapper, error) {
	c, err := run.NewExecutionsClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPExecutionsClientWrapper{client: c}, nil
}

// Real implementations
type GCPExecutionsClientWrapper struct {
	client *run.ExecutionsClient
}

func (w *GCPExecutionsClientWrapper) ListExecutions(ctx context.Context, req *runpb.ListExecutionsRequest, opts ...gax.CallOption) ExecutionIteratorWrapper {
	return &GCPExecutionIteratorWrapper{it: w.client.ListExecutions(ctx, req, opts...)}
}

func (w *GCPExecutionsClientWrapper) Close() error {
	return w.client.Close()
}

type GCPExecutionIteratorWrapper struct {
	it *run.ExecutionIterator
}

func (w *GCPExecutionIteratorWrapper) Next() (*runpb.Execution, error) {
	return w.it.Next()
}
