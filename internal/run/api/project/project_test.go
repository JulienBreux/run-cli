package project

import (
	"context"
	"errors"
	"testing"

	"cloud.google.com/go/resourcemanager/apiv3/resourcemanagerpb"
	"github.com/JulienBreux/run-cli/internal/run/api/client"
	model "github.com/JulienBreux/run-cli/internal/run/model/common/project"
	"github.com/googleapis/gax-go/v2"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// MockClient is a mock implementation of Client.
type MockClient struct {
	ListProjectsFunc func(ctx context.Context) ([]model.Project, error)
}

func (m *MockClient) ListProjects(ctx context.Context) ([]model.Project, error) {
	return m.ListProjectsFunc(ctx)
}

func TestList(t *testing.T) {
	// Backup original client and restore after test
	origClient := apiClient
	defer func() { apiClient = origClient }()

	t.Run("success", func(t *testing.T) {
		expectedProjects := []model.Project{
			{Name: "project-1", Number: 123},
			{Name: "project-2", Number: 456},
		}

		apiClient = &MockClient{
			ListProjectsFunc: func(ctx context.Context) ([]model.Project, error) {
				return expectedProjects, nil
			},
		}

		projects, err := List()
		assert.NoError(t, err)
		assert.Equal(t, expectedProjects, projects)
	})

	t.Run("error", func(t *testing.T) {
		expectedError := errors.New("something went wrong")

		apiClient = &MockClient{
			ListProjectsFunc: func(ctx context.Context) ([]model.Project, error) {
				return nil, expectedError
			},
		}

		projects, err := List()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, projects)
	})
}

// Mocks for GCPClient testing

type MockProjectsClientWrapper struct {
	SearchProjectsFunc func(ctx context.Context, req *resourcemanagerpb.SearchProjectsRequest, opts ...gax.CallOption) ProjectIteratorWrapper
	CloseFunc          func() error
}

func (m *MockProjectsClientWrapper) SearchProjects(ctx context.Context, req *resourcemanagerpb.SearchProjectsRequest, opts ...gax.CallOption) ProjectIteratorWrapper {
	if m.SearchProjectsFunc != nil {
		return m.SearchProjectsFunc(ctx, req, opts...)
	}
	return &MockProjectIteratorWrapper{}
}

func (m *MockProjectsClientWrapper) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

type MockProjectIteratorWrapper struct {
	Items []*resourcemanagerpb.Project
	Index int
	Err   error
}

func (m *MockProjectIteratorWrapper) Next() (*resourcemanagerpb.Project, error) {
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

func TestGCPClient_ListProjects(t *testing.T) {
	// Mock dependencies
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createProjectsClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createProjectsClient = origCreateClient
	}()

	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createProjectsClient = func(ctx context.Context, opts ...option.ClientOption) (ProjectsClientWrapper, error) {
			return &MockProjectsClientWrapper{
				SearchProjectsFunc: func(ctx context.Context, req *resourcemanagerpb.SearchProjectsRequest, opts ...gax.CallOption) ProjectIteratorWrapper {
					return &MockProjectIteratorWrapper{
						Items: []*resourcemanagerpb.Project{
							{ProjectId: "p1", Name: "projects/1"},
							{ProjectId: "p2", Name: "projects/2"},
						},
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}

		client := &GCPClient{}
		projects, err := client.ListProjects(context.Background())
		assert.NoError(t, err)
		assert.Len(t, projects, 2)
		assert.Equal(t, "p1", projects[0].Name)
		assert.Equal(t, 1, projects[0].Number)
	})

	t.Run("Auth Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("auth failed")
		}

		client := &GCPClient{}
		projects, err := client.ListProjects(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find default credentials")
		assert.Nil(t, projects)
	})

	t.Run("Client Creation Error", func(t *testing.T) {
		// Restore auth mock
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createProjectsClient = func(ctx context.Context, opts ...option.ClientOption) (ProjectsClientWrapper, error) {
			return nil, errors.New("client creation failed")
		}

		client := &GCPClient{}
		projects, err := client.ListProjects(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client creation failed")
		assert.Nil(t, projects)
	})

	t.Run("Iterator Error", func(t *testing.T) {
		createProjectsClient = func(ctx context.Context, opts ...option.ClientOption) (ProjectsClientWrapper, error) {
			return &MockProjectsClientWrapper{
				SearchProjectsFunc: func(ctx context.Context, req *resourcemanagerpb.SearchProjectsRequest, opts ...gax.CallOption) ProjectIteratorWrapper {
					return &MockProjectIteratorWrapper{
						Err: errors.New("iterator error"),
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}

		client := &GCPClient{}
		projects, err := client.ListProjects(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "iterator error")
		assert.Nil(t, projects)
	})

	t.Run("Auth Error In Loop", func(t *testing.T) {
		createProjectsClient = func(ctx context.Context, opts ...option.ClientOption) (ProjectsClientWrapper, error) {
			return &MockProjectsClientWrapper{
				SearchProjectsFunc: func(ctx context.Context, req *resourcemanagerpb.SearchProjectsRequest, opts ...gax.CallOption) ProjectIteratorWrapper {
					return &MockProjectIteratorWrapper{
						Err: errors.New("something Unauthenticated here"),
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}

		client := &GCPClient{}
		projects, err := client.ListProjects(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed")
		assert.Contains(t, err.Error(), "Unauthenticated")
		assert.Nil(t, projects)
	})
}

func TestMapProject(t *testing.T) {
	resp := &resourcemanagerpb.Project{
		ProjectId: "my-project",
		Name:      "projects/123456",
	}

	result := mapProject(resp)

		assert.Equal(t, "my-project", result.Name)

		assert.Equal(t, 123456, result.Number)

	}

	

	func TestWrappers_Delegation(t *testing.T) {

		// This test exercises the wrapper methods to ensure coverage.

		// Since we can't easily mock the underlying GCP client, we expect panics when calling methods on nil clients.

		// This confirms the wrappers are attempting to delegate.

	

		t.Run("GCPProjectsClientWrapper", func(t *testing.T) {

			w := &GCPProjectsClientWrapper{client: nil} // Nil client

	

			assert.Panics(t, func() {

				w.SearchProjects(context.Background(), nil)

			})

	

			assert.Panics(t, func() {

				_ = w.Close()

			})

		})

	

		t.Run("GCPProjectIteratorWrapper", func(t *testing.T) {

			it := &GCPProjectIteratorWrapper{it: nil} // Nil iterator

	

			assert.Panics(t, func() {

				_, _ = it.Next()

			})

		})

	}

	