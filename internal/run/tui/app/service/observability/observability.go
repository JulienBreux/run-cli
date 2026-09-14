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
	"fmt"
	"strings"
	"sync"
	"time"

	api_monitoring "github.com/JulienBreux/run-cli/internal/run/api/monitoring"
	model_monitoring "github.com/JulienBreux/run-cli/internal/run/model/monitoring"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FetchMetricsFunc allows dependency injection for tests.
var FetchMetricsFunc = api_monitoring.FetchMetrics

// ObservabilityComponent represents the live metrics and observability dashboard for a service.
type ObservabilityComponent struct {
	View *tview.Flex
	app  *tview.Application

	controlsView  *tview.TextView
	requestsCard  *tview.TextView
	errorRateCard *tview.TextView
	latencyCard   *tview.TextView
	instancesCard *tview.TextView
	telemetryView *tview.TextView
	logsView      *tview.TextView
	bannerView    *tview.TextView

	cardsFlex   *tview.Flex
	contentFlex *tview.Flex

	project     string
	region      string
	serviceName string
	window      model_monitoring.Window
	autoRefresh bool
	stopTicker  chan struct{}
	tickerMutex sync.Mutex

	summary   *model_monitoring.MetricsSummary
	isLoading bool
}

// NewObservabilityComponent initializes a new observability dashboard component.
func NewObservabilityComponent(app *tview.Application) *ObservabilityComponent {
	c := &ObservabilityComponent{
		app:    app,
		window: model_monitoring.Window1h,
	}

	c.controlsView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignLeft)

	c.requestsCard = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	c.requestsCard.SetBorder(true).SetTitle(" Total Requests ")

	c.errorRateCard = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	c.errorRateCard.SetBorder(true).SetTitle(" 5xx Error Rate ")

	c.latencyCard = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	c.latencyCard.SetBorder(true).SetTitle(" Latency (p95) ")

	c.instancesCard = tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter)
	c.instancesCard.SetBorder(true).SetTitle(" Instances & Usage ")

	c.cardsFlex = tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(c.requestsCard, 0, 1, false).
		AddItem(c.errorRateCard, 0, 1, false).
		AddItem(c.latencyCard, 0, 1, false).
		AddItem(c.instancesCard, 0, 1, false)

	c.telemetryView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	c.telemetryView.SetBorder(true).SetTitle(" Metrics & Trends ")

	c.logsView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	c.logsView.SetBorder(true).SetTitle(" Recent Error & Warning Logs ")

	c.bannerView = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)
	c.bannerView.SetBorder(true).SetTitle(" Cloud Monitoring Notice ")

	c.contentFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(c.cardsFlex, 5, 0, false).
		AddItem(c.telemetryView, 8, 0, false).
		AddItem(c.logsView, 0, 1, false)

	c.View = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(c.controlsView, 1, 0, false).
		AddItem(c.contentFlex, 0, 1, true)

	c.updateControls()
	c.renderEmpty()

	c.View.SetInputCapture(c.HandleInput)

	return c
}

// SetService sets the active service and reloads metrics.
func (c *ObservabilityComponent) SetService(project, region, serviceName string) {
	c.project = project
	c.region = region
	c.serviceName = serviceName
	c.Reload(nil)
}

// SetWindow changes the time window (1h, 6h, 24h) and reloads metrics.
func (c *ObservabilityComponent) SetWindow(w model_monitoring.Window) {
	if c.window == w {
		return
	}
	c.window = w
	c.updateControls()
	c.Reload(nil)
}

// ToggleAutoRefresh toggles the 30-second periodic refresh.
func (c *ObservabilityComponent) ToggleAutoRefresh() {
	c.tickerMutex.Lock()
	defer c.tickerMutex.Unlock()

	c.autoRefresh = !c.autoRefresh
	c.updateControls()

	if c.autoRefresh {
		c.stopTicker = make(chan struct{})
		go func(stopChan chan struct{}) {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-stopChan:
					return
				case <-ticker.C:
					c.Reload(nil)
				}
			}
		}(c.stopTicker)
	} else if c.stopTicker != nil {
		close(c.stopTicker)
		c.stopTicker = nil
	}
}

// Clear stops background timers and clears component state.
func (c *ObservabilityComponent) Clear() {
	c.tickerMutex.Lock()
	if c.autoRefresh && c.stopTicker != nil {
		close(c.stopTicker)
		c.stopTicker = nil
		c.autoRefresh = false
	}
	c.tickerMutex.Unlock()

	c.summary = nil
	c.isLoading = false
	c.renderEmpty()
	c.updateControls()
}

