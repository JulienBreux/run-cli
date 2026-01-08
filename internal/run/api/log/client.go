package log

import (
	"context"
	"fmt"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"google.golang.org/api/option"
)

// Client defines the interface for Logging operations.
type Client interface {
	Entries(ctx context.Context, opts ...interface{}) EntryIterator
	Close() error
}

// EntryIterator defines the interface for iterating over log entries.
type EntryIterator interface {
	Next() (*logging.Entry, error)
}

// ClientFactory is a function that returns a Client.
type ClientFactory func(ctx context.Context, projectID string) (Client, error)

var clientFactory ClientFactory = NewGCPClient

// Interfaces for mocking
type LogAdminClientWrapper interface {
	Entries(ctx context.Context, opts ...logadmin.EntriesOption) EntryIterator
	Close() error
}

// Variables for dependency injection
var createLogAdminClient = func(ctx context.Context, projectID string, opts ...option.ClientOption) (LogAdminClientWrapper, error) {
	c, err := logadmin.NewClient(ctx, projectID, opts...)
	if err != nil {
		return nil, err
	}
	return &RealLogAdminClient{client: c}, nil
}

// RealLogAdminClient wraps logadmin.Client.
type RealLogAdminClient struct {
	client *logadmin.Client
}

func (w *RealLogAdminClient) Entries(ctx context.Context, opts ...logadmin.EntriesOption) EntryIterator {
	return &GCPEntryIterator{it: w.client.Entries(ctx, opts...)}
}

func (w *RealLogAdminClient) Close() error {
	return w.client.Close()
}

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct {
	client LogAdminClientWrapper
}

// NewGCPClient creates a new GCPClient.
func NewGCPClient(ctx context.Context, projectID string) (Client, error) {
	creds, err := client.FindDefaultCredentials(ctx, logging.ReadScope)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	c, err := createLogAdminClient(ctx, projectID, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	return &GCPClient{client: c}, nil
}

func (c *GCPClient) Entries(ctx context.Context, opts ...interface{}) EntryIterator {
	var logOpts []logadmin.EntriesOption
	for _, o := range opts {
		if opt, ok := o.(logadmin.EntriesOption); ok {
			logOpts = append(logOpts, opt)
		}
	}
	return c.client.Entries(ctx, logOpts...)
}

func (c *GCPClient) Close() error {
	return c.client.Close()
}

// GCPEntryIterator wraps logadmin.EntryIterator.
type GCPEntryIterator struct {
	it *logadmin.EntryIterator
}

func (it *GCPEntryIterator) Next() (*logging.Entry, error) {
	return it.it.Next()
}
