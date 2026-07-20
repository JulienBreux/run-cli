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

package project

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"

	resourcemanager "cloud.google.com/go/resourcemanager/apiv3"
	resourcemanagerpb "cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	model "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Interfaces for mocking
type ProjectsClientWrapper interface {
	SearchProjects(ctx context.Context, req *resourcemanagerpb.SearchProjectsRequest, opts ...gax.CallOption) ProjectIteratorWrapper
	Close() error
}

type ProjectIteratorWrapper interface {
	Next() (*resourcemanagerpb.Project, error)
}

// Variables for dependency injection
var createProjectsClient = func(ctx context.Context, opts ...option.ClientOption) (ProjectsClientWrapper, error) {
	c, err := resourcemanager.NewProjectsClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPProjectsClientWrapper{client: c}, nil
}

// Real implementations
type GCPProjectsClientWrapper struct {
	client *resourcemanager.ProjectsClient
}

func (w *GCPProjectsClientWrapper) SearchProjects(ctx context.Context, req *resourcemanagerpb.SearchProjectsRequest, opts ...gax.CallOption) ProjectIteratorWrapper {
	return &GCPProjectIteratorWrapper{it: w.client.SearchProjects(ctx, req, opts...)}
}

func (w *GCPProjectsClientWrapper) Close() error {
	return w.client.Close()
}

type GCPProjectIteratorWrapper struct {
	it *resourcemanager.ProjectIterator
}

func (w *GCPProjectIteratorWrapper) Next() (*resourcemanagerpb.Project, error) {
	return w.it.Next()
}

// Client defines the interface for Cloud Resource Manager operations.
type Client interface {
	ListProjects(ctx context.Context) ([]model.Project, error)
}

var _ Client = (*GCPClient)(nil)

// GCPClient is the Google Cloud Platform implementation of Client.
// It caches the ProjectsClientWrapper to avoid redundant credential lookups
// and gRPC connection overhead in high-frequency TUI operations.
type GCPClient struct {
	mu           sync.Mutex
	cachedClient ProjectsClientWrapper
}

// getOrCreateClient returns a cached client or creates a new one if it doesn't exist.
// This reduces latency by avoiding repeated authentication and connection overhead.
func (c *GCPClient) getOrCreateClient(ctx context.Context) (ProjectsClientWrapper, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cachedClient != nil {
		return c.cachedClient, nil
	}

	bgCtx := context.Background()
	creds, err := client.FindDefaultCredentials(bgCtx, resourcemanager.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w. Tip: Try running 'gcloud auth application-default login' to authenticate the Go client", err)
	}

	cClient, err := createProjectsClient(bgCtx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}

	c.cachedClient = cClient
	return c.cachedClient, nil
}

// ListProjects lists projects for the current user.
func (c *GCPClient) ListProjects(ctx context.Context) ([]model.Project, error) {
	cClient, err := c.getOrCreateClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &resourcemanagerpb.SearchProjectsRequest{
		// Query: "", // Empty query lists all projects
	}

	var projects []model.Project
	it := cClient.SearchProjects(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, client.WrapError(err)
		}

		projects = append(projects, mapProject(resp))
	}

	return projects, nil
}

func mapProject(resp *resourcemanagerpb.Project) model.Project {
	// Parse Project Number from Name "projects/123456"
	parts := strings.Split(resp.Name, "/")
	numberStr := parts[len(parts)-1]
	number, _ := strconv.Atoi(numberStr)

	return model.Project{
		Name:   resp.ProjectId,
		Number: number,
	}
}
