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

package revision

import (
	"fmt"
	"strings"

	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_revision "github.com/JulienBreux/run-cli/internal/run/model/service/revision"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/table"
	"github.com/dustin/go-humanize"
	"github.com/rivo/tview"
)

// ListComponent represents a revision list component.
type ListComponent struct {
	*table.Table
	app *tview.Application
}

// DetailComponent represents a revision detail component.
type DetailComponent struct {
	*tview.TextView
}

// NewListComponent creates a new revision list component.
func NewListComponent(app *tview.Application) *ListComponent {
	t := table.New(" Revisions ")
	t.SetHeadersWithExpansions(
		[]string{"NAME", "TRAFFIC", "DEPLOYED", "REVISION TAGS"},
		[]int{2, 1, 1, 2},
	)

	return &ListComponent{
		Table: t,
		app:   app,
	}
}

// NewDetailComponent creates a new revision detail component.
func NewDetailComponent() *DetailComponent {
	tv := tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	tv.SetBorder(true).SetTitle(" Revision Details ")

	return &DetailComponent{
		TextView: tv,
	}
}

// Update updates the list component with revisions.
func (c *ListComponent) Update(service *model_service.Service, revisions []model_revision.Revision) {
	c.Table.Table.Clear()
	c.SetHeadersWithExpansions(
		[]string{"NAME", "TRAFFIC", "DEPLOYED", "REVISION TAGS"},
		[]int{2, 1, 1, 2},
	)

	for i, rev := range revisions {
		row := i + 1

		traffic := "0%"
		tags := ""
		for _, ts := range service.TrafficStatuses {
			isLatestMatch := ts.Type == "TRAFFIC_TARGET_ALLOCATION_TYPE_LATEST" && rev.Name == service.LatestReadyRevision
			isNamedMatch := ts.Revision == rev.Name

			if isLatestMatch || isNamedMatch {
				if ts.Percent > 0 {
					traffic = fmt.Sprintf("%d%%", ts.Percent)
					if isLatestMatch {
						traffic += " (to latest)"
					}
				}
				if ts.Tag != "" {
					if tags != "" {
						tags += ", "
					}
					tags += ts.Tag
				}
			}
		}

		c.Table.Table.SetCell(row, 0, tview.NewTableCell(rev.Name))
		c.Table.Table.SetCell(row, 1, tview.NewTableCell(traffic))
		c.Table.Table.SetCell(row, 2, tview.NewTableCell(humanize.Time(rev.CreateTime)))
		c.Table.Table.SetCell(row, 3, tview.NewTableCell(tags))
	}

	c.Table.Table.SetTitle(fmt.Sprintf(" Revisions (%d) ", len(revisions)))
}

// Update updates the detail component with a revision.
func (c *DetailComponent) Update(rev model_revision.Revision) {
	var sb strings.Builder
	fmt.Fprintf(&sb, "[lightcyan]Name:[white] %s\n", rev.Name)
	fmt.Fprintf(&sb, "[lightcyan]Author:[white] %s\n", rev.Author)
	fmt.Fprintf(&sb, "[lightcyan]Created:[white] %s\n", rev.CreateTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintln(&sb, "")
	fmt.Fprintln(&sb, "[yellow::b]General[white::-")

	// Billing
	billing := "CPU is always allocated"
	if rev.CpuIdle {
		billing = "CPU is only allocated during request processing"
	}
	fmt.Fprintf(&sb, "[lightcyan]Billing:[white] %s\n", billing)

	// Startup CPU Boost
	startupBoost := "Disabled"
	if rev.StartupCpuBoost {
		startupBoost = "Enabled"
	}
	fmt.Fprintf(&sb, "[lightcyan]Startup CPU boost:[white] %s\n", startupBoost)

	// Concurrency
	fmt.Fprintf(&sb, "[lightcyan]Concurrency:[white] %d\n", rev.MaxInstanceRequestConcurrency)

	// Request timeout
	fmt.Fprintf(&sb, "[lightcyan]Request timeout:[white] %s\n", rev.Timeout)

	// Execution environment
	execEnv := rev.ExecutionEnvironment
	switch execEnv {
	case "EXECUTION_ENVIRONMENT_UNSPECIFIED":
		execEnv = "Default"
	case "EXECUTION_ENVIRONMENT_GEN1":
		execEnv = "First Generation"
	case "EXECUTION_ENVIRONMENT_GEN2":
		execEnv = "Second Generation"
	}
	fmt.Fprintf(&sb, "[lightcyan]Execution environment:[white] %s\n", execEnv)

	fmt.Fprintln(&sb, "")
	fmt.Fprintln(&sb, "[yellow::b]Containers[white::-")
	for i, c := range rev.Containers {
		name := c.Name
		if name == "" {
			name = fmt.Sprintf("container-%d", i+1)
		}
		fmt.Fprintf(&sb, "[lightcyan]%s[white]\n", name)
		fmt.Fprintf(&sb, "  [lightcyan]Image:[white] %s\n", c.Image)

		if len(c.Ports) > 0 {
			fmt.Fprintf(&sb, "  [lightcyan]Port:[white] %d\n", c.Ports[0].ContainerPort)
		}

		if c.Resources != nil && len(c.Resources.Limits) > 0 {
			mem := c.Resources.Limits["memory"]
			cpu := c.Resources.Limits["cpu"]
			gpu := c.Resources.Limits["nvidia.com/gpu"]

			resStr := fmt.Sprintf("%s Memory, %s CPU", mem, cpu)
			if gpu != "" {
				gpuStr := gpu + " GPU"
				if rev.Accelerator != "" {
					gpuStr += " (" + rev.Accelerator + ")"
				}
				resStr += ", " + gpuStr
			}
			fmt.Fprintf(&sb, "  [lightcyan]Resources:[white] %s\n", resStr)
		}
		fmt.Fprintln(&sb, "")
	}

	c.SetText(sb.String())
}
