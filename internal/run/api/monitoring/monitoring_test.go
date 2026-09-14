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

package monitoring

import (
	"context"
	"errors"
	"testing"
	"time"

	"cloud.google.com/go/logging"
	api_client "github.com/JulienBreux/run-cli/internal/run/api/client"
	model_monitoring "github.com/JulienBreux/run-cli/internal/run/model/monitoring"
	"github.com/stretchr/testify/assert"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/monitoring/v3"
	"google.golang.org/api/option"
)

type MockTimeSeriesClient struct {
	ListTimeSeriesFunc func(ctx context.Context, projectID string, filter string, startTime, endTime time.Time, aligner, alignmentPeriod, reducer string, groupBy []string) ([]*monitoring.TimeSeries, error)
}

func (m *MockTimeSeriesClient) ListTimeSeries(ctx context.Context, projectID string, filter string, startTime, endTime time.Time, aligner, alignmentPeriod, reducer string, groupBy []string) ([]*monitoring.TimeSeries, error) {
	if m.ListTimeSeriesFunc != nil {
		return m.ListTimeSeriesFunc(ctx, projectID, filter, startTime, endTime, aligner, alignmentPeriod, reducer, groupBy)
	}
	return nil, nil
}

func TestIsPermissionOrDisabledError(t *testing.T) {
	assert.False(t, IsPermissionOrDisabledError(nil))
	assert.False(t, IsPermissionOrDisabledError(errors.New("network timeout")))

	assert.True(t, IsPermissionOrDisabledError(errors.New("monitoring.googleapis.com is not enabled")))
	assert.True(t, IsPermissionOrDisabledError(errors.New("googleapi: Error 403: PermissionDenied")))
	assert.True(t, IsPermissionOrDisabledError(errors.New("service_disabled in project")))
	assert.True(t, IsPermissionOrDisabledError(errors.New("Cloud Monitoring API has not been used in project")))
}

