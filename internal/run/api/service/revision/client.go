package revision

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
type RevisionsClientWrapper interface {
	ListRevisions(ctx context.Context, req *runpb.ListRevisionsRequest, opts ...gax.CallOption) RevisionIteratorWrapper
	Close() error
}

type RevisionIteratorWrapper interface {
	Next() (*runpb.Revision, error)
}

// Variables for dependency injection
var createRevisionsClient = func(ctx context.Context, opts ...option.ClientOption) (RevisionsClientWrapper, error) {
	c, err := run.NewRevisionsClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPRevisionsClientWrapper{client: c}, nil
}

// Real implementations
type GCPRevisionsClientWrapper struct {
	client *run.RevisionsClient
}

func (w *GCPRevisionsClientWrapper) ListRevisions(ctx context.Context, req *runpb.ListRevisionsRequest, opts ...gax.CallOption) RevisionIteratorWrapper {
	return &GCPRevisionIteratorWrapper{it: w.client.ListRevisions(ctx, req, opts...)}
}

func (w *GCPRevisionsClientWrapper) Close() error {
	return w.client.Close()
}

type GCPRevisionIteratorWrapper struct {
	it *run.RevisionIterator
}

func (w *GCPRevisionIteratorWrapper) Next() (*runpb.Revision, error) {
	return w.it.Next()
}

// Client defines the interface for Cloud Run Revision operations.
type Client interface {
	ListRevisions(ctx context.Context, project, region, service string) ([]*runpb.Revision, error)
}

var apiClient Client = &GCPClient{}

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct{}

// ListRevisions lists revisions for a service.
func (c *GCPClient) ListRevisions(ctx context.Context, project, region, service string) ([]*runpb.Revision, error) {
	creds, err := client.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	cClient, err := createRevisionsClient(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	defer func() { _ = cClient.Close() }()

	req := &runpb.ListRevisionsRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s/services/%s", project, region, service),
	}

	var revisions []*runpb.Revision
	it := cClient.ListRevisions(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		revisions = append(revisions, resp)
	}

	return revisions, nil
}