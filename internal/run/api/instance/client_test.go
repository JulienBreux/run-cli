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
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockClient is a mock implementation of Client.
type MockClient struct {
	ListInstancesFunc  func(ctx context.Context, project, region string) ([]*runpb.Instance, error)
	GetInstanceFunc    func(ctx context.Context, name string) (*runpb.Instance, error)
	StartInstanceFunc  func(ctx context.Context, name string) (*runpb.Instance, error)
	StopInstanceFunc   func(ctx context.Context, name string) (*runpb.Instance, error)
	DeleteInstanceFunc func(ctx context.Context, name string) (*runpb.Instance, error)
	UpdateAuthenticationFunc func(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error)
}

func (m *MockClient) ListInstances(ctx context.Context, project, region string) ([]*runpb.Instance, error) {
	if m.ListInstancesFunc != nil {
		return m.ListInstancesFunc(ctx, project, region)
	}
	return nil, nil
}

func (m *MockClient) GetInstance(ctx context.Context, name string) (*runpb.Instance, error) {
	if m.GetInstanceFunc != nil {
		return m.GetInstanceFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockClient) StartInstance(ctx context.Context, name string) (*runpb.Instance, error) {
	if m.StartInstanceFunc != nil {
		return m.StartInstanceFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockClient) StopInstance(ctx context.Context, name string) (*runpb.Instance, error) {
	if m.StopInstanceFunc != nil {
		return m.StopInstanceFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockClient) DeleteInstance(ctx context.Context, name string) (*runpb.Instance, error) {
	if m.DeleteInstanceFunc != nil {
		return m.DeleteInstanceFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockClient) UpdateAuthentication(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error) {
	if m.UpdateAuthenticationFunc != nil {
		return m.UpdateAuthenticationFunc(ctx, name, allowUnauthenticated)
	}
	return nil, nil
}

func TestMapInstance(t *testing.T) {
	now := time.Now()
	resp := &runpb.Instance{
		Name:        "projects/my-project/locations/us-central1/instances/my-inst",
		Description: "Instance desc",
		Uid:         "uid-123",
		Generation:  2,
		CreateTime:  timestamppb.New(now),
		UpdateTime:  timestamppb.New(now),
		Creator:     "user@example.com",
		Labels: map[string]string{
			"env": "prod",
		},
		Containers: []*runpb.Container{
			{Name: "c1", Image: "gcr.io/test/app:v1"},
		},
		ContainerStatuses: []*runpb.ContainerStatus{
			{Name: "c1", ImageDigest: "sha256:123"},
		},
		Urls: []string{"https://my-inst.run.app"},
		TerminalCondition: &runpb.Condition{
			State:              runpb.Condition_CONDITION_SUCCEEDED,
			Message:            "Instance ready",
			LastTransitionTime: timestamppb.New(now),
		},
		Conditions: []*runpb.Condition{
			{
				Type:               "Ready",
				State:              runpb.Condition_CONDITION_SUCCEEDED,
				Message:            "Ready",
				LastTransitionTime: timestamppb.New(now),
			},
		},
		LogUri: "https://console.cloud.google.com/logs",
	}

	result := mapInstance(resp, "us-central1", "my-project")

	assert.Equal(t, resp.Name, result.Name)
	assert.Equal(t, "my-inst", result.ID())
	assert.Equal(t, "Instance desc", result.Description)
	assert.Equal(t, "uid-123", result.UID)
	assert.Equal(t, int64(2), result.Generation)
	assert.Equal(t, "us-central1", result.Region)
	assert.Equal(t, "my-project", result.Project)
	assert.Equal(t, "user@example.com", result.Creator)
	assert.Equal(t, "Ready", result.Status())
	assert.Len(t, result.Containers, 1)
	assert.Equal(t, "c1", result.Containers[0].Name)
	assert.Len(t, result.ContainerStatuses, 1)
	assert.Equal(t, "sha256:123", result.ContainerStatuses[0].ImageDigest)
	assert.Equal(t, []string{"https://my-inst.run.app"}, result.URLs)
}

func TestMapInstance_NilFields(t *testing.T) {
	resp := &runpb.Instance{
		Name: "projects/p/locations/r/instances/inst1",
	}

	result := mapInstance(resp, "r", "p")

	assert.Equal(t, resp.Name, result.Name)
	assert.Nil(t, result.TerminalCondition)
	assert.Empty(t, result.Containers)
}

func TestList(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListInstancesFunc = func(ctx context.Context, project, region string) ([]*runpb.Instance, error) {
		return []*runpb.Instance{
			{Name: "projects/p/locations/r/instances/inst1"},
			{Name: "projects/p/locations/r/instances/inst2"},
		}, nil
	}

	instances, err := List("p", "r")
	assert.NoError(t, err)
	assert.Len(t, instances, 2)
	assert.Equal(t, "inst1", instances[0].ID())
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListInstancesFunc = func(ctx context.Context, project, region string) ([]*runpb.Instance, error) {
		return nil, assert.AnError
	}

	instances, err := List("p", "r")
	assert.Error(t, err)
	assert.Nil(t, instances)
}

func TestList_AllRegions(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListInstancesFunc = func(ctx context.Context, project, region string) ([]*runpb.Instance, error) {
		if region == "us-central1" {
			return []*runpb.Instance{{Name: "projects/p/locations/us-central1/instances/inst-us"}}, nil
		}
		return []*runpb.Instance{}, nil
	}

	instances, err := List("p", api_region.ALL)
	assert.NoError(t, err)

	found := false
	for _, inst := range instances {
		if inst.ID() == "inst-us" && inst.Region == "us-central1" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestGet(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.GetInstanceFunc = func(ctx context.Context, name string) (*runpb.Instance, error) {
		assert.Equal(t, "projects/p/locations/r/instances/inst1", name)
		return &runpb.Instance{Name: name}, nil
	}

	inst, err := Get("p", "r", "inst1")
	assert.NoError(t, err)
	assert.NotNil(t, inst)
	assert.Equal(t, "inst1", inst.ID())
}

func TestGet_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.GetInstanceFunc = func(ctx context.Context, name string) (*runpb.Instance, error) {
		return nil, assert.AnError
	}

	inst, err := Get("p", "r", "inst1")
	assert.Error(t, err)
	assert.Nil(t, inst)
}

func TestStart(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.StartInstanceFunc = func(ctx context.Context, name string) (*runpb.Instance, error) {
		assert.Equal(t, "projects/p/locations/r/instances/inst1", name)
		return &runpb.Instance{Name: name}, nil
	}

	err := Start("p", "r", "inst1")
	assert.NoError(t, err)
}

func TestStart_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.StartInstanceFunc = func(ctx context.Context, name string) (*runpb.Instance, error) {
		return nil, assert.AnError
	}

	err := Start("p", "r", "inst1")
	assert.Error(t, err)
}

func TestStop(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.StopInstanceFunc = func(ctx context.Context, name string) (*runpb.Instance, error) {
		assert.Equal(t, "projects/p/locations/r/instances/inst1", name)
		return &runpb.Instance{Name: name}, nil
	}

	err := Stop("p", "r", "inst1")
	assert.NoError(t, err)
}

func TestStop_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.StopInstanceFunc = func(ctx context.Context, name string) (*runpb.Instance, error) {
		return nil, assert.AnError
	}

	err := Stop("p", "r", "inst1")
	assert.Error(t, err)
}

func TestDelete(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.DeleteInstanceFunc = func(ctx context.Context, name string) (*runpb.Instance, error) {
		assert.Equal(t, "projects/p/locations/r/instances/inst1", name)
		return &runpb.Instance{Name: name}, nil
	}

	err := Delete("p", "r", "inst1")
	assert.NoError(t, err)
}

func TestDelete_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.DeleteInstanceFunc = func(ctx context.Context, name string) (*runpb.Instance, error) {
		return nil, assert.AnError
	}

	err := Delete("p", "r", "inst1")
	assert.Error(t, err)
}

func TestUpdateAuthentication(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.UpdateAuthenticationFunc = func(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error) {
		assert.Equal(t, "projects/p/locations/r/instances/inst1", name)
		assert.True(t, allowUnauthenticated)
		return &runpb.Instance{
			Name:               name,
			InvokerIamDisabled: allowUnauthenticated,
		}, nil
	}

	inst, err := UpdateAuthentication(context.Background(), "p", "r", "inst1", true)
	assert.NoError(t, err)
	assert.NotNil(t, inst)
	assert.True(t, inst.InvokerIamDisabled)
}

func TestUpdateAuthentication_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.UpdateAuthenticationFunc = func(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error) {
		return nil, assert.AnError
	}

	inst, err := UpdateAuthentication(context.Background(), "p", "r", "inst1", false)
	assert.Error(t, err)
	assert.Nil(t, inst)
}

// --- Mocks for GCPClient testing ---

type MockInstancesClientWrapper struct {
	ListInstancesFunc  func(ctx context.Context, req *runpb.ListInstancesRequest, opts ...gax.CallOption) InstanceIteratorWrapper
	GetInstanceFunc    func(ctx context.Context, req *runpb.GetInstanceRequest, opts ...gax.CallOption) (*runpb.Instance, error)
	StartInstanceFunc  func(ctx context.Context, req *runpb.StartInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error)
	StopInstanceFunc   func(ctx context.Context, req *runpb.StopInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error)
	DeleteInstanceFunc func(ctx context.Context, req *runpb.DeleteInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error)
	UpdateAuthenticationFunc func(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error)
	CloseFunc          func() error
}

func (m *MockInstancesClientWrapper) ListInstances(ctx context.Context, req *runpb.ListInstancesRequest, opts ...gax.CallOption) InstanceIteratorWrapper {
	if m.ListInstancesFunc != nil {
		return m.ListInstancesFunc(ctx, req, opts...)
	}
	return &MockInstanceIteratorWrapper{}
}

func (m *MockInstancesClientWrapper) GetInstance(ctx context.Context, req *runpb.GetInstanceRequest, opts ...gax.CallOption) (*runpb.Instance, error) {
	if m.GetInstanceFunc != nil {
		return m.GetInstanceFunc(ctx, req, opts...)
	}
	return nil, nil
}

func (m *MockInstancesClientWrapper) StartInstance(ctx context.Context, req *runpb.StartInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
	if m.StartInstanceFunc != nil {
		return m.StartInstanceFunc(ctx, req, opts...)
	}
	return nil, nil
}

func (m *MockInstancesClientWrapper) StopInstance(ctx context.Context, req *runpb.StopInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
	if m.StopInstanceFunc != nil {
		return m.StopInstanceFunc(ctx, req, opts...)
	}
	return nil, nil
}

func (m *MockInstancesClientWrapper) DeleteInstance(ctx context.Context, req *runpb.DeleteInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
	if m.DeleteInstanceFunc != nil {
		return m.DeleteInstanceFunc(ctx, req, opts...)
	}
	return nil, nil
}

func (m *MockInstancesClientWrapper) UpdateAuthentication(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error) {
	if m.UpdateAuthenticationFunc != nil {
		return m.UpdateAuthenticationFunc(ctx, name, allowUnauthenticated)
	}
	return nil, nil
}

func (m *MockInstancesClientWrapper) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

type MockInstanceIteratorWrapper struct {
	Items []*runpb.Instance
	Index int
	Err   error
}

func (m *MockInstanceIteratorWrapper) Next() (*runpb.Instance, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.Index >= len(m.Items) {
		return nil, iterator.Done
	}
	item := m.Items[m.Index]
	m.Index++
	return item, nil
}

type MockInstanceOperationWrapper struct {
	WaitFunc func(ctx context.Context, opts ...gax.CallOption) (*runpb.Instance, error)
}

func (m *MockInstanceOperationWrapper) Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Instance, error) {
	if m.WaitFunc != nil {
		return m.WaitFunc(ctx, opts...)
	}
	return nil, nil
}

