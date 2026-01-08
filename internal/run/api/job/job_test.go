package job

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

// MockClient is a mock implementation of the Client interface (High Level).
type MockClient struct {
	ListJobsFunc func(ctx context.Context, project, region string) ([]*runpb.Job, error)
	RunJobFunc   func(ctx context.Context, name string) (*runpb.Execution, error)
}

func (m *MockClient) ListJobs(ctx context.Context, project, region string) ([]*runpb.Job, error) {
	if m.ListJobsFunc != nil {
		return m.ListJobsFunc(ctx, project, region)
	}
	return nil, nil
}

func (m *MockClient) RunJob(ctx context.Context, name string) (*runpb.Execution, error) {
	if m.RunJobFunc != nil {
		return m.RunJobFunc(ctx, name)
	}
	return nil, nil
}

func TestMapJob(t *testing.T) {
	now := time.Now()
	resp := &runpb.Job{
		Name:    "projects/my-project/locations/us-central1/jobs/my-job",
		Creator: "user@example.com",
		LatestCreatedExecution: &runpb.ExecutionReference{
			Name:       "projects/my-project/locations/us-central1/executions/my-job-exec",
			CreateTime: timestamppb.New(now),
		},
		TerminalCondition: &runpb.Condition{
			State:              runpb.Condition_CONDITION_SUCCEEDED,
			Message:            "All good",
			LastTransitionTime: timestamppb.New(now),
		},
	}

	result := mapJob(resp, "us-central1")

	assert.Equal(t, resp.Name, result.Name)
	assert.Equal(t, "user@example.com", result.Creator)
	assert.Equal(t, "us-central1", result.Region)
	
	// Execution
	assert.NotNil(t, result.LatestCreatedExecution)
	assert.Equal(t, resp.LatestCreatedExecution.Name, result.LatestCreatedExecution.Name)
	assert.Equal(t, now.Unix(), result.LatestCreatedExecution.CreateTime.Unix())
	
	// Condition
	assert.NotNil(t, result.TerminalCondition)
	assert.Equal(t, "CONDITION_SUCCEEDED", result.TerminalCondition.State)
	assert.Equal(t, "All good", result.TerminalCondition.Message)
}

func TestMapJob_NilFields(t *testing.T) {
	resp := &runpb.Job{
		Name:    "projects/my-project/locations/us-central1/jobs/my-job",
		Creator: "user@example.com",
	}

	result := mapJob(resp, "us-central1")

	assert.Equal(t, resp.Name, result.Name)
	assert.Equal(t, "user@example.com", result.Creator)
	assert.Nil(t, result.LatestCreatedExecution)
	assert.Nil(t, result.TerminalCondition)
}

func TestList(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListJobsFunc = func(ctx context.Context, project, region string) ([]*runpb.Job, error) {
		return []*runpb.Job{
			{Name: "job1"},
			{Name: "job2"},
		}, nil
	}

	jobs, err := List("p", "r")
	assert.NoError(t, err)
	assert.Len(t, jobs, 2)
	assert.Equal(t, "job1", jobs[0].Name)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListJobsFunc = func(ctx context.Context, project, region string) ([]*runpb.Job, error) {
		return nil, assert.AnError
	}

	jobs, err := List("p", "r")
	assert.Error(t, err)
	assert.Nil(t, jobs)
}

func TestExecute(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.RunJobFunc = func(ctx context.Context, name string) (*runpb.Execution, error) {
		assert.Equal(t, "projects/p/locations/r/jobs/myjob", name)
		return &runpb.Execution{Name: "exec1"}, nil
	}

	exec, err := Execute("p", "r", "myjob")
	assert.NoError(t, err)
	assert.NotNil(t, exec)
	assert.Equal(t, "exec1", exec.Name)
}

func TestExecute_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.RunJobFunc = func(ctx context.Context, name string) (*runpb.Execution, error) {
		return nil, assert.AnError
	}

	exec, err := Execute("p", "r", "myjob")
	assert.Error(t, err)
	assert.Nil(t, exec)
}

func TestList_AllRegions(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListJobsFunc = func(ctx context.Context, project, region string) ([]*runpb.Job, error) {
		if region == "us-central1" {
			return []*runpb.Job{{Name: "job-us"}}, nil
		}
		return []*runpb.Job{}, nil
	}

	jobs, err := List("p", api_region.ALL)
	assert.NoError(t, err)
	
	found := false
	for _, j := range jobs {
		if j.Name == "job-us" && j.Region == "us-central1" {
			found = true
			break
		}
	}
	assert.True(t, found)
}

// --- Mocks for GCPClient testing ---

type MockJobsClientWrapper struct {
	ListJobsFunc func(ctx context.Context, req *runpb.ListJobsRequest, opts ...gax.CallOption) JobIteratorWrapper
	RunJobFunc   func(ctx context.Context, req *runpb.RunJobRequest, opts ...gax.CallOption) (RunJobOperationWrapper, error)
	CloseFunc    func() error
}

func (m *MockJobsClientWrapper) ListJobs(ctx context.Context, req *runpb.ListJobsRequest, opts ...gax.CallOption) JobIteratorWrapper {
	if m.ListJobsFunc != nil {
		return m.ListJobsFunc(ctx, req, opts...)
	}
	return &MockJobIteratorWrapper{}
}

func (m *MockJobsClientWrapper) RunJob(ctx context.Context, req *runpb.RunJobRequest, opts ...gax.CallOption) (RunJobOperationWrapper, error) {
	if m.RunJobFunc != nil {
		return m.RunJobFunc(ctx, req, opts...)
	}
	return nil, nil
}

