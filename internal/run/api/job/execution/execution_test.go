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

package execution

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
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockClient is a mock implementation of high-level Client interface.
type MockClient struct {
	ListExecutionsFunc func(ctx context.Context, project, region, jobName string) ([]*runpb.Execution, error)
}

func (m *MockClient) ListExecutions(ctx context.Context, project, region, jobName string) ([]*runpb.Execution, error) {
	if m.ListExecutionsFunc != nil {
		return m.ListExecutionsFunc(ctx, project, region, jobName)
	}
	return nil, nil
}

func TestMapExecution(t *testing.T) {
	now := time.Now()
	resp := &runpb.Execution{
		Name:           "projects/my-project/locations/us-central1/executions/my-execution",
		Job:            "my-job",
		CreateTime:     timestamppb.New(now),
		StartTime:      timestamppb.New(now),
		CompletionTime: timestamppb.New(now),
		DeleteTime:     timestamppb.New(now),
		ExpireTime:     timestamppb.New(now),
		TaskCount:      10,
		SucceededCount: 8,
		FailedCount:    1,
		RunningCount:   0,
		CancelledCount: 1,
		RetriedCount:   0,
		LogUri:         "https://console.cloud.google.com/logs/viewer",
		Conditions: []*runpb.Condition{
			{
				Type:               "Completed",
				State:              runpb.Condition_CONDITION_SUCCEEDED,
				Message:            "Success",
				LastTransitionTime: timestamppb.New(now),
			},
		},
	}

	result := mapExecution(resp, "us-central1")

	assert.Equal(t, resp.Name, result.Name)
	assert.Equal(t, resp.Job, result.Job)
	assert.Equal(t, now.Unix(), result.CreateTime.Unix())
	assert.Equal(t, now.Unix(), result.StartTime.Unix())
	assert.Equal(t, now.Unix(), result.CompletionTime.Unix())
	assert.Equal(t, now.Unix(), result.DeleteTime.Unix())
	assert.Equal(t, now.Unix(), result.ExpireTime.Unix())
	assert.Equal(t, int32(10), result.TaskCount)
	assert.Equal(t, int32(8), result.SucceededCount)
	assert.Equal(t, int32(1), result.FailedCount)
	assert.Equal(t, int32(0), result.RunningCount)
	assert.Equal(t, int32(1), result.CancelledCount)
	assert.Equal(t, int32(0), result.RetriedCount)
	assert.Equal(t, "https://console.cloud.google.com/logs/viewer", result.LogURI)
	assert.Equal(t, "us-central1", result.Region)

	assert.Len(t, result.Conditions, 1)
	assert.Equal(t, "Completed", result.Conditions[0].Type)
	assert.Equal(t, "CONDITION_SUCCEEDED", result.Conditions[0].State)
	assert.Equal(t, "Success", result.Conditions[0].Message)

	assert.NotNil(t, result.TerminalCondition)
	assert.Equal(t, "Completed", result.TerminalCondition.Type)
	assert.Equal(t, "CONDITION_SUCCEEDED", result.TerminalCondition.State)
	assert.Equal(t, "Success", result.TerminalCondition.Message)
}

func TestList(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListExecutionsFunc = func(ctx context.Context, project, region, jobName string) ([]*runpb.Execution, error) {
		now := timestamppb.New(time.Now())
		return []*runpb.Execution{
			{
				Name:           "exec1",
				CreateTime:     now,
				StartTime:      now,
				CompletionTime: now,
			},
			{
				Name:           "exec2",
				CreateTime:     now,
				StartTime:      now,
				CompletionTime: now,
			},
		}, nil
	}

	executions, err := List("p", "r", "job")
	assert.NoError(t, err)
	assert.Len(t, executions, 2)
	assert.Equal(t, "exec1", executions[0].Name)
}

func TestList_Error(t *testing.T) {
	originalClient := apiClient
	defer func() { apiClient = originalClient }()

	mock := &MockClient{}
	apiClient = mock

	mock.ListExecutionsFunc = func(ctx context.Context, project, region, jobName string) ([]*runpb.Execution, error) {
		return nil, assert.AnError
	}

	executions, err := List("p", "r", "job")
	assert.Error(t, err)
	assert.Nil(t, executions)
}

// --- Mocks for GCPClient testing ---

type MockExecutionsClientWrapper struct {
	ListExecutionsFunc func(ctx context.Context, req *runpb.ListExecutionsRequest, opts ...gax.CallOption) ExecutionIteratorWrapper
	CloseFunc          func() error
}

