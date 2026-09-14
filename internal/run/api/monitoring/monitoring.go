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
	"sort"
	"strings"
	"time"

	api_log "github.com/JulienBreux/run-cli/internal/run/api/log"
	model_monitoring "github.com/JulienBreux/run-cli/internal/run/model/monitoring"
	"google.golang.org/api/monitoring/v3"
)

var (
	// FetchRecentLogsFunc is a variable for dependency injection in tests.
	FetchRecentLogsFunc = api_log.FetchRecentLogs
)

// FetchMetrics collects telemetry and metrics for a specific Cloud Run service.
func FetchMetrics(ctx context.Context, projectID, location, serviceName string, window model_monitoring.Window) (*model_monitoring.MetricsSummary, error) {
	summary := &model_monitoring.MetricsSummary{
		Window:      window,
		LastUpdated: time.Now(),
	}

	client, err := createTimeSeriesClient(ctx)
	if err != nil {
		if IsPermissionOrDisabledError(err) {
			summary.HasPermissionIssue = true
			summary.PermissionErrorMessage = fmt.Sprintf("Cloud Monitoring API is disabled or permissions are missing. Run:\ngcloud services enable monitoring.googleapis.com --project=%s", projectID)
			fetchLogs(ctx, projectID, serviceName, summary)
			return summary, nil
		}
		return nil, err
	}

	endTime := time.Now().UTC()
	startTime := endTime.Add(-window.Duration())
	alignmentPeriod := window.AlignmentPeriod()

	// 1. Request counts by status code
	reqFilter := fmt.Sprintf(`metric.type = "run.googleapis.com/request_count" AND resource.type = "cloud_run_revision" AND resource.labels.service_name = "%s"`, serviceName)
	seriesList, err := client.ListTimeSeries(ctx, projectID, reqFilter, startTime, endTime, "ALIGN_DELTA", alignmentPeriod, "REDUCE_SUM", []string{"metric.labels.response_code_class"})
	if err != nil {
		if IsPermissionOrDisabledError(err) {
			summary.HasPermissionIssue = true
			summary.PermissionErrorMessage = fmt.Sprintf("Cloud Monitoring API is disabled or permissions are missing. Run:\ngcloud services enable monitoring.googleapis.com --project=%s", projectID)
			fetchLogs(ctx, projectID, serviceName, summary)
			return summary, nil
		}
		return nil, err
	}

	processRequestSeries(seriesList, summary)

	// 2. Request latency (p95)
	latFilter := fmt.Sprintf(`metric.type = "run.googleapis.com/request_latencies" AND resource.type = "cloud_run_revision" AND resource.labels.service_name = "%s"`, serviceName)
	latSeries, err := client.ListTimeSeries(ctx, projectID, latFilter, startTime, endTime, "ALIGN_PERCENTILE_95", alignmentPeriod, "REDUCE_MEAN", nil)
	if err == nil && len(latSeries) > 0 {
		processLatencySeries(latSeries, summary)
	}

	// 3. Instance count
	instFilter := fmt.Sprintf(`metric.type = "run.googleapis.com/container/instance_count" AND resource.type = "cloud_run_revision" AND resource.labels.service_name = "%s"`, serviceName)
	instSeries, err := client.ListTimeSeries(ctx, projectID, instFilter, startTime, endTime, "ALIGN_MAX", alignmentPeriod, "REDUCE_SUM", []string{"metric.labels.state"})
	if err == nil && len(instSeries) > 0 {
		processInstanceSeries(instSeries, summary)
	}

	// 4. CPU & Memory utilization
	cpuFilter := fmt.Sprintf(`metric.type = "run.googleapis.com/container/cpu/utilizations" AND resource.type = "cloud_run_revision" AND resource.labels.service_name = "%s"`, serviceName)
	cpuSeries, err := client.ListTimeSeries(ctx, projectID, cpuFilter, startTime, endTime, "ALIGN_MEAN", alignmentPeriod, "REDUCE_MEAN", nil)
	if err == nil && len(cpuSeries) > 0 {
		summary.AvgCPUPercent = extractLatestValue(cpuSeries) * 100
	}

	memFilter := fmt.Sprintf(`metric.type = "run.googleapis.com/container/memory/utilizations" AND resource.type = "cloud_run_revision" AND resource.labels.service_name = "%s"`, serviceName)
	memSeries, err := client.ListTimeSeries(ctx, projectID, memFilter, startTime, endTime, "ALIGN_MEAN", alignmentPeriod, "REDUCE_MEAN", nil)
	if err == nil && len(memSeries) > 0 {
		summary.AvgMemoryPercent = extractLatestValue(memSeries) * 100
	}

	// 5. Recent Logs
	fetchLogs(ctx, projectID, serviceName, summary)

	return summary, nil
}

