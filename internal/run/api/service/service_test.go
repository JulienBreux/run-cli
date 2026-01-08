package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	api_region "github.com/JulienBreux/run-cli/internal/run/api/region"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestMapService(t *testing.T) {
	// Setup input
	now := time.Now()
	resp := &runpb.Service{
		Name:         "projects/my-project/locations/us-central1/services/my-service",
		Uri:          "https://my-service.run.app",
		LastModifier: "user@example.com",
		UpdateTime:   timestamppb.New(now),
		Scaling: &runpb.ServiceScaling{
			MinInstanceCount: 1,
			MaxInstanceCount: 5,
			ScalingMode:      runpb.ServiceScaling_AUTOMATIC,
		},
		TrafficStatuses: []*runpb.TrafficTargetStatus{
			{
				Type:     runpb.TrafficTargetAllocationType_TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST,
				Percent:  100,
				Revision: "my-service-v1",
			},
		},
		LatestReadyRevision:   "my-service-v1",
		LatestCreatedRevision: "my-service-v2",
		Ingress:               runpb.IngressTraffic_INGRESS_TRAFFIC_ALL,
	}

	// Call function
	result := mapService(resp, "my-project", "us-central1")

	// Assertions
	assert.Equal(t, "my-service", result.Name)
	assert.Equal(t, "https://my-service.run.app", result.URI)
	assert.Equal(t, "user@example.com", result.LastModifier)
	assert.Equal(t, "my-project", result.Project)
	assert.Equal(t, "us-central1", result.Region)
	
	// Scaling
	assert.Equal(t, "AUTOMATIC", result.Scaling.ScalingMode)
	assert.Equal(t, int32(1), result.Scaling.MinInstances)
	assert.Equal(t, int32(5), result.Scaling.MaxInstances)

	// Traffic
	assert.Len(t, result.TrafficStatuses, 1)
	assert.Equal(t, "my-service-v1", result.TrafficStatuses[0].Revision)
	
	// Revisions
	assert.Equal(t, "my-service-v1", result.LatestReadyRevision)
	assert.Equal(t, "my-service-v2", result.LatestCreatedRevision)

	// Networking
	assert.Equal(t, "INGRESS_TRAFFIC_ALL", result.Networking.Ingress)
}

func TestMapService_Complex(t *testing.T) {
	instanceCount := int32(3)
	resp := &runpb.Service{
		Name: "projects/p/locations/r/services/s-complex",
		Scaling: &runpb.ServiceScaling{
			ScalingMode:         runpb.ServiceScaling_MANUAL,
			ManualInstanceCount: &instanceCount,
		},
		Template: &runpb.RevisionTemplate{
			VpcAccess: &runpb.VpcAccess{
				Connector: "connector1",
				Egress:    runpb.VpcAccess_ALL_TRAFFIC,
			},
			ServiceAccount: "sa@example.com",
			EncryptionKey:  "key1",
		},
		BinaryAuthorization: &runpb.BinaryAuthorization{
			BinauthzMethod: &runpb.BinaryAuthorization_UseDefault{
				UseDefault: true,
			},
			BreakglassJustification: "emergency",
		},
	}

	result := mapService(resp, "p", "r")

	// Scaling
	assert.Equal(t, "MANUAL", result.Scaling.ScalingMode)
	assert.Equal(t, int32(3), result.Scaling.ManualInstanceCount)

	// Networking / VPC
	assert.NotNil(t, result.Networking.VpcAccess)
	assert.Equal(t, "connector1", result.Networking.VpcAccess.Connector)
	assert.Equal(t, "ALL_TRAFFIC", result.Networking.VpcAccess.Egress)

	// Security
	assert.Equal(t, "sa@example.com", result.Security.ServiceAccount)
	assert.Equal(t, "key1", result.Security.EncryptionKey)
	assert.Equal(t, "default", result.Security.BinaryAuthorization)
	assert.Equal(t, "emergency", result.Security.BreakglassJustification)
}

// MockClient is a mock implementation of the Client interface.
type MockClient struct {
	ListServicesFunc  func(ctx context.Context, project, region string) ([]*runpb.Service, error)
	GetServiceFunc    func(ctx context.Context, name string) (*runpb.Service, error)
	UpdateServiceFunc func(ctx context.Context, service *runpb.Service) (*runpb.Service, error)
}

func (m *MockClient) ListServices(ctx context.Context, project, region string) ([]*runpb.Service, error) {
	if m.ListServicesFunc != nil {
		return m.ListServicesFunc(ctx, project, region)
	}
	return nil, nil
}

func (m *MockClient) GetService(ctx context.Context, name string) (*runpb.Service, error) {
	if m.GetServiceFunc != nil {
		return m.GetServiceFunc(ctx, name)
	}
	return nil, nil
}

func (m *MockClient) UpdateService(ctx context.Context, service *runpb.Service) (*runpb.Service, error) {
	if m.UpdateServiceFunc != nil {
		return m.UpdateServiceFunc(ctx, service)
	}
	return nil, nil
}

