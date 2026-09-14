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
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/api/client"
	"google.golang.org/api/monitoring/v3"
	"google.golang.org/api/option"
)

// TimeSeriesClientWrapper defines the interface for querying Cloud Monitoring time series.
type TimeSeriesClientWrapper interface {
	ListTimeSeries(ctx context.Context, projectID string, filter string, startTime, endTime time.Time, aligner, alignmentPeriod, reducer string, groupBy []string) ([]*monitoring.TimeSeries, error)
}

var (
	clientCache *monitoring.Service
	clientMutex sync.Mutex

	newMonitoringService = monitoring.NewService

	createMonitoringService = func(ctx context.Context, opts ...option.ClientOption) (*monitoring.Service, error) {
		clientMutex.Lock()
		defer clientMutex.Unlock()

		if clientCache != nil {
			return clientCache, nil
		}

		creds, err := client.FindDefaultCredentials(ctx, monitoring.MonitoringReadScope)
		if err != nil {
			return nil, client.WrapError(err)
		}

		var allOpts []option.ClientOption
		allOpts = append(allOpts, option.WithCredentials(creds))
		allOpts = append(allOpts, opts...)

		svc, err := newMonitoringService(ctx, allOpts...)
		if err != nil {
			return nil, client.WrapError(err)
		}
		clientCache = svc
		return svc, nil
	}


	createTimeSeriesClient = func(ctx context.Context) (TimeSeriesClientWrapper, error) {
		svc, err := createMonitoringService(ctx)
		if err != nil {
			return nil, err
		}
		return &GCPTimeSeriesClient{service: svc}, nil
	}
)

// GCPTimeSeriesClient is the concrete implementation of TimeSeriesClientWrapper using GCP Monitoring API.
type GCPTimeSeriesClient struct {
	service *monitoring.Service
}

func (c *GCPTimeSeriesClient) ListTimeSeries(ctx context.Context, projectID string, filter string, startTime, endTime time.Time, aligner, alignmentPeriod, reducer string, groupBy []string) ([]*monitoring.TimeSeries, error) {
	name := fmt.Sprintf("projects/%s", projectID)
	req := c.service.Projects.TimeSeries.List(name).
		Filter(filter).
		IntervalStartTime(startTime.Format(time.RFC3339)).
		IntervalEndTime(endTime.Format(time.RFC3339)).
		Context(ctx)

	if aligner != "" {
		req = req.AggregationPerSeriesAligner(aligner)
	}
	if alignmentPeriod != "" {
		req = req.AggregationAlignmentPeriod(alignmentPeriod)
	}
	if reducer != "" {
		req = req.AggregationCrossSeriesReducer(reducer)
	}
	if len(groupBy) > 0 {
		req = req.AggregationGroupByFields(groupBy...)
	}

	resp, err := req.Do()
	if err != nil {
		return nil, err
	}
	return resp.TimeSeries, nil
}

// IsPermissionOrDisabledError checks whether the error indicates a disabled API or permission issue.
func IsPermissionOrDisabledError(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "monitoring.googleapis.com") ||
		strings.Contains(msg, "permissiondenied") ||
		strings.Contains(msg, "permission denied") ||
		strings.Contains(msg, "service_disabled") ||
		strings.Contains(msg, "403") ||
		strings.Contains(msg, "has not been used in project")
}