func TestFetchMetrics_Success(t *testing.T) {
	origCreateClient := createTimeSeriesClient
	origFetchLogs := FetchRecentLogsFunc
	defer func() {
		createTimeSeriesClient = origCreateClient
		FetchRecentLogsFunc = origFetchLogs
	}()

	intVal := func(v int64) *int64 { return &v }
	doubleVal := func(v float64) *float64 { return &v }

	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339)

	mockClient := &MockTimeSeriesClient{
		ListTimeSeriesFunc: func(ctx context.Context, projectID string, filter string, startTime, endTime time.Time, aligner, alignmentPeriod, reducer string, groupBy []string) ([]*monitoring.TimeSeries, error) {
			if aligner == "ALIGN_DELTA" {
				// Requests
				return []*monitoring.TimeSeries{
					{
						Metric: &monitoring.Metric{
							Labels: map[string]string{"response_code_class": "2xx"},
						},
						Points: []*monitoring.Point{
							{
								Interval: &monitoring.TimeInterval{EndTime: nowStr},
								Value:    &monitoring.TypedValue{Int64Value: intVal(90)},
							},
						},
					},
					{
						Metric: &monitoring.Metric{
							Labels: map[string]string{"response_code_class": "5xx"},
						},
						Points: []*monitoring.Point{
							{
								Interval: &monitoring.TimeInterval{EndTime: nowStr},
								Value:    &monitoring.TypedValue{Int64Value: intVal(10)},
							},
						},
					},
				}, nil
			}
			if aligner == "ALIGN_PERCENTILE_95" {
				// Latency
				return []*monitoring.TimeSeries{
					{
						Points: []*monitoring.Point{
							{
								Interval: &monitoring.TimeInterval{EndTime: nowStr},
								Value:    &monitoring.TypedValue{DoubleValue: doubleVal(125.5)},
							},
						},
					},
				}, nil
			}
			if aligner == "ALIGN_MAX" {
				// Instances
				return []*monitoring.TimeSeries{
					{
						Metric: &monitoring.Metric{
							Labels: map[string]string{"state": "active"},
						},
						Points: []*monitoring.Point{
							{
								Interval: &monitoring.TimeInterval{EndTime: nowStr},
								Value:    &monitoring.TypedValue{Int64Value: intVal(4)},
							},
						},
					},
					{
						Metric: &monitoring.Metric{
							Labels: map[string]string{"state": "idle"},
						},
						Points: []*monitoring.Point{
							{
								Interval: &monitoring.TimeInterval{EndTime: nowStr},
								Value:    &monitoring.TypedValue{Int64Value: intVal(1)},
							},
						},
					},
				}, nil
			}
			if aligner == "ALIGN_MEAN" {
				// CPU / Mem
				return []*monitoring.TimeSeries{
					{
						Points: []*monitoring.Point{
							{
								Value: &monitoring.TypedValue{DoubleValue: doubleVal(0.42)},
							},
						},
					},
				}, nil
			}
			return nil, nil
		},
	}

	createTimeSeriesClient = func(ctx context.Context) (TimeSeriesClientWrapper, error) {
		return mockClient, nil
	}

	FetchRecentLogsFunc = func(ctx context.Context, projectID, filter string, limit int) ([]*logging.Entry, error) {
		return []*logging.Entry{
			{
				Timestamp: now,
				Severity:  logging.Error,
				Payload:   "test error payload",
			},
		}, nil
	}

	summary, err := FetchMetrics(context.Background(), "my-proj", "us-central1", "my-svc", model_monitoring.Window1h)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.False(t, summary.HasPermissionIssue)
	assert.Equal(t, int64(100), summary.TotalRequests)
	assert.Equal(t, int64(90), summary.Requests2xx)
	assert.Equal(t, int64(10), summary.Requests5xx)
	assert.InDelta(t, 10.0, summary.ErrorRate, 0.001)
	assert.InDelta(t, 125.5, summary.LatencyP95Ms, 0.001)
	assert.Equal(t, int64(4), summary.ActiveInstances)
	assert.Equal(t, int64(1), summary.IdleInstances)
	assert.InDelta(t, 42.0, summary.AvgCPUPercent, 0.001)
	assert.Len(t, summary.RecentLogs, 1)
	assert.Equal(t, "ERROR", summary.RecentLogs[0].Severity)
}

func TestFetchMetrics_PermissionOrDisabled(t *testing.T) {
	origCreateClient := createTimeSeriesClient
	origFetchLogs := FetchRecentLogsFunc
	defer func() {
		createTimeSeriesClient = origCreateClient
		FetchRecentLogsFunc = origFetchLogs
	}()

	createTimeSeriesClient = func(ctx context.Context) (TimeSeriesClientWrapper, error) {
		return nil, errors.New("googleapi: Error 403: monitoring.googleapis.com is not enabled")
	}

	FetchRecentLogsFunc = func(ctx context.Context, projectID, filter string, limit int) ([]*logging.Entry, error) {
		return nil, nil
	}

	summary, err := FetchMetrics(context.Background(), "my-proj", "us-central1", "my-svc", model_monitoring.Window1h)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.True(t, summary.HasPermissionIssue)
	assert.Contains(t, summary.PermissionErrorMessage, "gcloud services enable monitoring.googleapis.com")
}

func TestFetchMetrics_OtherClientError(t *testing.T) {
	origCreateClient := createTimeSeriesClient
	defer func() {
		createTimeSeriesClient = origCreateClient
	}()

	createTimeSeriesClient = func(ctx context.Context) (TimeSeriesClientWrapper, error) {
		return nil, errors.New("unrecoverable internal error")
	}

	summary, err := FetchMetrics(context.Background(), "my-proj", "us-central1", "my-svc", model_monitoring.Window1h)
	assert.Error(t, err)
	assert.Nil(t, summary)
}