func TestList(t *testing.T) {
	// Save original client and restore after test
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListServicesFunc = func(ctx context.Context, project, region string) ([]*runpb.Service, error) {
		return []*runpb.Service{
			{Name: "projects/p/locations/r/services/s1", Uri: "uri1"},
			{Name: "projects/p/locations/r/services/s2", Uri: "uri2"},
		}, nil
	}

	services, err := List("p", "r")

	assert.NoError(t, err)
	assert.Len(t, services, 2)
	assert.Equal(t, "s1", services[0].Name)
	assert.Equal(t, "s2", services[1].Name)
}

func TestUpdateScaling(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	// Setup GetService mock
	mock.GetServiceFunc = func(ctx context.Context, name string) (*runpb.Service, error) {
		return &runpb.Service{
			Name: name,
			Scaling: &runpb.ServiceScaling{
				MinInstanceCount: 0,
				MaxInstanceCount: 1,
			},
		}, nil
	}

	// Setup UpdateService mock
	mock.UpdateServiceFunc = func(ctx context.Context, service *runpb.Service) (*runpb.Service, error) {
		// Assert that scaling was updated correctly
		assert.Equal(t, int32(2), service.Scaling.MinInstanceCount)
		assert.Equal(t, int32(5), service.Scaling.MaxInstanceCount)
		return service, nil
	}

	result, err := UpdateScaling(context.Background(), "p", "r", "s1", 2, 5, 0)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "s1", result.Name)
	assert.Equal(t, int32(2), result.Scaling.MinInstances)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListServicesFunc = func(ctx context.Context, project, region string) ([]*runpb.Service, error) {
		return nil, assert.AnError
	}

	services, err := List("p", "r")
	assert.Error(t, err)
	assert.Nil(t, services)
}

func TestUpdateScaling_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	// Test GetService Error
	mock.GetServiceFunc = func(ctx context.Context, name string) (*runpb.Service, error) {
		return nil, assert.AnError
	}
	_, err := UpdateScaling(context.Background(), "p", "r", "s", 1, 2, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to get service")

	// Test UpdateService Error
	mock.GetServiceFunc = func(ctx context.Context, name string) (*runpb.Service, error) {
		return &runpb.Service{Name: name}, nil
	}
	mock.UpdateServiceFunc = func(ctx context.Context, service *runpb.Service) (*runpb.Service, error) {
		return nil, assert.AnError
	}
	_, err = UpdateScaling(context.Background(), "p", "r", "s", 1, 2, 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to update service")
}

func TestList_AllRegions(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListServicesFunc = func(ctx context.Context, project, region string) ([]*runpb.Service, error) {
		if region == "us-central1" {
			return []*runpb.Service{
				{Name: "projects/p/locations/us-central1/services/s1"},
			}, nil
		}
		return []*runpb.Service{}, nil
	}

	services, err := List("p", api_region.ALL)
	assert.NoError(t, err)
	
	// We expect at least 1 service
	found := false
	for _, s := range services {
		if s.Name == "s1" && s.Region == "us-central1" {
			found = true
			break
		}
	}
	assert.True(t, found, "Expected to find service s1 in us-central1")
}

// --- Mocks for GCPClient testing ---

type MockServicesClientWrapper struct {
	ListServicesFunc  func(ctx context.Context, req *runpb.ListServicesRequest, opts ...gax.CallOption) ServiceIteratorWrapper
	GetServiceFunc    func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error)
	UpdateServiceFunc func(ctx context.Context, req *runpb.UpdateServiceRequest, opts ...gax.CallOption) (UpdateServiceOperationWrapper, error)
	CloseFunc         func() error
}

func (m *MockServicesClientWrapper) ListServices(ctx context.Context, req *runpb.ListServicesRequest, opts ...gax.CallOption) ServiceIteratorWrapper {
	if m.ListServicesFunc != nil {
		return m.ListServicesFunc(ctx, req, opts...)
	}
	return &MockServiceIteratorWrapper{}
}

func (m *MockServicesClientWrapper) GetService(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
	if m.GetServiceFunc != nil {
		return m.GetServiceFunc(ctx, req, opts...)
	}
	return nil, nil
}

func (m *MockServicesClientWrapper) UpdateService(ctx context.Context, req *runpb.UpdateServiceRequest, opts ...gax.CallOption) (UpdateServiceOperationWrapper, error) {
	if m.UpdateServiceFunc != nil {
		return m.UpdateServiceFunc(ctx, req, opts...)
	}
	return nil, nil
}

func (m *MockServicesClientWrapper) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

type MockServiceIteratorWrapper struct {
	Items []*runpb.Service
	Index int
	Err   error
}

func (m *MockServiceIteratorWrapper) Next() (*runpb.Service, error) {
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

type MockUpdateServiceOperationWrapper struct {
	WaitFunc func(ctx context.Context, opts ...gax.CallOption) (*runpb.Service, error)
}

func (m *MockUpdateServiceOperationWrapper) Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Service, error) {
	if m.WaitFunc != nil {
		return m.WaitFunc(ctx, opts...)
	}
	return nil, nil
}