func (m *MockJobsClientWrapper) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

type MockJobIteratorWrapper struct {
	Items []*runpb.Job
	Index int
	Err   error
}

func (m *MockJobIteratorWrapper) Next() (*runpb.Job, error) {
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

type MockRunJobOperationWrapper struct {
	WaitFunc func(ctx context.Context, opts ...gax.CallOption) (*runpb.Execution, error)
}

func (m *MockRunJobOperationWrapper) Wait(ctx context.Context, opts ...gax.CallOption) (*runpb.Execution, error) {
	if m.WaitFunc != nil {
		return m.WaitFunc(ctx, opts...)
	}
	return nil, nil
}

func TestGCPClient_ListJobs(t *testing.T) {
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createJobsClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createJobsClient = origCreateClient
	}()

	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createJobsClient = func(ctx context.Context, opts ...option.ClientOption) (JobsClientWrapper, error) {
			return &MockJobsClientWrapper{
				ListJobsFunc: func(ctx context.Context, req *runpb.ListJobsRequest, opts ...gax.CallOption) JobIteratorWrapper {
					return &MockJobIteratorWrapper{
						Items: []*runpb.Job{{Name: "job1"}},
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}

		client := &GCPClient{}
		jobs, err := client.ListJobs(context.Background(), "p", "r")
		assert.NoError(t, err)
		assert.Len(t, jobs, 1)
		assert.Equal(t, "job1", jobs[0].Name)
	})

	t.Run("Auth Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("auth failed")
		}
		client := &GCPClient{}
		_, err := client.ListJobs(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find default credentials")
	})
	
	t.Run("Client Creation Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createJobsClient = func(ctx context.Context, opts ...option.ClientOption) (JobsClientWrapper, error) {
			return nil, errors.New("client error")
		}
		client := &GCPClient{}
		_, err := client.ListJobs(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client error")
	})

	t.Run("Iterator Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createJobsClient = func(ctx context.Context, opts ...option.ClientOption) (JobsClientWrapper, error) {
			return &MockJobsClientWrapper{
				ListJobsFunc: func(ctx context.Context, req *runpb.ListJobsRequest, opts ...gax.CallOption) JobIteratorWrapper {
					return &MockJobIteratorWrapper{
						Err: errors.New("iter error"),
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		client := &GCPClient{}
		_, err := client.ListJobs(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "iter error")
	})
	
	t.Run("Iterator Auth Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createJobsClient = func(ctx context.Context, opts ...option.ClientOption) (JobsClientWrapper, error) {
			return &MockJobsClientWrapper{
				ListJobsFunc: func(ctx context.Context, req *runpb.ListJobsRequest, opts ...gax.CallOption) JobIteratorWrapper {
					return &MockJobIteratorWrapper{
						Err: errors.New("Unauthenticated request"),
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		client := &GCPClient{}
		_, err := client.ListJobs(context.Background(), "p", "r")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed")
	})
}

func TestGCPClient_RunJob(t *testing.T) {
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createJobsClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createJobsClient = origCreateClient
	}()
	
	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success", func(t *testing.T) {
		createJobsClient = func(ctx context.Context, opts ...option.ClientOption) (JobsClientWrapper, error) {
			return &MockJobsClientWrapper{
				RunJobFunc: func(ctx context.Context, req *runpb.RunJobRequest, opts ...gax.CallOption) (RunJobOperationWrapper, error) {
					return &MockRunJobOperationWrapper{
						WaitFunc: func(ctx context.Context, opts ...gax.CallOption) (*runpb.Execution, error) {
							return &runpb.Execution{Name: "exec-1"}, nil
						},
					}, nil
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		
		client := &GCPClient{}
		exec, err := client.RunJob(context.Background(), "job1")
		assert.NoError(t, err)
		assert.Equal(t, "exec-1", exec.Name)
	})
	
	t.Run("Run Error", func(t *testing.T) {
		createJobsClient = func(ctx context.Context, opts ...option.ClientOption) (JobsClientWrapper, error) {
			return &MockJobsClientWrapper{
				RunJobFunc: func(ctx context.Context, req *runpb.RunJobRequest, opts ...gax.CallOption) (RunJobOperationWrapper, error) {
					return nil, errors.New("run failed")
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		
		client := &GCPClient{}
		_, err := client.RunJob(context.Background(), "job1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "run failed")
	})
	
	t.Run("Client Creation Error", func(t *testing.T) {
		createJobsClient = func(ctx context.Context, opts ...option.ClientOption) (JobsClientWrapper, error) {
			return nil, errors.New("client creation error")
		}
		
		client := &GCPClient{}
		_, err := client.RunJob(context.Background(), "job1")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client creation error")
	})
}

func TestWrappers_Delegation(t *testing.T) {
	// Expect panics because nil clients are used
	
	t.Run("GCPJobsClientWrapper", func(t *testing.T) {
		w := &GCPJobsClientWrapper{client: nil}
		assert.Panics(t, func() { w.ListJobs(context.Background(), nil) })
		assert.Panics(t, func() { w.RunJob(context.Background(), nil) })
		assert.Panics(t, func() { w.Close() })
	})
	
	t.Run("GCPJobIteratorWrapper", func(t *testing.T) {
		it := &GCPJobIteratorWrapper{it: nil}
		assert.Panics(t, func() { it.Next() })
	})
	
	t.Run("GCPRunJobOperationWrapper", func(t *testing.T) {
		op := &GCPRunJobOperationWrapper{op: nil}
		assert.Panics(t, func() { op.Wait(context.Background()) })
	})
}