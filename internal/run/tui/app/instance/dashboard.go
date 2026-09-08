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

package instance

import (
	"bytes"
	"fmt"
	"strings"

	api_instance "github.com/JulienBreux/run-cli/internal/run/api/instance"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model "github.com/JulienBreux/run-cli/internal/run/model/instance"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/shortcut"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"sigs.k8s.io/yaml"
)

const (
	DASHBOARD_PAGE_ID = "instance-dashboard"

	yamlHighlightLexer     = "yaml"
	yamlHighlightStyle     = "solarized-dark"
	yamlHighlightFormatter = "terminal256"
)

var (
	dashboardFlex      *tview.Flex
	dashboardHeader    *tview.TextView
	dashboardTabs      *tview.TextView
	dashboardPages     *tview.Pages
	dashboardInstance  *model.Instance

	overviewDetail   *tview.TextView
	containersDetail *tview.TextView
	conditionsDetail *tview.TextView
	yamlDetail       *tview.TextView

	activeTab = 0
	tabs      = []string{"Overview", "Containers", "Conditions", "YAML"}

	GetInstanceFunc = api_instance.Get
	OnDashboardBack func()
)

// Dashboard returns the dashboard primitive.
func Dashboard(app *tview.Application) *tview.Flex {
	dashboardHeader = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	dashboardTabs = tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false)

	dashboardPages = tview.NewPages()

	// Overview Tab
	dashboardPages.AddPage(tabs[0], buildOverviewTab(), true, true)
	// Containers Tab
	dashboardPages.AddPage(tabs[1], buildContainersTab(), true, false)
	// Conditions Tab
	dashboardPages.AddPage(tabs[2], buildConditionsTab(), true, false)
	// YAML Tab
	dashboardPages.AddPage(tabs[3], buildYAMLTab(), true, false)

	dashboardFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(dashboardHeader, 1, 0, false).
		AddItem(dashboardTabs, 1, 0, false).
		AddItem(dashboardPages, 0, 1, true)

	dashboardFlex.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyRight {
			activeTab = (activeTab + 1) % len(tabs)
			updateTabs()
			return nil
		}
		if event.Key() == tcell.KeyBacktab || event.Key() == tcell.KeyLeft {
			activeTab = (activeTab - 1 + len(tabs)) % len(tabs)
			updateTabs()
			return nil
		}
		if event.Key() == tcell.KeyEscape {
			if OnDashboardBack != nil {
				OnDashboardBack()
			}
			return nil
		}
		return event
	})

	updateTabs()
	return dashboardFlex
}

func buildOverviewTab() tview.Primitive {
	overviewDetail = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	overviewDetail.SetBorder(true).SetTitle(" Overview ")
	return overviewDetail
}

func buildContainersTab() tview.Primitive {
	containersDetail = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	containersDetail.SetBorder(true).SetTitle(" Containers ")
	return containersDetail
}

func buildConditionsTab() tview.Primitive {
	conditionsDetail = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	conditionsDetail.SetBorder(true).SetTitle(" Conditions & Networking ")
	return conditionsDetail
}

func buildYAMLTab() tview.Primitive {
	yamlDetail = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(false)
	yamlDetail.SetBorder(true).SetTitle(" Raw YAML ")
	return yamlDetail
}

func updateTabs() {
	dashboardTabs.Clear()
	for i, tab := range tabs {
		if i == activeTab {
			_, _ = fmt.Fprintf(dashboardTabs, `["%s"][black:lightcyan] %s [white:-]`, tab, tab)
		} else {
			_, _ = fmt.Fprintf(dashboardTabs, `["%s"] %s `, tab, tab)
		}
	}
	dashboardPages.SwitchToPage(tabs[activeTab])
}