// Reload fetches the latest metrics asynchronously.
func (c *ObservabilityComponent) Reload(onComplete func(error)) {
	if c.serviceName == "" || c.project == "" {
		if onComplete != nil {
			onComplete(nil)
		}
		return
	}

	c.isLoading = true
	c.updateControls()

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		summary, err := FetchMetricsFunc(ctx, c.project, c.region, c.serviceName, c.window)

		c.app.QueueUpdateDraw(func() {
			c.isLoading = false
			c.summary = summary
			c.updateControls()

			if err != nil {
				c.renderError(err)
				if onComplete != nil {
					onComplete(err)
				}
				return
			}

			c.renderSummary(summary)
			if onComplete != nil {
				onComplete(nil)
			}
		})
	}()
}

// HandleInput processes keyboard events for observability shortcuts.
func (c *ObservabilityComponent) HandleInput(event *tcell.EventKey) *tcell.EventKey {
	switch event.Rune() {
	case '1':
		c.SetWindow(model_monitoring.Window1h)
		return nil
	case '6':
		c.SetWindow(model_monitoring.Window6h)
		return nil
	case 'd', '2':
		c.SetWindow(model_monitoring.Window24h)
		return nil
	case 'r':
		c.Reload(nil)
		return nil
	case 'a':
		c.ToggleAutoRefresh()
		return nil
	}
	return event
}

func (c *ObservabilityComponent) updateControls() {
	var sb strings.Builder

	// Window buttons
	sb.WriteString(" [yellow::b]Window:[white::-] ")
	for _, w := range []model_monitoring.Window{model_monitoring.Window1h, model_monitoring.Window6h, model_monitoring.Window24h} {
		if c.window == w {
			fmt.Fprintf(&sb, "[black:dodgerblue] %s [-:-] ", w)
		} else {
			fmt.Fprintf(&sb, "[gray] %s [-] ", w)
		}
	}

	// Auto-refresh badge
	sb.WriteString("  [yellow::b]Auto-Refresh:[white::-] ")
	if c.autoRefresh {
		sb.WriteString("[black:green] ON (30s) [-:-] ")
	} else {
		sb.WriteString("[gray] OFF [-] ")
	}

	// Last update or loading status
	if c.isLoading {
		sb.WriteString("  [yellow::b]⟳ Loading metrics...[-:-]")
	} else if c.summary != nil && !c.summary.LastUpdated.IsZero() {
		fmt.Fprintf(&sb, "  [gray]Updated: %s[-]", c.summary.LastUpdated.Format("15:04:05"))
	}

	sb.WriteString("  [darkgray](Keys: [white]1[darkgray]:1h  [white]6[darkgray]:6h  [white]d[darkgray]:24h  [white]a[darkgray]:auto  [white]r[darkgray]:refresh)[-]")

	c.controlsView.SetText(sb.String())
}

func (c *ObservabilityComponent) renderEmpty() {
	c.requestsCard.SetText("\n[gray]—[-]")
	c.errorRateCard.SetText("\n[gray]—[-]")
	c.latencyCard.SetText("\n[gray]—[-]")
	c.instancesCard.SetText("\n[gray]—[-]")
	c.telemetryView.SetText("\n  [gray]No metrics loaded yet. Press [white]r[gray] to refresh.[-]")
	c.logsView.SetText("\n  [gray]No recent error logs.[-]")
}

func (c *ObservabilityComponent) renderError(err error) {
	c.telemetryView.SetText(fmt.Sprintf("\n  [red]Failed to load metrics: %v[-]", err))
}

