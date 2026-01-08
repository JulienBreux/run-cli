package workerpool

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockClient (High Level)
type MockClient struct {
	ListWorkerPoolsFunc  func(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error)
	GetWorkerPoolFunc    func(ctx context.Context, name string) (*runpb.WorkerPool, error)
	UpdateWorkerPoolFunc func(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error)
}

func (m *MockClient) ListWorkerPools(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
	if m.ListWorkerPoolsFunc != nil {
		return m.ListWorkerPoolsFunc(ctx, project, region)
	}
	return nil, nil
}

func (m *MockClient) GetWorkerPool(ctx context.Context, name string) (*runpb.WorkerPool, error) {
	if m.GetWorkerPoolFunc != nil {
		return m.GetWorkerPoolFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockClient) UpdateWorkerPool(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error) {
	if m.UpdateWorkerPoolFunc != nil {
		return m.UpdateWorkerPoolFunc(ctx, workerPool)
	}
	return nil, nil
}

func TestMapWorkerPool(t *testing.T) {
	now := time.Now()
	instanceCount := int32(2)
	resp := &runpb.WorkerPool{
		Name:         "projects/my-project/locations/us-central1/workerPools/my-pool",
		LastModifier: "user@example.com",
		UpdateTime:   timestamppb.New(now),
		Scaling: &runpb.WorkerPoolScaling{
			ManualInstanceCount: &instanceCount,
		},
		Labels: map[string]string{
			"env": "prod",
		},
	}

	result := mapWorkerPool(resp, "my-project", "us-central1")

	assert.Equal(t, "my-pool", result.DisplayName)
	assert.Equal(t, resp.Name, result.Name)
	assert.Equal(t, "user@example.com", result.LastModifier)
	assert.Equal(t, now.Unix(), result.UpdateTime.Unix())
	assert.Equal(t, "my-project", result.Project)
	assert.Equal(t, "us-central1", result.Region)
	
	// Scaling
	assert.NotNil(t, result.Scaling)
	assert.Equal(t, int32(2), result.Scaling.ManualInstanceCount)
	
	// Labels
	assert.Equal(t, "prod", result.Labels["env"])
}

func TestList(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListWorkerPoolsFunc = func(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
		return []*runpb.WorkerPool{{Name: "pool1"}}, nil
	}

	pools, err := List("p", "r")
	assert.NoError(t, err)
	assert.Len(t, pools, 1)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListWorkerPoolsFunc = func(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
		return nil, assert.AnError
	}

	pools, err := List("p", "r")
	assert.Error(t, err)
	assert.Nil(t, pools)
}

func TestList_AllRegions(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListWorkerPoolsFunc = func(ctx context.Context, project, region string) ([]*runpb.WorkerPool, error) {
		return []*runpb.WorkerPool{{Name: "projects/" + project + "/locations/" + region + "/workerPools/pool-" + region}}, nil
	}

	pools, err := List("p", api_region.ALL)
	assert.NoError(t, err)
	// We expect at least one pool per region. api_region.List() has > 0 regions.
	assert.NotEmpty(t, pools)
}

func TestUpdateScaling(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.GetWorkerPoolFunc = func(ctx context.Context, name string) (*runpb.WorkerPool, error) {
		return &runpb.WorkerPool{Name: name}, nil
	}
	mock.UpdateWorkerPoolFunc = func(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error) {
		assert.Equal(t, int32(5), *workerPool.Scaling.ManualInstanceCount)
		return workerPool, nil
	}

	result, err := UpdateScaling(context.Background(), "p", "r", "pool1", 5)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int32(5), result.Scaling.ManualInstanceCount)
}

func TestUpdateScaling_GetError(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.GetWorkerPoolFunc = func(ctx context.Context, name string) (*runpb.WorkerPool, error) {
		return nil, assert.AnError
	}

	result, err := UpdateScaling(context.Background(), "p", "r", "pool1", 5)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to get worker pool")
}

func TestUpdateScaling_UpdateError(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.GetWorkerPoolFunc = func(ctx context.Context, name string) (*runpb.WorkerPool, error) {
		return &runpb.WorkerPool{Name: name}, nil
	}
	mock.UpdateWorkerPoolFunc = func(ctx context.Context, workerPool *runpb.WorkerPool) (*runpb.WorkerPool, error) {
		return nil, assert.AnError
	}

	result, err := UpdateScaling(context.Background(), "p", "r", "pool1", 5)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to update worker pool")
}

// --- Mocks for GCPClient testing ---

