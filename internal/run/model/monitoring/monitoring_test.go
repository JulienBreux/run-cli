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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWindow_Duration(t *testing.T) {
	tests := []struct {
		window   Window
		expected time.Duration
	}{
		{Window1h, 1 * time.Hour},
		{Window6h, 6 * time.Hour},
		{Window24h, 24 * time.Hour},
		{Window("unknown"), 1 * time.Hour},
	}

	for _, tt := range tests {
		t.Run(string(tt.window), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.window.Duration())
		})
	}
}

func TestWindow_AlignmentPeriod(t *testing.T) {
	tests := []struct {
		window   Window
		expected string
	}{
		{Window1h, "60s"},
		{Window6h, "300s"},
		{Window24h, "1200s"},
		{Window("unknown"), "60s"},
	}

	for _, tt := range tests {
		t.Run(string(tt.window), func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.window.AlignmentPeriod())
		})
	}
}

func TestMetricsSummary_Initialization(t *testing.T) {
	now := time.Now()
	summary := MetricsSummary{
		Window:        Window1h,
		LastUpdated:   now,
		TotalRequests: 1000,
		Requests2xx:   950,
		Requests4xx:   40,
		Requests5xx:   10,
		ErrorRate:     1.0,
		LatencyP50Ms:  45.5,
		LatencyP95Ms:  120.0,
		LatencyP99Ms:  350.0,
		ActiveInstances: 3,
		IdleInstances:   1,
		AvgCPUPercent:   35.2,
		AvgMemoryPercent: 60.1,
		RequestPoints: []DataPoint{
			{Timestamp: now.Add(-time.Minute), Value: 50},
			{Timestamp: now, Value: 60},
		},
		RecentLogs: []LogEntry{
			{Timestamp: now, Severity: "ERROR", Message: "Something failed"},
		},
	}

	assert.Equal(t, Window1h, summary.Window)
	assert.Equal(t, int64(1000), summary.TotalRequests)
	assert.Equal(t, 1.0, summary.ErrorRate)
	assert.Len(t, summary.RequestPoints, 2)
	assert.Len(t, summary.RecentLogs, 1)
}