func updateOverviewTab() {
	if dashboardInstance == nil {
		overviewDetail.SetText("No instance details available")
		return
	}

	inst := dashboardInstance
	var sb strings.Builder

	fmt.Fprintln(&sb, "[yellow::b]Resource Identity[white::-]")
	fmt.Fprintf(&sb, "  [lightcyan]Name:[white]             %s\n", inst.Name)
	fmt.Fprintf(&sb, "  [lightcyan]ID:[white]               %s\n", inst.ID())
	fmt.Fprintf(&sb, "  [lightcyan]Region:[white]           %s\n", inst.Region)
	fmt.Fprintf(&sb, "  [lightcyan]Project:[white]          %s\n", inst.Project)
	if inst.UID != "" {
		fmt.Fprintf(&sb, "  [lightcyan]UID:[white]              %s\n", inst.UID)
	}
	if inst.Generation > 0 {
		fmt.Fprintf(&sb, "  [lightcyan]Generation:[white]       %d\n", inst.Generation)
	}
	fmt.Fprintln(&sb, "")

	fmt.Fprintln(&sb, "[yellow::b]Status & Lifecycle[white::-]")
	fmt.Fprintf(&sb, "  [lightcyan]Health Status:[white]    %s\n", inst.Status())
	if inst.Reconciling {
		fmt.Fprintln(&sb, "  [yellow]Reconciling:[white]      True")
	}
	if !inst.CreateTime.IsZero() {
		fmt.Fprintf(&sb, "  [lightcyan]Created:[white]          %s (%s)\n", inst.CreateTime.Format("2006-01-02 15:04:05"), humanize.Time(inst.CreateTime))
	}
	if !inst.UpdateTime.IsZero() {
		fmt.Fprintf(&sb, "  [lightcyan]Updated:[white]          %s (%s)\n", inst.UpdateTime.Format("2006-01-02 15:04:05"), humanize.Time(inst.UpdateTime))
	}
	fmt.Fprintln(&sb, "")

	fmt.Fprintln(&sb, "[yellow::b]Security & Identity[white::-]")
	sa := "Default compute service account"
	if inst.ServiceAccount != "" {
		sa = inst.ServiceAccount
	}
	fmt.Fprintf(&sb, "  [lightcyan]Service Account:[white]  %s\n", sa)
	if len(inst.URLs) > 0 {
		fmt.Fprintf(&sb, "  [lightcyan]URLs:[white]             %s\n", strings.Join(inst.URLs, ", "))
	}
	fmt.Fprintln(&sb, "")

	if len(inst.Labels) > 0 {
		fmt.Fprintln(&sb, "[yellow::b]Labels[white::-]")
		for k, v := range inst.Labels {
			fmt.Fprintf(&sb, "  [lightcyan]%s:[white] %s\n", k, v)
		}
		fmt.Fprintln(&sb, "")
	}

	if len(inst.Annotations) > 0 {
		fmt.Fprintln(&sb, "[yellow::b]Annotations[white::-]")
		for k, v := range inst.Annotations {
			fmt.Fprintf(&sb, "  [lightcyan]%s:[white] %s\n", k, v)
		}
	}

	overviewDetail.SetText(sb.String())
}

func updateContainersTab() {
	if dashboardInstance == nil || len(dashboardInstance.Containers) == 0 {
		containersDetail.SetText("No container details available")
		return
	}

	var sb strings.Builder
	for i, c := range dashboardInstance.Containers {
		if i > 0 {
			fmt.Fprint(&sb, "\n----------------------------------------\n\n")
		}
		fmt.Fprintf(&sb, "[yellow::b]Container: %s[white::-]\n", c.Name)
		fmt.Fprintf(&sb, "  [lightcyan]Image:[white]           %s\n", c.Image)

		if c.Resources != nil && len(c.Resources.Limits) > 0 {
			fmt.Fprintln(&sb, "  [lightcyan]Resource Limits:[white]")
			for k, v := range c.Resources.Limits {
				fmt.Fprintf(&sb, "    - %s: %s\n", k, v)
			}
		}

		if len(c.Ports) > 0 {
			fmt.Fprintln(&sb, "  [lightcyan]Ports:[white]")
			for _, p := range c.Ports {
				fmt.Fprintf(&sb, "    - %s: %d\n", p.Name, p.ContainerPort)
			}
		}

		if len(c.Env) > 0 {
			fmt.Fprintln(&sb, "  [lightcyan]Environment Variables:[white]")
			for _, env := range c.Env {
				fmt.Fprintf(&sb, "    - %s: %s\n", env.Name, env.Value)
			}
		}
	}

	if len(dashboardInstance.ContainerStatuses) > 0 {
		fmt.Fprintln(&sb, "\n[yellow::b]Container Runtime Statuses[white::-]")
		for _, cs := range dashboardInstance.ContainerStatuses {
			fmt.Fprintf(&sb, "  [lightcyan]%s:[white] Digest: %s\n", cs.Name, cs.ImageDigest)
		}
	}

	containersDetail.SetText(sb.String())
}