func TestGCPClient_ListInstances(t *testing.T) {
	ctx := context.Background()
	mockWrapper := &MockInstancesClientWrapper{
		ListInstancesFunc: func(ctx context.Context, req *runpb.ListInstancesRequest, opts ...gax.CallOption) InstanceIteratorWrapper {
			assert.Equal(t, "projects/p/locations/r", req.Parent)
			return &MockInstanceIteratorWrapper{
				Items: []*runpb.Instance{
					{Name: "projects/p/locations/r/instances/i1"},
					{Name: "projects/p/locations/r/instances/i2"},
				},
			}
		},
	}

	client := &GCPClient{client: mockWrapper}
	instances, err := client.ListInstances(ctx, "p", "r")
	assert.NoError(t, err)
	assert.Len(t, instances, 2)
}

func TestGCPClient_ListInstances_IteratorError(t *testing.T) {
	ctx := context.Background()
	mockWrapper := &MockInstancesClientWrapper{
		ListInstancesFunc: func(ctx context.Context, req *runpb.ListInstancesRequest, opts ...gax.CallOption) InstanceIteratorWrapper {
			return &MockInstanceIteratorWrapper{
				Err: errors.New("rpc error"),
			}
		},
	}

	client := &GCPClient{client: mockWrapper}
	instances, err := client.ListInstances(ctx, "p", "r")
	assert.Error(t, err)
	assert.Nil(t, instances)
}

