package revision

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/run/apiv2/runpb"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockClient is a mock implementation of Client.
type MockClient struct {
	ListRevisionsFunc func(ctx context.Context, project, region, service string) ([]*runpb.Revision, error)
}

func (m *MockClient) ListRevisions(ctx context.Context, project, region, service string) ([]*runpb.Revision, error) {
	if m.ListRevisionsFunc != nil {
		return m.ListRevisionsFunc(ctx, project, region, service)
	}
	return nil, nil
}

func TestList(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListRevisionsFunc = func(ctx context.Context, project, region, service string) ([]*runpb.Revision, error) {
		return []*runpb.Revision{{Name: "rev1"}}, nil
	}

	revisions, err := List("p", "r", "s")
	assert.NoError(t, err)
	assert.Len(t, revisions, 1)
	assert.Equal(t, "rev1", revisions[0].Name)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListRevisionsFunc = func(ctx context.Context, project, region, service string) ([]*runpb.Revision, error) {
		return nil, assert.AnError
	}

	revisions, err := List("p", "r", "s")
	assert.Error(t, err)
	assert.Nil(t, revisions)
}

func TestMapRevision(t *testing.T) {
	now := time.Now()
	resp := &runpb.Revision{
		Name:       "projects/p/locations/l/services/s/revisions/my-rev",
		CreateTime: timestamppb.New(now),
		UpdateTime: timestamppb.New(now),
		Containers: []*runpb.Container{
			{
				Name:  "c1",
				Image: "img:latest",
				Resources: &runpb.ResourceRequirements{
					CpuIdle:         true,
					StartupCpuBoost: true,
				},
				Ports: []*runpb.ContainerPort{
					{
						ContainerPort: 8080,
					},
				},
			},
		},
		ExecutionEnvironment:          runpb.ExecutionEnvironment_EXECUTION_ENVIRONMENT_GEN2,
		MaxInstanceRequestConcurrency: 80,
		Timeout:                       durationpb.New(time.Second * 30),
		NodeSelector: &runpb.NodeSelector{
			Accelerator: "nvidia-tesla-t4",
		},
	}

	result := mapRevision(resp, "my-service")

	assert.Equal(t, "my-rev", result.Name)
	assert.Equal(t, "my-service", result.Service)
	assert.Equal(t, now.Unix(), result.CreateTime.Unix())
	
	// Containers
	assert.Len(t, result.Containers, 1)
	assert.Equal(t, "c1", result.Containers[0].Name)
	assert.True(t, result.Containers[0].Resources.CPUIdle)
	
	// Env
	assert.Equal(t, "EXECUTION_ENVIRONMENT_GEN2", result.ExecutionEnvironment)
	assert.Equal(t, int32(80), result.MaxInstanceRequestConcurrency)
	assert.Equal(t, 30*time.Second, result.Timeout)
	
	// Accelerator
	assert.Equal(t, "nvidia-tesla-t4", result.Accelerator)
	
	// Top level shortcuts
	assert.True(t, result.CpuIdle)
	assert.True(t, result.StartupCpuBoost)
}

func TestMapRevision_NilFields(t *testing.T) {
	resp := &runpb.Revision{
		Name: "projects/p/locations/l/services/s/revisions/my-rev",
	}

	result := mapRevision(resp, "my-service")

	assert.Equal(t, "my-rev", result.Name)
	assert.Empty(t, result.Containers)
	assert.False(t, result.CpuIdle)
	assert.Empty(t, result.Accelerator)
}

// --- GCPClient Tests ---

type MockRevisionsClientWrapper struct {
	ListRevisionsFunc func(ctx context.Context, req *runpb.ListRevisionsRequest, opts ...gax.CallOption) RevisionIteratorWrapper
	CloseFunc         func() error
}

func (m *MockRevisionsClientWrapper) ListRevisions(ctx context.Context, req *runpb.ListRevisionsRequest, opts ...gax.CallOption) RevisionIteratorWrapper {
	if m.ListRevisionsFunc != nil {
		return m.ListRevisionsFunc(ctx, req, opts...)
	}
	return &MockRevisionIteratorWrapper{}
}

func (m *MockRevisionsClientWrapper) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

type MockRevisionIteratorWrapper struct {
	Items []*runpb.Revision
	Index int
	Err   error
}

func (m *MockRevisionIteratorWrapper) Next() (*runpb.Revision, error) {
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

func TestGCPClient_ListRevisions(t *testing.T) {
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createRevisionsClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createRevisionsClient = origCreateClient
	}()

	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createRevisionsClient = func(ctx context.Context, opts ...option.ClientOption) (RevisionsClientWrapper, error) {
			return &MockRevisionsClientWrapper{
				ListRevisionsFunc: func(ctx context.Context, req *runpb.ListRevisionsRequest, opts ...gax.CallOption) RevisionIteratorWrapper {
					return &MockRevisionIteratorWrapper{
						Items: []*runpb.Revision{{Name: "rev1"}},
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}

		c := &GCPClient{}
		revs, err := c.ListRevisions(context.Background(), "p", "r", "s")
		assert.NoError(t, err)
		assert.Len(t, revs, 1)
		assert.Equal(t, "rev1", revs[0].Name)
	})

	t.Run("AuthError", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("auth failed")
		}
		c := &GCPClient{}
		_, err := c.ListRevisions(context.Background(), "p", "r", "s")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find default credentials")
	})

	t.Run("ClientCreationError", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createRevisionsClient = func(ctx context.Context, opts ...option.ClientOption) (RevisionsClientWrapper, error) {
			return nil, errors.New("creation failed")
		}
		c := &GCPClient{}
		_, err := c.ListRevisions(context.Background(), "p", "r", "s")
		assert.Error(t, err)
	})

	t.Run("IteratorError", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createRevisionsClient = func(ctx context.Context, opts ...option.ClientOption) (RevisionsClientWrapper, error) {
			return &MockRevisionsClientWrapper{
				ListRevisionsFunc: func(ctx context.Context, req *runpb.ListRevisionsRequest, opts ...gax.CallOption) RevisionIteratorWrapper {
					return &MockRevisionIteratorWrapper{Err: errors.New("iter failed")}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		c := &GCPClient{}
		_, err := c.ListRevisions(context.Background(), "p", "r", "s")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "iter failed")
	})
}

func TestWrappers_Delegation(t *testing.T) {
	t.Run("GCPRevisionsClientWrapper", func(t *testing.T) {
		w := &GCPRevisionsClientWrapper{client: nil}
		assert.Panics(t, func() { _ = w.ListRevisions(context.Background(), nil) })
		assert.Panics(t, func() { _ = w.Close() })
	})

	t.Run("GCPRevisionIteratorWrapper", func(t *testing.T) {
		it := &GCPRevisionIteratorWrapper{it: nil}
		assert.Panics(t, func() { _, _ = it.Next() })
	})
}