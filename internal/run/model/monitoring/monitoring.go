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

import "time"

// Window represents the time window for metrics collection.
type Window string

const (
	Window1h  Window = "1h"
	Window6h  Window = "6h"
	Window24h Window = "24h"
)

// Duration returns the time.Duration corresponding to the window.
func (w Window) Duration() time.Duration {
	switch w {
	case Window6h:
		return 6 * time.Hour
	case Window24h:
		return 24 * time.Hour
	case Window1h:
		fallthrough
	default:
		return 1 * time.Hour
	}
}

// AlignmentPeriod returns the Cloud Monitoring alignment period string for the window.
func (w Window) AlignmentPeriod() string {
	switch w {
	case Window6h:
		return "300s"
	case Window24h:
		return "1200s"
	case Window1h:
		fallthrough
	default:
		return "60s"
	}
}

// DataPoint represents a single metric data point in time.
type DataPoint struct {
	Timestamp time.Time
	Value     float64
}

// LogEntry represents a recent log entry relevant to observability.
type LogEntry struct {
	Timestamp time.Time
	Severity  string
	Message   string
}

// MetricsSummary encapsulates all observability metrics for a Cloud Run service.
type MetricsSummary struct {
	Window           Window
	LastUpdated      time.Time
	TotalRequests    int64
	Requests2xx      int64
	Requests4xx      int64
	Requests5xx      int64
	ErrorRate        float64
	LatencyP50Ms     float64
	LatencyP95Ms     float64
	LatencyP99Ms     float64
	ActiveInstances  int64
	IdleInstances    int64
	AvgCPUPercent    float64
	AvgMemoryPercent float64

	RequestPoints  []DataPoint
	LatencyPoints  []DataPoint
	InstancePoints []DataPoint
	RecentLogs     []LogEntry

	HasPermissionIssue     bool
	PermissionErrorMessage string
}
