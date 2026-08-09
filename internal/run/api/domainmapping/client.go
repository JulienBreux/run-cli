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

package domainmapping

import (
	"context"
	"fmt"
	"sync"

	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/run/v1"
)

// DomainMappingsClientWrapper defines the interface for the DomainMappings API interactions.
type DomainMappingsClientWrapper interface {
	List(parent string, pageToken string) (*run.ListDomainMappingsResponse, error)
}

// variable for dependency injection
var createClient = func(ctx context.Context, creds *google.Credentials) (DomainMappingsClientWrapper, error) {
	s, err := run.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	return &GCPDomainMappingsClient{service: s}, nil
}

// GCPDomainMappingsClient is the real implementation using the Google Cloud Run API.
type GCPDomainMappingsClient struct {
	service *run.APIService
}

func (c *GCPDomainMappingsClient) List(parent string, pageToken string) (*run.ListDomainMappingsResponse, error) {
	call := c.service.Projects.Locations.Domainmappings.List(parent)
	if pageToken != "" {
		call.Continue(pageToken)
	}
	return call.Do()
}

// Client defines the interface for the DomainMapping API client.
type Client interface {
	ListDomainMappings(ctx context.Context, project, region string) ([]*run.DomainMapping, error)
}

// GCPClient is the Google Cloud Platform implementation of the Client interface.
type GCPClient struct {
	mu     sync.Mutex
	client DomainMappingsClientWrapper
}

// getClient lazily initializes and returns the cached DomainMappingsClientWrapper in a thread-safe manner.
// Using context.Background() ensures the credentials and clients remain valid and are not canceled
// with individual short-lived request contexts.
func (c *GCPClient) getClient(ctx context.Context) (DomainMappingsClientWrapper, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil {
		return c.client, nil
	}

	bgCtx := context.Background()
	creds, err := client.FindDefaultCredentials(bgCtx, run.CloudPlatformScope)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	dmClient, err := createClient(bgCtx, creds)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain mappings client: %w", err)
	}

	c.client = dmClient
	return c.client, nil
}

// ListDomainMappings lists domain mappings for a given project and region.
func (c *GCPClient) ListDomainMappings(ctx context.Context, project, region string) ([]*run.DomainMapping, error) {
	dmClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	parent := fmt.Sprintf("projects/%s/locations/%s", project, region)

	var domainMappings []*run.DomainMapping
	pageToken := ""

	for {
		resp, err := dmClient.List(parent, pageToken)
		if err != nil {
			return nil, fmt.Errorf("failed to list domain mappings: %w", err)
		}

		domainMappings = append(domainMappings, resp.Items...)

		if resp.Metadata != nil {
			pageToken = resp.Metadata.Continue
		} else {
			pageToken = ""
		}

		if pageToken == "" {
			break
		}
	}

	return domainMappings, nil
}
