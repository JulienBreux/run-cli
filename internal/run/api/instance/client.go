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

package instance

import (
	"context"
	"fmt"
	"sync"

	run "cloud.google.com/go/run/apiv2"
	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	runv2 "google.golang.org/api/run/v2"
)

// Interfaces for mocking
type InstancesClientWrapper interface {
	ListInstances(ctx context.Context, req *runpb.ListInstancesRequest, opts ...gax.CallOption) InstanceIteratorWrapper
	GetInstance(ctx context.Context, req *runpb.GetInstanceRequest, opts ...gax.CallOption) (*runpb.Instance, error)
	StartInstance(ctx context.Context, req *runpb.StartInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error)
	StopInstance(ctx context.Context, req *runpb.StopInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error)
	DeleteInstance(ctx context.Context, req *runpb.DeleteInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error)
	UpdateAuthentication(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error)
	Close() error
}

type InstanceIteratorWrapper interface {
	Next() (*runpb.Instance, error)
}

type InstanceOperationWrapper interface {
	Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Instance, error)
}

// Variables for dependency injection
var createInstancesClient = func(ctx context.Context, opts ...option.ClientOption) (InstancesClientWrapper, error) {
	c, err := run.NewInstancesClient(ctx, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPInstancesClientWrapper{client: c}, nil
}

// Real implementations
type GCPInstancesClientWrapper struct {
	client *run.InstancesClient
}

func (w *GCPInstancesClientWrapper) ListInstances(ctx context.Context, req *runpb.ListInstancesRequest, opts ...gax.CallOption) InstanceIteratorWrapper {
	return &GCPInstanceIteratorWrapper{it: w.client.ListInstances(ctx, req, opts...)}
}

func (w *GCPInstancesClientWrapper) GetInstance(ctx context.Context, req *runpb.GetInstanceRequest, opts ...gax.CallOption) (*runpb.Instance, error) {
	return w.client.GetInstance(ctx, req, opts...)
}

func (w *GCPInstancesClientWrapper) StartInstance(ctx context.Context, req *runpb.StartInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
	op, err := w.client.StartInstance(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPInstanceOperationWrapper{waitFunc: op.Wait}, nil
}

func (w *GCPInstancesClientWrapper) StopInstance(ctx context.Context, req *runpb.StopInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
	op, err := w.client.StopInstance(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPInstanceOperationWrapper{waitFunc: op.Wait}, nil
}

func (w *GCPInstancesClientWrapper) DeleteInstance(ctx context.Context, req *runpb.DeleteInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
	op, err := w.client.DeleteInstance(ctx, req, opts...)
	if err != nil {
		return nil, err
	}
	return &GCPInstanceOperationWrapper{waitFunc: op.Wait}, nil
}

func (w *GCPInstancesClientWrapper) UpdateAuthentication(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error) {
	runService, err := runv2.NewService(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create run service: %w", err)
	}
	inst := &runv2.GoogleCloudRunV2Instance{
		InvokerIamDisabled: allowUnauthenticated,
	}
	_, err = runService.Projects.Locations.Instances.Patch(name, inst).UpdateMask("invoker_iam_disabled").Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to update instance authentication: %w", err)
	}
	return w.GetInstance(ctx, &runpb.GetInstanceRequest{Name: name})
}

func (w *GCPInstancesClientWrapper) Close() error {
	return w.client.Close()
}

type GCPInstanceIteratorWrapper struct {
	it *run.InstanceIterator
}

func (w *GCPInstanceIteratorWrapper) Next() (*runpb.Instance, error) {
	return w.it.Next()
}

type GCPInstanceOperationWrapper struct {
	waitFunc func(ctx context.Context, opts ...gax.CallOption) (*runpb.Instance, error)
}

func (w *GCPInstanceOperationWrapper) Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Instance, error) {
	return w.waitFunc(ctx, opts...)
}

// Client defines the interface for Cloud Run Instance operations.
type Client interface {
	ListInstances(ctx context.Context, project, region string) ([]*runpb.Instance, error)
	GetInstance(ctx context.Context, name string) (*runpb.Instance, error)
	StartInstance(ctx context.Context, name string) (*runpb.Instance, error)
	StopInstance(ctx context.Context, name string) (*runpb.Instance, error)
	DeleteInstance(ctx context.Context, name string) (*runpb.Instance, error)
	UpdateAuthentication(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error)
}

var _ Client = (*GCPClient)(nil)

// GCPClient is the Google Cloud Platform implementation of Client.
type GCPClient struct {
	mu     sync.Mutex
	client InstancesClientWrapper
}

func (c *GCPClient) getClient(ctx context.Context) (InstancesClientWrapper, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.client != nil {
		return c.client, nil
	}

	bgCtx := context.Background()
	creds, err := client.FindDefaultCredentials(bgCtx, run.DefaultAuthScopes()...)
	if err != nil {
		return nil, fmt.Errorf("failed to find default credentials: %w", err)
	}

	cClient, err := createInstancesClient(bgCtx, option.WithCredentials(creds))
	if err != nil {
		return nil, err
	}
	c.client = cClient
	return c.client, nil
}

// ListInstances lists instances for a project and region.
func (c *GCPClient) ListInstances(ctx context.Context, project, region string) ([]*runpb.Instance, error) {
	cClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	req := &runpb.ListInstancesRequest{
		Parent: fmt.Sprintf("projects/%s/locations/%s", project, region),
	}

	var instances []*runpb.Instance
	it := cClient.ListInstances(ctx, req)
	for {
		resp, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, client.WrapError(err)
		}
		instances = append(instances, resp)
	}

	return instances, nil
}

// GetInstance gets an instance by full resource name.
func (c *GCPClient) GetInstance(ctx context.Context, name string) (*runpb.Instance, error) {
	cClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cClient.GetInstance(ctx, &runpb.GetInstanceRequest{Name: name})
	if err != nil {
		return nil, client.WrapError(err)
	}
	return resp, nil
}

// StartInstance starts an instance.
func (c *GCPClient) StartInstance(ctx context.Context, name string) (*runpb.Instance, error) {
	cClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	op, err := cClient.StartInstance(ctx, &runpb.StartInstanceRequest{Name: name})
	if err != nil {
		return nil, client.WrapError(err)
	}

	return op.Wait(ctx)
}

// StopInstance stops an instance.
func (c *GCPClient) StopInstance(ctx context.Context, name string) (*runpb.Instance, error) {
	cClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	op, err := cClient.StopInstance(ctx, &runpb.StopInstanceRequest{Name: name})
	if err != nil {
		return nil, client.WrapError(err)
	}

	return op.Wait(ctx)
}

// DeleteInstance deletes an instance.
func (c *GCPClient) DeleteInstance(ctx context.Context, name string) (*runpb.Instance, error) {
	cClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	op, err := cClient.DeleteInstance(ctx, &runpb.DeleteInstanceRequest{Name: name})
	if err != nil {
		return nil, client.WrapError(err)
	}

	return op.Wait(ctx)
}

// UpdateAuthentication updates the authentication setting for an instance.
func (c *GCPClient) UpdateAuthentication(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error) {
	cClient, err := c.getClient(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := cClient.UpdateAuthentication(ctx, name, allowUnauthenticated)
	if err != nil {
		return nil, client.WrapError(err)
	}

	return resp, nil
}