type MockWorkerPoolsClientWrapper struct {
	ListWorkerPoolsFunc  func(ctx context.Context, req *runpb.ListWorkerPoolsRequest, opts ...gax.CallOption) WorkerPoolIteratorWrapper
	GetWorkerPoolFunc    func(ctx context.Context, req *runpb.GetWorkerPoolRequest, opts ...gax.CallOption) (*runpb.WorkerPool, error)
	UpdateWorkerPoolFunc func(ctx context.Context, req *runpb.UpdateWorkerPoolRequest, opts ...gax.CallOption) (UpdateWorkerPoolOperationWrapper, error)
	CloseFunc            func() error
}

func (m *MockWorkerPoolsClientWrapper) ListWorkerPools(ctx context.Context, req *runpb.ListWorkerPoolsRequest, opts ...gax.CallOption) WorkerPoolIteratorWrapper {
	if m.ListWorkerPoolsFunc != nil {
		return m.ListWorkerPoolsFunc(ctx, req, opts...)
	}
	return &MockWorkerPoolIteratorWrapper{}
}

func (m *MockWorkerPoolsClientWrapper) GetWorkerPool(ctx context.Context, req *runpb.GetWorkerPoolRequest, opts ...gax.CallOption) (*runpb.WorkerPool, error) {
	if m.GetWorkerPoolFunc != nil {
		return m.GetWorkerPoolFunc(ctx, req, opts...)
	}
	return nil, nil
}

func (m *MockWorkerPoolsClientWrapper) UpdateWorkerPool(ctx context.Context, req *runpb.UpdateWorkerPoolRequest, opts ...gax.CallOption) (UpdateWorkerPoolOperationWrapper, error) {
	if m.UpdateWorkerPoolFunc != nil {
		return m.UpdateWorkerPoolFunc(ctx, req, opts...)
	}
	return nil, nil
}

func (m *MockWorkerPoolsClientWrapper) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

type MockWorkerPoolIteratorWrapper struct {
	Items []*runpb.WorkerPool
	Index int
	Err   error
}

func (m *MockWorkerPoolIteratorWrapper) Next() (*runpb.WorkerPool, error) {
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

type MockUpdateWorkerPoolOperationWrapper struct {
	WaitFunc func(ctx context.Context, opts ...gax.CallOption) (*runpb.WorkerPool, error)
}

func (m *MockUpdateWorkerPoolOperationWrapper) Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.WorkerPool, error) {
	if m.WaitFunc != nil {
		return m.WaitFunc(ctx, opts...)
	}
	return nil, nil
}

