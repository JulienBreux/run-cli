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

package service

import (
	"context"
	"fmt"
	"sync"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"github.com/googleapis/gax-go/v2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// Interfaces for mocking
type ServicesClientWrapper interface {
	ListServices(ctx context.Context, req *runpb.ListServicesRequest, opts ...gax.CallOption) ServiceIteratorWrapper
	GetService(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error)
	UpdateService(ctx context.Context, req *runpb.UpdateServiceRequest, opts ...gax.CallOption) (UpdateServiceOperationWrapper, error)
	Close() error
}

type ServiceIteratorWrapper interface {
	Next() (*runpb.Service, error)
}

type UpdateServiceOperationWrapper interface {
	Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Service, error)
}

// Variables for dependency injection
var createServicesClient = func(ctx context.Context, opts ...option.ClientOption) (ServicesClientWrapper, error) {
	c, err := run.NewServicesClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPServicesClientWrapper{client: c}, nil
}

// Real implementations
type GCPServicesClientWrapper struct {
	client *run.ServicesClient
}

func (w *GCPServicesClientWrapper) ListServices(ctx context.Context, req *runpb.ListServicesRequest, opts ...gax.CallOption) ServiceIteratorWrapper {
	return &GCPServiceIteratorWrapper{it: w.client.ListServices(ctx, req, opts...)}
}

func (w *GCPServicesClientWrapper) GetService(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
	return w.client.GetService(ctx, req, opts...)
}

func (w *GCPServicesClientWrapper) UpdateService(ctx context.Context, req *runpb.UpdateServiceRequest, opts ...gax.CallOption) (UpdateServiceOperationWrapper, error) {
	op, err := w.client.UpdateService(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPUpdateServiceOperationWrapper{op: op}, nil
}

func (w *GCPServicesClientWrapper) Close() error {
	return w.client.Close()
}

type GCPServiceIteratorWrapper struct {
	it *run.ServiceIterator
}

func (w *GCPServiceIteratorWrapper) Next() (*runpb.Service, error) {
	return w.it.Next()
}

type GCPUpdateServiceOperationWrapper struct {
	op *run.UpdateServiceOperation
}

func (w *GCPUpdateServiceOperationWrapper) Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Service, error) {
	return w.op.Wait(ctx, opts...)
}

// Client defines the interface for Cloud Run Service operations.
type Client interface {
	ListServices(ctx context.Context, project, region string) ([]*runpb.Service, error)
	GetService(ctx context.Context, name string) (*runpb.Service, error)
	UpdateService(ctx context.Context, service *runpb.Service) (*runpb.Service, error)
}

// Ensure GCPClient implements Client
var _ Client = (*GCPClient)(nil)

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct {
	mu     sync.Mutex
	creds  *google.Credentials
	client ServicesClientWrapper
}

func (c *GCPClient) getClient(ctx context.Context) (ServicesClientWrapper, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil {
		return c.client, nil
	}

	if c.creds == nil {
		creds, err := client.FindDefaultCredentials(ctx, run.DefaultAuthScopes()...)
		if err != nil {
			return nil, fmt.Errorf("failed to find default credentials: %w", err)
		}
		c.creds = creds
	}

	cClient, err := createServicesClient(ctx, option.WithCredentials(c.creds))
	if err != nil {
		return nil, err
	}
	c.client = cClient

	return c.client, nil
}

// ListServices lists services for a project and region.
func (c *GCPClient) ListServices(ctx context.Context, project, region string) ([]*runpb.Service, error) {
	cClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &runpb.ListServicesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var services []*runpb.Service
	it := cClient.ListServices(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, client.WrapError(err)
		}
		services = append(services, resp)
	}

	return services, nil
}

// GetService gets a single service.
func (c *GCPClient) GetService(ctx context.Context, name string) (*runpb.Service, error) {
	cClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	return cClient.GetService(ctx, &runpb.GetServiceRequest{Name: name})
}

// UpdateService updates a service.
func (c *GCPClient) UpdateService(ctx context.Context, service *runpb.Service) (*runpb.Service, error) {
	cClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	op, err := cClient.UpdateService(ctx, &runpb.UpdateServiceRequest{Service: service})
	if err != nil {
		return nil, err
	}

	return op.Wait(ctx)
}