func TestGCPClient_ListServices(t *testing.T) {
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createServicesClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createServicesClient = origCreateClient
	}()
	
	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}
	
	t.Run("Success", func(t *testing.T) {
		createServicesClient = func(ctx context.Context, opts ...option.ClientOption) (ServicesClientWrapper, error) {
			return &MockServicesClientWrapper{
				ListServicesFunc: func(ctx context.Context, req *runpb.ListServicesRequest, opts ...gax.CallOption) ServiceIteratorWrapper {
					return &MockServiceIteratorWrapper{
						Items: []*runpb.Service{{Name: "s1"}},
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		
		client := &GCPClient{}
		services, err := client.ListServices(context.Background(), "p", "r")
		assert.NoError(t, err)
		assert.Len(t, services, 1)
		assert.Equal(t, "s1", services[0].Name)
	})
	
	t.Run("Auth Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("auth failed")
		}
		client := &GCPClient{}
		_, err := client.ListServices(context.Background(), "p", "r")
		assert.Error(t, err)
	})
	
	t.Run("Iterator Auth Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createServicesClient = func(ctx context.Context, opts ...option.ClientOption) (ServicesClientWrapper, error) {
			return &MockServicesClientWrapper{
				ListServicesFunc: func(ctx context.Context, req *runpb.ListServicesRequest, opts ...gax.CallOption) ServiceIteratorWrapper {
					return &MockServiceIteratorWrapper{
						Err: errors.New("Unauthenticated"),
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		client := &GCPClient{}
		_, err := client.ListServices(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed")
	})
}

func TestGCPClient_GetService(t *testing.T) {
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createServicesClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createServicesClient = origCreateClient
	}()

	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createServicesClient = func(ctx context.Context, opts ...option.ClientOption) (ServicesClientWrapper, error) {
			return &MockServicesClientWrapper{
				GetServiceFunc: func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return &runpb.Service{Name: "s1"}, nil
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		
		client := &GCPClient{}
		s, err := client.GetService(context.Background(), "s1")
		assert.NoError(t, err)
		assert.Equal(t, "s1", s.Name)
	})
	
	t.Run("Get Error", func(t *testing.T) {
		createServicesClient = func(ctx context.Context, opts ...option.ClientOption) (ServicesClientWrapper, error) {
			return &MockServicesClientWrapper{
				GetServiceFunc: func(ctx context.Context, req *runpb.GetServiceRequest, opts ...gax.CallOption) (*runpb.Service, error) {
					return nil, errors.New("get error")
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		client := &GCPClient{}
		_, err := client.GetService(context.Background(), "s1")
		assert.Error(t, err)
	})
}

func TestGCPClient_UpdateService(t *testing.T) {
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createServicesClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createServicesClient = origCreateClient
	}()

	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createServicesClient = func(ctx context.Context, opts ...option.ClientOption) (ServicesClientWrapper, error) {
			return &MockServicesClientWrapper{
				UpdateServiceFunc: func(ctx context.Context, req *runpb.UpdateServiceRequest, opts ...gax.CallOption) (UpdateServiceOperationWrapper, error) {
					return &MockUpdateServiceOperationWrapper{
						WaitFunc: func(ctx context.Context, opts ...gax.CallOption) (*runpb.Service, error) {
							return &runpb.Service{Name: "s1-updated"}, nil
						},
					}, nil
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		
		client := &GCPClient{}
		s, err := client.UpdateService(context.Background(), &runpb.Service{Name: "s1"})
		assert.NoError(t, err)
		assert.Equal(t, "s1-updated", s.Name)
	})
	
	t.Run("Update Start Error", func(t *testing.T) {
		createServicesClient = func(ctx context.Context, opts ...option.ClientOption) (ServicesClientWrapper, error) {
			return &MockServicesClientWrapper{
				UpdateServiceFunc: func(ctx context.Context, req *runpb.UpdateServiceRequest, opts ...gax.CallOption) (UpdateServiceOperationWrapper, error) {
					return nil, errors.New("update start error")
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		client := &GCPClient{}
		_, err := client.UpdateService(context.Background(), &runpb.Service{Name: "s1"})
		assert.Error(t, err)
	})
}

func TestWrappers_Delegation(t *testing.T) {
	// Expect panics because nil clients are used
	
	t.Run("GCPServicesClientWrapper", func(t *testing.T) {
		w := &GCPServicesClientWrapper{client: nil}
		assert.Panics(t, func() { w.ListServices(context.Background(), nil) })
		assert.Panics(t, func() { w.GetService(context.Background(), nil) })
		assert.Panics(t, func() { w.UpdateService(context.Background(), nil) })
		assert.Panics(t, func() { w.Close() })
	})
	
	t.Run("GCPServiceIteratorWrapper", func(t *testing.T) {
		it := &GCPServiceIteratorWrapper{it: nil}
		assert.Panics(t, func() { it.Next() })
	})
	
	t.Run("GCPUpdateServiceOperationWrapper", func(t *testing.T) {
		op := &GCPUpdateServiceOperationWrapper{op: nil}
		assert.Panics(t, func() { op.Wait(context.Background()) })
	})
}