func TestGCPClient_ListWorkerPools(t *testing.T) {
	origFindCreds := findDefaultCredentials
	origCreateClient := createWorkerPoolsClient
	defer func() {
		findDefaultCredentials = origFindCreds
		createWorkerPoolsClient = origCreateClient
	}()

	findDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createWorkerPoolsClient = func(ctx context.Context, opts ...option.ClientOption) (WorkerPoolsClientWrapper, error) {
			return &MockWorkerPoolsClientWrapper{
				ListWorkerPoolsFunc: func(ctx context.Context, req *runpb.ListWorkerPoolsRequest, opts ...gax.CallOption) WorkerPoolIteratorWrapper {
					return &MockWorkerPoolIteratorWrapper{
						Items: []*runpb.WorkerPool{{Name: "pool1"}},
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}

		client := &GCPClient{}
		pools, err := client.ListWorkerPools(context.Background(), "p", "r")
		assert.NoError(t, err)
		assert.Len(t, pools, 1)
		assert.Equal(t, "pool1", pools[0].Name)
	})

	t.Run("Auth Error", func(t *testing.T) {
		findDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("auth failed")
		}
		client := &GCPClient{}
		_, err := client.ListWorkerPools(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find default credentials")
	})

	t.Run("Iterator Auth Error", func(t *testing.T) {
		findDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createWorkerPoolsClient = func(ctx context.Context, opts ...option.ClientOption) (WorkerPoolsClientWrapper, error) {
			return &MockWorkerPoolsClientWrapper{
				ListWorkerPoolsFunc: func(ctx context.Context, req *runpb.ListWorkerPoolsRequest, opts ...gax.CallOption) WorkerPoolIteratorWrapper {
					return &MockWorkerPoolIteratorWrapper{
						Err: errors.New("Unauthenticated request"),
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		client := &GCPClient{}
		_, err := client.ListWorkerPools(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed")
	})
}

func TestGCPClient_GetWorkerPool(t *testing.T) {
	origFindCreds := findDefaultCredentials
	origCreateClient := createWorkerPoolsClient
	defer func() {
		findDefaultCredentials = origFindCreds
		createWorkerPoolsClient = origCreateClient
	}()

	findDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createWorkerPoolsClient = func(ctx context.Context, opts ...option.ClientOption) (WorkerPoolsClientWrapper, error) {
			return &MockWorkerPoolsClientWrapper{
				GetWorkerPoolFunc: func(ctx context.Context, req *runpb.GetWorkerPoolRequest, opts ...gax.CallOption) (*runpb.WorkerPool, error) {
					return &runpb.WorkerPool{Name: "pool1"}, nil
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		
		client := &GCPClient{}
		pool, err := client.GetWorkerPool(context.Background(), "pool1")
		assert.NoError(t, err)
		assert.Equal(t, "pool1", pool.Name)
	})
	
	t.Run("Get Error", func(t *testing.T) {
		createWorkerPoolsClient = func(ctx context.Context, opts ...option.ClientOption) (WorkerPoolsClientWrapper, error) {
			return &MockWorkerPoolsClientWrapper{
				GetWorkerPoolFunc: func(ctx context.Context, req *runpb.GetWorkerPoolRequest, opts ...gax.CallOption) (*runpb.WorkerPool, error) {
					return nil, errors.New("get error")
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		
		client := &GCPClient{}
		_, err := client.GetWorkerPool(context.Background(), "pool1")
		assert.Error(t, err)
	})
}

func TestGCPClient_UpdateWorkerPool(t *testing.T) {
	origFindCreds := findDefaultCredentials
	origCreateClient := createWorkerPoolsClient
	defer func() {
		findDefaultCredentials = origFindCreds
		createWorkerPoolsClient = origCreateClient
	}()

	findDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createWorkerPoolsClient = func(ctx context.Context, opts ...option.ClientOption) (WorkerPoolsClientWrapper, error) {
			return &MockWorkerPoolsClientWrapper{
				UpdateWorkerPoolFunc: func(ctx context.Context, req *runpb.UpdateWorkerPoolRequest, opts ...gax.CallOption) (UpdateWorkerPoolOperationWrapper, error) {
					return &MockUpdateWorkerPoolOperationWrapper{
						WaitFunc: func(ctx context.Context, opts ...gax.CallOption) (*runpb.WorkerPool, error) {
							return &runpb.WorkerPool{Name: "pool1-updated"}, nil
						},
					}, nil
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		
		client := &GCPClient{}
		pool, err := client.UpdateWorkerPool(context.Background(), &runpb.WorkerPool{Name: "pool1"})
		assert.NoError(t, err)
		assert.Equal(t, "pool1-updated", pool.Name)
	})

	t.Run("Update Op Start Error", func(t *testing.T) {
		createWorkerPoolsClient = func(ctx context.Context, opts ...option.ClientOption) (WorkerPoolsClientWrapper, error) {
			return &MockWorkerPoolsClientWrapper{
				UpdateWorkerPoolFunc: func(ctx context.Context, req *runpb.UpdateWorkerPoolRequest, opts ...gax.CallOption) (UpdateWorkerPoolOperationWrapper, error) {
					return nil, errors.New("update start error")
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		
		client := &GCPClient{}
		_, err := client.UpdateWorkerPool(context.Background(), &runpb.WorkerPool{Name: "pool1"})
		assert.Error(t, err)
	})
}

func TestWrappers_Delegation(t *testing.T) {
	// Expect panics because nil clients are used
	
	t.Run("GCPWorkerPoolsClientWrapper", func(t *testing.T) {
		w := &GCPWorkerPoolsClientWrapper{client: nil}
		assert.Panics(t, func() { w.ListWorkerPools(context.Background(), nil) })
		assert.Panics(t, func() { w.GetWorkerPool(context.Background(), nil) })
		assert.Panics(t, func() { w.UpdateWorkerPool(context.Background(), nil) })
		assert.Panics(t, func() { w.Close() })
	})
	
	t.Run("GCPWorkerPoolIteratorWrapper", func(t *testing.T) {
		it := &GCPWorkerPoolIteratorWrapper{it: nil}
		assert.Panics(t, func() { it.Next() })
	})
	
	t.Run("GCPUpdateWorkerPoolOperationWrapper", func(t *testing.T) {
		op := &GCPUpdateWorkerPoolOperationWrapper{op: nil}
		assert.Panics(t, func() { op.Wait(context.Background()) })
	})
}