func (m *MockExecutionsClientWrapper) ListExecutions(ctx context.Context, req *runpb.ListExecutionsRequest, opts ...gax.CallOption) ExecutionIteratorWrapper {
	if m.ListExecutionsFunc != nil {
		return m.ListExecutionsFunc(ctx, req, opts...)
	}
	return &MockExecutionIteratorWrapper{}
}

func (m *MockExecutionsClientWrapper) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

type MockExecutionIteratorWrapper struct {
	Items []*runpb.Execution
	Index int
	Err   error
}

func (m *MockExecutionIteratorWrapper) Next() (*runpb.Execution, error) {
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

func TestGCPClient_ListExecutions(t *testing.T) {
	origFindCreds := client.FindDefaultCredentials
	origCreateClient := createExecutionsClient
	defer func() {
		client.FindDefaultCredentials = origFindCreds
		createExecutionsClient = origCreateClient
	}()

	client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
		return &google.Credentials{}, nil
	}

	t.Run("Success and Client Caching", func(t *testing.T) {
		createCount := 0
		createExecutionsClient = func(ctx context.Context, opts ...option.ClientOption) (ExecutionsClientWrapper, error) {
			createCount++
			return &MockExecutionsClientWrapper{
				ListExecutionsFunc: func(ctx context.Context, req *runpb.ListExecutionsRequest, opts ...gax.CallOption) ExecutionIteratorWrapper {
					return &MockExecutionIteratorWrapper{
						Items: []*runpb.Execution{{Name: "exec1"}},
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}

		gcpClient := &GCPClient{}

		// First call should create the client
		execs, err := gcpClient.ListExecutions(context.Background(), "p", "r", "job")
		assert.NoError(t, err)
		assert.Len(t, execs, 1)
		assert.Equal(t, "exec1", execs[0].Name)
		assert.Equal(t, 1, createCount)

		// Second call should reuse cached client, not calling createExecutionsClient again
		execs2, err := gcpClient.ListExecutions(context.Background(), "p", "r", "job")
		assert.NoError(t, err)
		assert.Len(t, execs2, 1)
		assert.Equal(t, "exec1", execs2[0].Name)
		assert.Equal(t, 1, createCount)
	})

	t.Run("Auth Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("auth failed")
		}
		gcpClient := &GCPClient{}
		_, err := gcpClient.ListExecutions(context.Background(), "p", "r", "job")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to find default credentials")
	})

	t.Run("Client Creation Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createExecutionsClient = func(ctx context.Context, opts ...option.ClientOption) (ExecutionsClientWrapper, error) {
			return nil, errors.New("client error")
		}
		gcpClient := &GCPClient{}
		_, err := gcpClient.ListExecutions(context.Background(), "p", "r", "job")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "client error")
	})

	t.Run("Iterator Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createExecutionsClient = func(ctx context.Context, opts ...option.ClientOption) (ExecutionsClientWrapper, error) {
			return &MockExecutionsClientWrapper{
				ListExecutionsFunc: func(ctx context.Context, req *runpb.ListExecutionsRequest, opts ...gax.CallOption) ExecutionIteratorWrapper {
					return &MockExecutionIteratorWrapper{
						Err: errors.New("iter error"),
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		gcpClient := &GCPClient{}
		_, err := gcpClient.ListExecutions(context.Background(), "p", "r", "job")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "iter error")
	})

	t.Run("Iterator Auth Error", func(t *testing.T) {
		client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{}, nil
		}
		createExecutionsClient = func(ctx context.Context, opts ...option.ClientOption) (ExecutionsClientWrapper, error) {
			return &MockExecutionsClientWrapper{
				ListExecutionsFunc: func(ctx context.Context, req *runpb.ListExecutionsRequest, opts ...gax.CallOption) ExecutionIteratorWrapper {
					return &MockExecutionIteratorWrapper{
						Err: errors.New("Unauthenticated request"),
					}
				},
				CloseFunc: func() error { return nil },
			}, nil
		}
		gcpClient := &GCPClient{}
		_, err := gcpClient.ListExecutions(context.Background(), "p", "r", "job")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "authentication failed")
	})
}

func TestWrappers_Delegation(t *testing.T) {
	t.Run("GCPExecutionsClientWrapper", func(t *testing.T) {
		w := &GCPExecutionsClientWrapper{client: nil}
		assert.Panics(t, func() { _ = w.ListExecutions(context.Background(), nil) })
		assert.Panics(t, func() { _ = w.Close() })
	})

	t.Run("GCPExecutionIteratorWrapper", func(t *testing.T) {
		it := &GCPExecutionIteratorWrapper{it: nil}
		assert.Panics(t, func() { _, _ = it.Next() })
	})
}