func processRequestSeries(seriesList []*monitoring.TimeSeries, summary *model_monitoring.MetricsSummary) {
	timePointsMap := make(map[string]float64)

	for _, s := range seriesList {
		codeClass := ""
		if s.Metric != nil && s.Metric.Labels != nil {
			codeClass = s.Metric.Labels["response_code_class"]
		}

		var classTotal int64
		for _, pt := range s.Points {
			val := getPointValue(pt.Value)
			classTotal += int64(val)

			if pt.Interval != nil && pt.Interval.EndTime != "" {
				timePointsMap[pt.Interval.EndTime] += val
			}
		}

		switch codeClass {
		case "2xx":
			summary.Requests2xx += classTotal
		case "4xx":
			summary.Requests4xx += classTotal
		case "5xx":
			summary.Requests5xx += classTotal
		default:
			// If not categorized, consider as 2xx or other
			summary.Requests2xx += classTotal
		}
	}

	summary.TotalRequests = summary.Requests2xx + summary.Requests4xx + summary.Requests5xx
	if summary.TotalRequests > 0 {
		summary.ErrorRate = (float64(summary.Requests5xx) / float64(summary.TotalRequests)) * 100.0
	}

	// Sort request points chronologically
	type kv struct {
		time time.Time
		val  float64
	}
	var pts []kv
	for tsStr, val := range timePointsMap {
		t, err := time.Parse(time.RFC3339, tsStr)
		if err == nil {
			pts = append(pts, kv{time: t, val: val})
		}
	}
	sort.Slice(pts, func(i, j int) bool {
		return pts[i].time.Before(pts[j].time)
	})

	for _, p := range pts {
		summary.RequestPoints = append(summary.RequestPoints, model_monitoring.DataPoint{
			Timestamp: p.time,
			Value:     p.val,
		})
	}
}

func processLatencySeries(seriesList []*monitoring.TimeSeries, summary *model_monitoring.MetricsSummary) {
	for _, s := range seriesList {
		var pts []model_monitoring.DataPoint
		for _, pt := range s.Points {
			val := getPointValue(pt.Value)
			// Latency in Cloud Run is reported in milliseconds or seconds depending on API;
			// If value < 10, likely in seconds, convert to ms
			if val > 0 && val < 5.0 {
				val *= 1000.0
			}
			t := time.Now()
			if pt.Interval != nil && pt.Interval.EndTime != "" {
				if parsed, err := time.Parse(time.RFC3339, pt.Interval.EndTime); err == nil {
					t = parsed
				}
			}
			pts = append(pts, model_monitoring.DataPoint{
				Timestamp: t,
				Value:     val,
			})
		}

		// Sort points chronologically
		sort.Slice(pts, func(i, j int) bool {
			return pts[i].Timestamp.Before(pts[j].Timestamp)
		})
		summary.LatencyPoints = pts

		if len(pts) > 0 {
			// Latest latency
			summary.LatencyP95Ms = pts[len(pts)-1].Value
		}
	}
}

func processInstanceSeries(seriesList []*monitoring.TimeSeries, summary *model_monitoring.MetricsSummary) {
	for _, s := range seriesList {
		state := ""
		if s.Metric != nil && s.Metric.Labels != nil {
			state = s.Metric.Labels["state"]
		}

		latestVal := int64(extractLatestValue([]*monitoring.TimeSeries{s}))
		switch state {
		case "active":
			summary.ActiveInstances = latestVal
		case "idle":
			summary.IdleInstances = latestVal
		}

		// Also record points
		for _, pt := range s.Points {
			t := time.Now()
			if pt.Interval != nil && pt.Interval.EndTime != "" {
				if parsed, err := time.Parse(time.RFC3339, pt.Interval.EndTime); err == nil {
					t = parsed
				}
			}
			summary.InstancePoints = append(summary.InstancePoints, model_monitoring.DataPoint{
				Timestamp: t,
				Value:     getPointValue(pt.Value),
			})
		}
	}
}

func extractLatestValue(series []*monitoring.TimeSeries) float64 {
	if len(series) == 0 || len(series[0].Points) == 0 {
		return 0
	}
	return getPointValue(series[0].Points[0].Value)
}

func getPointValue(v *monitoring.TypedValue) float64 {
	if v == nil {
		return 0
	}
	if v.DoubleValue != nil {
		return *v.DoubleValue
	}
	if v.Int64Value != nil {
		return float64(*v.Int64Value)
	}
	if v.DistributionValue != nil {
		return v.DistributionValue.Mean
	}
	return 0
}

func fetchLogs(ctx context.Context, projectID, serviceName string, summary *model_monitoring.MetricsSummary) {
	if FetchRecentLogsFunc == nil {
		return
	}
	filter := fmt.Sprintf(`resource.type = "cloud_run_revision" AND resource.labels.service_name = "%s" AND severity >= WARNING`, serviceName)
	entries, err := FetchRecentLogsFunc(ctx, projectID, filter, 10)
	if err != nil {
		return
	}

	for _, entry := range entries {
		summary.RecentLogs = append(summary.RecentLogs, model_monitoring.LogEntry{
			Timestamp: entry.Timestamp,
			Severity:  strings.ToUpper(entry.Severity.String()),
			Message:   fmt.Sprintf("%v", entry.Payload),
		})
	}
}