func (c *ObservabilityComponent) renderSummary(s *model_monitoring.MetricsSummary) {
	if s == nil {
		c.renderEmpty()
		return
	}

	// If there's a permission issue or Cloud Monitoring is disabled
	if s.HasPermissionIssue {
		c.showPermissionBanner(s.PermissionErrorMessage)
		// Still show any logs that were fetched
		c.renderLogs(s.RecentLogs)
		return
	}

	c.hidePermissionBanner()

	// 1. Requests Card
	c.requestsCard.SetText(fmt.Sprintf("\n[white::b]%s[-:-]\n[lightgreen]2xx: %s[white] | [yellow]4xx: %s[white] | [red]5xx: %s",
		humanize.Comma(s.TotalRequests),
		humanize.Comma(s.Requests2xx),
		humanize.Comma(s.Requests4xx),
		humanize.Comma(s.Requests5xx),
	))

	// 2. Error Rate Card
	errColor := "green"
	if s.ErrorRate >= 5.0 {
		errColor = "red"
	} else if s.ErrorRate > 0.0 {
		errColor = "yellow"
	}
	c.errorRateCard.SetText(fmt.Sprintf("\n[%s::b]%.2f%%[-:-]\n[gray]%s 5xx errors[-]",
		errColor,
		s.ErrorRate,
		humanize.Comma(s.Requests5xx),
	))

	// 3. Latency Card
	latColor := "green"
	if s.LatencyP95Ms >= 1000.0 {
		latColor = "red"
	} else if s.LatencyP95Ms >= 300.0 {
		latColor = "yellow"
	}
	c.latencyCard.SetText(fmt.Sprintf("\n[%s::b]%.1f ms[-:-]\n[gray]p95 request latency[-]",
		latColor,
		s.LatencyP95Ms,
	))

	// 4. Instances & Usage Card
	c.instancesCard.SetText(fmt.Sprintf("\n[white::b]%d active[-:-] [gray]/ %d idle[-]\n[lightcyan]CPU: %.1f%%[white] | [lightcyan]Mem: %.1f%%",
		s.ActiveInstances,
		s.IdleInstances,
		s.AvgCPUPercent,
		s.AvgMemoryPercent,
	))

	// 5. Telemetry & Sparklines
	var sb strings.Builder
	reqVals := extractValues(s.RequestPoints)
	latVals := extractValues(s.LatencyPoints)
	instVals := extractValues(s.InstancePoints)

	fmt.Fprintf(&sb, " [lightcyan]Requests:  [white]%s  [gray](%s data points, window: %s)[-]\n",
		RenderSparkline(reqVals, 50),
		humanize.Comma(int64(len(reqVals))),
		s.Window,
	)
	fmt.Fprintf(&sb, " [lightcyan]Latency:   [white]%s  [gray](p95: %.1f ms)[-]\n",
		RenderColoredSparkline(latVals, 50, 300.0, 1000.0),
		s.LatencyP95Ms,
	)
	fmt.Fprintf(&sb, " [lightcyan]Instances: [white]%s  [gray](Current: %d active)[-]\n",
		RenderSparkline(instVals, 50),
		s.ActiveInstances,
	)
	fmt.Fprintf(&sb, " [lightcyan]Resources: [gray]Avg CPU: [white]%.1f%%[gray]  |  Avg Memory: [white]%.1f%%[-]",
		s.AvgCPUPercent,
		s.AvgMemoryPercent,
	)
	c.telemetryView.SetText(sb.String())

	// 6. Recent Logs
	c.renderLogs(s.RecentLogs)
}

func (c *ObservabilityComponent) renderLogs(logs []model_monitoring.LogEntry) {
	if len(logs) == 0 {
		c.logsView.SetText("\n  [green]✓ No recent warning or error logs detected in this service.[-]")
		return
	}

	var sb strings.Builder
	for _, entry := range logs {
		sevBadge := "[red::b]ERROR[-:-]"
		if strings.Contains(strings.ToUpper(entry.Severity), "WARN") {
			sevBadge = "[yellow::b]WARN [-:-]"
		}
		fmt.Fprintf(&sb, " [gray]%s[-]  %s  [white]%s[-]\n",
			entry.Timestamp.Format("15:04:05"),
			sevBadge,
			entry.Message,
		)
	}
	c.logsView.SetText(sb.String())
}

func (c *ObservabilityComponent) showPermissionBanner(message string) {
	c.bannerView.SetText(fmt.Sprintf("\n[yellow::b]⚠️  Cloud Monitoring Unavailable[-:-]\n\n[white]%s[-]\n\n[darkgray]Press [white]r[darkgray] to retry once enabled.[-]", message))

	c.contentFlex.Clear()
	c.contentFlex.
		AddItem(c.bannerView, 9, 0, false).
		AddItem(c.logsView, 0, 1, false)
}

func (c *ObservabilityComponent) hidePermissionBanner() {
	c.contentFlex.Clear()
	c.contentFlex.
		AddItem(c.cardsFlex, 5, 0, false).
		AddItem(c.telemetryView, 8, 0, false).
		AddItem(c.logsView, 0, 1, false)
}

func extractValues(points []model_monitoring.DataPoint) []float64 {
	vals := make([]float64, len(points))
	for i, p := range points {
		vals[i] = p.Value
	}
	return vals
}
