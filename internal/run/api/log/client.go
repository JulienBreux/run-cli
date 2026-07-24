/*
Copyright 2026 Julien Breux

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package log

import (
	"context"
	"fmt"
	"sync"

	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"golang.org/x/oauth2/google"
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

// Global cache for thread-safe GCP Logging Clients to eliminate repetitive
// credentials discovery and client creation overhead (~300ms latency per call).
var (
	logClients     = make(map[string]*GCPClient)
	logClientCreds *google.Credentials
	logClientMu    sync.Mutex
)

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct {
	client LogAdminClientWrapper
}

// NewGCPClient creates a new GCPClient, using cached clients and credentials
// where available to eliminate repetitive discovery and connection overhead (~300ms latency).
func NewGCPClient(ctx context.Context, projectID string) (Client, error) {
	logClientMu.Lock()
	defer logClientMu.Unlock()

	if logClientCreds == nil {
		// Discover credentials using a background context to ensure they remain valid
		// even if the calling request context is cancelled.
		creds, err := client.FindDefaultCredentials(context.Background(), logging.ReadScope)
		if err != nil {
			return nil, fmt.Errorf("failed to find default credentials: %w", err)
		}
		logClientCreds = creds
	}

	if cachedClient, exists := logClients[projectID]; exists {
		return cachedClient, nil
	}

	// Create new client using a background context to ensure connection pool longevity
	c, err := createLogAdminClient(context.Background(), projectID, option.WithCredentials(logClientCreds))
	if err != nil {
		return nil, err
	}

	gcpClient := &GCPClient{client: c}
	logClients[projectID] = gcpClient

	return gcpClient, nil
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
	// Close is a no-op to avoid closing shared connections
	return nil
}

// GCPEntryIterator wraps logadmin.EntryIterator.
type GCPEntryIterator struct {
	it *logadmin.EntryIterator
}

func (it *GCPEntryIterator) Next() (*logging.Entry, error) {
	return it.it.Next()
}