func TestFetchMetrics_ListError_Permission(t *testing.T) {
	origCreateClient := createTimeSeriesClient
	origFetchLogs := FetchRecentLogsFunc
	defer func() {
		createTimeSeriesClient = origCreateClient
		FetchRecentLogsFunc = origFetchLogs
	}()

	mockClient := &MockTimeSeriesClient{
		ListTimeSeriesFunc: func(ctx context.Context, projectID string, filter string, startTime, endTime time.Time, aligner, alignmentPeriod, reducer string, groupBy []string) ([]*monitoring.TimeSeries, error) {
			return nil, errors.New("PermissionDenied: caller does not have permission")
		},
	}

	createTimeSeriesClient = func(ctx context.Context) (TimeSeriesClientWrapper, error) {
		return mockClient, nil
	}

	FetchRecentLogsFunc = func(ctx context.Context, projectID, filter string, limit int) ([]*logging.Entry, error) {
		return nil, nil
	}

	summary, err := FetchMetrics(context.Background(), "my-proj", "us-central1", "my-svc", model_monitoring.Window1h)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	assert.True(t, summary.HasPermissionIssue)
}

func TestFetchMetrics_ListError_Generic(t *testing.T) {
	origCreateClient := createTimeSeriesClient
	defer func() {
		createTimeSeriesClient = origCreateClient
	}()

	mockClient := &MockTimeSeriesClient{
		ListTimeSeriesFunc: func(ctx context.Context, projectID string, filter string, startTime, endTime time.Time, aligner, alignmentPeriod, reducer string, groupBy []string) ([]*monitoring.TimeSeries, error) {
			return nil, errors.New("database connection failed")
		},
	}

	createTimeSeriesClient = func(ctx context.Context) (TimeSeriesClientWrapper, error) {
		return mockClient, nil
	}

	summary, err := FetchMetrics(context.Background(), "my-proj", "us-central1", "my-svc", model_monitoring.Window1h)
	assert.Error(t, err)
	assert.Nil(t, summary)
}

func TestGetPointValue_Distribution(t *testing.T) {
	val := &monitoring.TypedValue{
		DistributionValue: &monitoring.Distribution{
			Mean: 78.4,
		},
	}
	assert.InDelta(t, 78.4, getPointValue(val), 0.001)
	assert.Equal(t, 0.0, getPointValue(nil))
}

func TestCreateMonitoringService(t *testing.T) {
	origFind := api_client.FindDefaultCredentials
	origNew := newMonitoringService
	defer func() {
		api_client.FindDefaultCredentials = origFind
		newMonitoringService = origNew
		clientCache = nil
	}()

	t.Run("FindDefaultCredentials_Error", func(t *testing.T) {
		clientCache = nil
		api_client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return nil, errors.New("creds failed")
		}

		svc, err := createMonitoringService(context.Background())
		assert.Error(t, err)
		assert.Nil(t, svc)
	})

	t.Run("Success_And_Cache", func(t *testing.T) {
		clientCache = nil
		api_client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{
				ProjectID: "test-proj",
			}, nil
		}

		callCount := 0
		dummySvc := &monitoring.Service{}
		newMonitoringService = func(ctx context.Context, opts ...option.ClientOption) (*monitoring.Service, error) {
			callCount++
			return dummySvc, nil
		}

		svc, err := createMonitoringService(context.Background())
		assert.NoError(t, err)
		assert.Same(t, dummySvc, svc)
		assert.Equal(t, 1, callCount)

		// Second call should return cached without calling newMonitoringService again
		svc2, err2 := createMonitoringService(context.Background())
		assert.NoError(t, err2)
		assert.Same(t, dummySvc, svc2)
		assert.Equal(t, 1, callCount)
	})

	t.Run("NewMonitoringService_Error", func(t *testing.T) {
		clientCache = nil
		api_client.FindDefaultCredentials = func(ctx context.Context, scopes ...string) (*google.Credentials, error) {
			return &google.Credentials{ProjectID: "test-proj"}, nil
		}
		newMonitoringService = func(ctx context.Context, opts ...option.ClientOption) (*monitoring.Service, error) {
			return nil, errors.New("service creation failed")
		}

		svc, err := createMonitoringService(context.Background())
		assert.Error(t, err)
		assert.Nil(t, svc)
	})
}