func updateConditionsTab() {
	if dashboardInstance == nil {
		conditionsDetail.SetText("No condition details available")
		return
	}

	var sb strings.Builder
	fmt.Fprintln(&sb, "[yellow::b]Health Conditions[white::-]")
	if len(dashboardInstance.Conditions) == 0 && dashboardInstance.TerminalCondition == nil {
		fmt.Fprintln(&sb, "  No condition records found.")
	} else {
		for _, c := range dashboardInstance.Conditions {
			stateColor := "green"
			switch c.State {
			case "CONDITION_FAILED":
				stateColor = "red"
			case "CONDITION_PENDING":
				stateColor = "yellow"
			}
			fmt.Fprintf(&sb, "  [%s::b]%s[white::-] (%s)\n", stateColor, c.Type, c.State)
			if c.Reason != "" {
				fmt.Fprintf(&sb, "    Reason:  %s\n", c.Reason)
			}
			if c.Message != "" {
				fmt.Fprintf(&sb, "    Message: %s\n", c.Message)
			}
			if !c.LastTransitionTime.IsZero() {
				fmt.Fprintf(&sb, "    Transition: %s (%s)\n", c.LastTransitionTime.Format("2006-01-02 15:04:05"), humanize.Time(c.LastTransitionTime))
			}
			fmt.Fprintln(&sb, "")
		}
	}

	fmt.Fprintln(&sb, "[yellow::b]Networking & Ingress[white::-]")
	auth := "Require authentication"
	if dashboardInstance.InvokerIamDisabled {
		auth = "Allow unauthenticated invocations"
	}
	fmt.Fprintf(&sb, "  [lightcyan]Access:[white]       %s\n", auth)
	if dashboardInstance.IapEnabled {
		fmt.Fprintln(&sb, "  [lightcyan]IAP:[white]          Enabled")
	}
	if len(dashboardInstance.URLs) > 0 {
		fmt.Fprintf(&sb, "  [lightcyan]Endpoints:[white]    %s\n", strings.Join(dashboardInstance.URLs, ", "))
	}

	conditionsDetail.SetText(sb.String())
}

func updateYAMLTab() {
	if dashboardInstance == nil {
		yamlDetail.SetText("No instance details available")
		return
	}

	yamlBytes, err := yaml.Marshal(dashboardInstance)
	if err != nil {
		yamlDetail.SetText(fmt.Sprintf("Error marshalling YAML: %v", err))
		return
	}

	var buf bytes.Buffer
	_ = quick.Highlight(&buf, string(yamlBytes), yamlHighlightLexer, yamlHighlightFormatter, yamlHighlightStyle)
	yamlDetail.SetText(tview.TranslateANSI(buf.String()))
}

// DashboardReload reloads instance details and refreshes all dashboard tabs.
func DashboardReload(app *tview.Application, currentInfo info.Info, inst *model.Instance, onResult func(error)) {
	if inst == nil {
		if onResult != nil {
			onResult(fmt.Errorf("no instance selected"))
		}
		return
	}

	dashboardHeader.SetText(fmt.Sprintf("[yellow::b]%s[white::-] ([lightcyan]%s[white])", inst.ID(), inst.Region))

	go func() {
		project := currentInfo.Project
		if inst.Project != "" {
			project = inst.Project
		}
		region := currentInfo.Region
		if inst.Region != "" {
			region = inst.Region
		}

		loadedInst, err := GetInstanceFunc(project, region, inst.ID())

		app.QueueUpdateDraw(func() {
			if err != nil {
				dashboardInstance = nil
				if onResult != nil {
					onResult(err)
				}
				return
			}

			dashboardInstance = loadedInst
			dashboardHeader.SetText(fmt.Sprintf("[yellow::b]%s[white::-] ([lightcyan]%s[white]) - Health: [green]%s[white]",
				dashboardInstance.ID(), dashboardInstance.Region, dashboardInstance.Status()))

			updateOverviewTab()
			updateContainersTab()
			updateConditionsTab()
			updateYAMLTab()
			DashboardShortcuts()

			if onResult != nil {
				onResult(nil)
			}
		})
	}()
}

// DashboardShortcuts updates the footer shortcuts for the dashboard.
func DashboardShortcuts() {
	if footer.ContextShortcutView == nil {
		return
	}
	footer.ContextShortcutView.Clear()

	s := shortcut.FormatByCategory(shortcut.CategoryInstanceDashboard, nil)
	footer.ContextShortcutView.SetText(s)
}

// DashboardClear clears the dashboard state and components.
func DashboardClear() {
	dashboardInstance = nil
	if overviewDetail != nil {
		overviewDetail.Clear()
	}
	if containersDetail != nil {
		containersDetail.Clear()
	}
	if conditionsDetail != nil {
		conditionsDetail.Clear()
	}
	if yamlDetail != nil {
		yamlDetail.Clear()
	}
	if dashboardHeader != nil {
		dashboardHeader.Clear()
	}
}