func TestGCPClient_GetInstance(t *testing.T) {
	ctx := context.Background()
	mockWrapper := &MockInstancesClientWrapper{
		GetInstanceFunc: func(ctx context.Context, req *runpb.GetInstanceRequest, opts ...gax.CallOption) (*runpb.Instance, error) {
			assert.Equal(t, "projects/p/locations/r/instances/i1", req.Name)
			return &runpb.Instance{Name: req.Name}, nil
		},
	}

	client := &GCPClient{client: mockWrapper}
	inst, err := client.GetInstance(ctx, "projects/p/locations/r/instances/i1")
	assert.NoError(t, err)
	assert.NotNil(t, inst)
}

func TestGCPClient_StartInstance(t *testing.T) {
	ctx := context.Background()
	mockWrapper := &MockInstancesClientWrapper{
		StartInstanceFunc: func(ctx context.Context, req *runpb.StartInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
			assert.Equal(t, "inst1", req.Name)
			return &MockInstanceOperationWrapper{
				WaitFunc: func(ctx context.Context, opts ...gax.CallOption) (*runpb.Instance, error) {
					return &runpb.Instance{Name: "inst1"}, nil
				},
			}, nil
		},
	}

	client := &GCPClient{client: mockWrapper}
	res, err := client.StartInstance(ctx, "inst1")
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGCPClient_StopInstance(t *testing.T) {
	ctx := context.Background()
	mockWrapper := &MockInstancesClientWrapper{
		StopInstanceFunc: func(ctx context.Context, req *runpb.StopInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
			assert.Equal(t, "inst1", req.Name)
			return &MockInstanceOperationWrapper{
				WaitFunc: func(ctx context.Context, opts ...gax.CallOption) (*runpb.Instance, error) {
					return &runpb.Instance{Name: "inst1"}, nil
				},
			}, nil
		},
	}

	client := &GCPClient{client: mockWrapper}
	res, err := client.StopInstance(ctx, "inst1")
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGCPClient_DeleteInstance(t *testing.T) {
	ctx := context.Background()
	mockWrapper := &MockInstancesClientWrapper{
		DeleteInstanceFunc: func(ctx context.Context, req *runpb.DeleteInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
			assert.Equal(t, "inst1", req.Name)
			return &MockInstanceOperationWrapper{
				WaitFunc: func(ctx context.Context, opts ...gax.CallOption) (*runpb.Instance, error) {
					return &runpb.Instance{Name: "inst1"}, nil
				},
			}, nil
		},
	}

	client := &GCPClient{client: mockWrapper}
	res, err := client.DeleteInstance(ctx, "inst1")
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestGCPClient_UpdateAuthentication(t *testing.T) {
	ctx := context.Background()
	mockWrapper := &MockInstancesClientWrapper{
		UpdateAuthenticationFunc: func(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error) {
			assert.Equal(t, "inst1", name)
			assert.True(t, allowUnauthenticated)
			return &runpb.Instance{Name: "inst1", InvokerIamDisabled: true}, nil
		},
	}

	client := &GCPClient{client: mockWrapper}
	res, err := client.UpdateAuthentication(ctx, "inst1", true)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.True(t, res.InvokerIamDisabled)
}

func TestGCPClient_Errors(t *testing.T) {
	ctx := context.Background()
	mockWrapper := &MockInstancesClientWrapper{
		GetInstanceFunc: func(ctx context.Context, req *runpb.GetInstanceRequest, opts ...gax.CallOption) (*runpb.Instance, error) {
			return nil, errors.New("get error")
		},
		StartInstanceFunc: func(ctx context.Context, req *runpb.StartInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
			return nil, errors.New("start error")
		},
		StopInstanceFunc: func(ctx context.Context, req *runpb.StopInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
			return nil, errors.New("stop error")
		},
		DeleteInstanceFunc: func(ctx context.Context, req *runpb.DeleteInstanceRequest, opts ...gax.CallOption) (InstanceOperationWrapper, error) {
			return nil, errors.New("delete error")
		},
		UpdateAuthenticationFunc: func(ctx context.Context, name string, allowUnauthenticated bool) (*runpb.Instance, error) {
			return nil, errors.New("update auth error")
		},
	}

	client := &GCPClient{client: mockWrapper}

	_, err := client.GetInstance(ctx, "name")
	assert.Error(t, err)

	_, err = client.StartInstance(ctx, "name")
	assert.Error(t, err)

	_, err = client.StopInstance(ctx, "name")
	assert.Error(t, err)

	_, err = client.DeleteInstance(ctx, "name")
	assert.Error(t, err)

	_, err = client.UpdateAuthentication(ctx, "name", true)
	assert.Error(t, err)
}

func TestFormatInstanceName(t *testing.T) {
	// Already formatted
	fullName := "projects/p/locations/r/instances/i1"
	assert.Equal(t, fullName, formatInstanceName("p", "r", fullName))

	// Short name
	assert.Equal(t, fullName, formatInstanceName("p", "r", "i1"))
}

func TestGCPClient_GetClient_Cached(t *testing.T) {
	mockWrapper := &MockInstancesClientWrapper{}
	c := &GCPClient{client: mockWrapper}
	res, err := c.getClient(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, mockWrapper, res)
}

func TestWrappers_Delegation(t *testing.T) {
	t.Run("GCPInstancesClientWrapper", func(t *testing.T) {
		w := &GCPInstancesClientWrapper{client: nil}
		assert.Panics(t, func() { _ = w.ListInstances(context.Background(), nil) })
		assert.Panics(t, func() { _, _ = w.GetInstance(context.Background(), nil) })
		assert.Panics(t, func() { _, _ = w.StartInstance(context.Background(), nil) })
		assert.Panics(t, func() { _, _ = w.StopInstance(context.Background(), nil) })
		assert.Panics(t, func() { _, _ = w.DeleteInstance(context.Background(), nil) })
		assert.Panics(t, func() { _ = w.Close() })
		_, err := w.UpdateAuthentication(context.Background(), "name", true)
		assert.Error(t, err)
	})

	t.Run("GCPInstanceIteratorWrapper", func(t *testing.T) {
		it := &GCPInstanceIteratorWrapper{it: nil}
		assert.Panics(t, func() { _, _ = it.Next() })
	})

	t.Run("GCPInstanceOperationWrapper", func(t *testing.T) {
		op := &GCPInstanceOperationWrapper{
			waitFunc: func(ctx context.Context, opts ...gax.CallOption) (*runpb.Instance, error) {
				return nil, errors.New("op wait error")
			},
		}
		_, err := op.Wait(context.Background())
		assert.Error(t, err)
	})
}


