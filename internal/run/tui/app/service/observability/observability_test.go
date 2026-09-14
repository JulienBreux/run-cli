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

package observability

import (
	"context"
	"errors"
	"testing"
	"time"

	model_monitoring "github.com/JulienBreux/run-cli/internal/run/model/monitoring"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNewObservabilityComponent(t *testing.T) {
	app := tview.NewApplication()
	comp := NewObservabilityComponent(app)

	assert.NotNil(t, comp)
	assert.NotNil(t, comp.View)
	assert.Equal(t, model_monitoring.Window1h, comp.window)
	assert.False(t, comp.autoRefresh)
}

func TestObservabilityComponent_SetWindow(t *testing.T) {
	app := tview.NewApplication()
	comp := NewObservabilityComponent(app)

	// Mock FetchMetricsFunc
	origFetch := FetchMetricsFunc
	defer func() { FetchMetricsFunc = origFetch }()

	called := false
	FetchMetricsFunc = func(ctx context.Context, projectID, location, serviceName string, window model_monitoring.Window) (*model_monitoring.MetricsSummary, error) {
		called = true
		assert.Equal(t, model_monitoring.Window6h, window)
		return &model_monitoring.MetricsSummary{Window: window}, nil
	}

	comp.project = "test-proj"
	comp.serviceName = "test-svc"
	comp.SetWindow(model_monitoring.Window6h)

	assert.Equal(t, model_monitoring.Window6h, comp.window)
	// Calling with same window should no-op
	called = false
	comp.SetWindow(model_monitoring.Window6h)
	assert.False(t, called)
}

func TestObservabilityComponent_ToggleAutoRefresh_And_Clear(t *testing.T) {
	app := tview.NewApplication()
	comp := NewObservabilityComponent(app)

	assert.False(t, comp.autoRefresh)
	comp.ToggleAutoRefresh()
	assert.True(t, comp.autoRefresh)
	assert.NotNil(t, comp.stopTicker)

	// Toggle off
	comp.ToggleAutoRefresh()
	assert.False(t, comp.autoRefresh)
	assert.Nil(t, comp.stopTicker)

	// Toggle on again, then Clear
	comp.ToggleAutoRefresh()
	assert.True(t, comp.autoRefresh)
	comp.Clear()
	assert.False(t, comp.autoRefresh)
	assert.Nil(t, comp.stopTicker)
}

func TestObservabilityComponent_HandleInput(t *testing.T) {
	app := tview.NewApplication()
	comp := NewObservabilityComponent(app)

	// '1' -> 1h
	comp.window = model_monitoring.Window24h
	ev := comp.HandleInput(tcell.NewEventKey(tcell.KeyRune, '1', tcell.ModNone))
	assert.Nil(t, ev)
	assert.Equal(t, model_monitoring.Window1h, comp.window)

	// '6' -> 6h
	ev = comp.HandleInput(tcell.NewEventKey(tcell.KeyRune, '6', tcell.ModNone))
	assert.Nil(t, ev)
	assert.Equal(t, model_monitoring.Window6h, comp.window)

	// 'd' -> 24h
	ev = comp.HandleInput(tcell.NewEventKey(tcell.KeyRune, 'd', tcell.ModNone))
	assert.Nil(t, ev)
	assert.Equal(t, model_monitoring.Window24h, comp.window)

	// '2' -> 24h
	comp.window = model_monitoring.Window1h
	ev = comp.HandleInput(tcell.NewEventKey(tcell.KeyRune, '2', tcell.ModNone))
	assert.Nil(t, ev)
	assert.Equal(t, model_monitoring.Window24h, comp.window)

	// 'a' -> Toggle auto refresh
	ev = comp.HandleInput(tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone))
	assert.Nil(t, ev)
	assert.True(t, comp.autoRefresh)
	comp.ToggleAutoRefresh() // turn off

	// 'r' -> Reload
	ev = comp.HandleInput(tcell.NewEventKey(tcell.KeyRune, 'r', tcell.ModNone))
	assert.Nil(t, ev)

	// Unhandled key passes through
	ev = comp.HandleInput(tcell.NewEventKey(tcell.KeyRune, 'z', tcell.ModNone))
	assert.NotNil(t, ev)
	assert.Equal(t, 'z', ev.Rune())
}

func TestObservabilityComponent_RenderSummary_WithPermissionIssue(t *testing.T) {
	app := tview.NewApplication()
	comp := NewObservabilityComponent(app)

	summary := &model_monitoring.MetricsSummary{
		HasPermissionIssue:     true,
		PermissionErrorMessage: "Please enable monitoring.googleapis.com",
		RecentLogs: []model_monitoring.LogEntry{
			{Timestamp: time.Now(), Severity: "WARN", Message: "Warning event"},
		},
	}

	comp.renderSummary(summary)
	assert.Contains(t, comp.bannerView.GetText(true), "Cloud Monitoring Unavailable")
	assert.Contains(t, comp.logsView.GetText(true), "Warning event")
}

func TestObservabilityComponent_RenderSummary_Healthy(t *testing.T) {
	app := tview.NewApplication()
	comp := NewObservabilityComponent(app)

	now := time.Now()
	summary := &model_monitoring.MetricsSummary{
		Window:        model_monitoring.Window1h,
		LastUpdated:   now,
		TotalRequests: 54000,
		Requests2xx:   53000,
		Requests4xx:   800,
		Requests5xx:   200,
		ErrorRate:     0.37,
		LatencyP95Ms:  125.4,
		ActiveInstances: 3,
		IdleInstances:   1,
		AvgCPUPercent:   45.2,
		AvgMemoryPercent: 62.0,
		RequestPoints: []model_monitoring.DataPoint{
			{Timestamp: now.Add(-2 * time.Minute), Value: 100},
			{Timestamp: now.Add(-1 * time.Minute), Value: 250},
			{Timestamp: now, Value: 300},
		},
		LatencyPoints: []model_monitoring.DataPoint{
			{Timestamp: now, Value: 125.4},
		},
		InstancePoints: []model_monitoring.DataPoint{
			{Timestamp: now, Value: 3},
		},
		RecentLogs: []model_monitoring.LogEntry{
			{Timestamp: now, Severity: "ERROR", Message: "Error in request handler"},
		},
	}

	comp.renderSummary(summary)

	assert.Contains(t, comp.requestsCard.GetText(true), "54,000")
	assert.Contains(t, comp.errorRateCard.GetText(true), "0.37%")
	assert.Contains(t, comp.latencyCard.GetText(true), "125.4 ms")
	assert.Contains(t, comp.instancesCard.GetText(true), "3 active")
	assert.Contains(t, comp.telemetryView.GetText(true), "Requests:")
	assert.Contains(t, comp.logsView.GetText(true), "Error in request handler")

	// Test nil summary
	comp.renderSummary(nil)
	assert.Contains(t, comp.requestsCard.GetText(true), "—")

	// Test render error
	comp.renderError(errors.New("custom error"))
	assert.Contains(t, comp.telemetryView.GetText(true), "Failed to load metrics")
}

func TestObservabilityComponent_Reload_EmptyService(t *testing.T) {
	app := tview.NewApplication()
	comp := NewObservabilityComponent(app)

	called := false
	comp.Reload(func(err error) {
		called = true
		assert.NoError(t, err)
	})
	assert.True(t, called)
}
