package service

import (
	"fmt"
	"strings"

	api_revision "github.com/JulienBreux/run-cli/internal/run/api/service/revision"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_revision "github.com/JulienBreux/run-cli/internal/run/model/service/revision"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/header"
	"github.com/JulienBreux/run-cli/internal/run/ui/component/table"
	"github.com/dustin/go-humanize"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	DASHBOARD_PAGE_ID = "service-dashboard"
)

var (
	dashboardFlex      *tview.Flex
	dashboardHeader    *tview.TextView
	dashboardTabs      *tview.TextView
	dashboardPages     *tview.Pages
	dashboardService   *model_service.Service
	dashboardRevisions []model_revision.Revision

	// Revisions tab components
	revisionsTable  *table.Table
	revisionsDetail *tview.TextView

	// Networking tab components
	networkingDetail *tview.TextView

	activeTab = 0
	tabs      = []string{"Revisions", "Observability", "Networking", "Security"}
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

	// Revisions Tab
	dashboardPages.AddPage(tabs[0], buildRevisionsTab(app), true, true)
	// Observability Tab
	dashboardPages.AddPage(tabs[1], tview.NewBox().SetTitle(" Observability (Placeholder) ").SetBorder(true), true, false)
	// Networking Tab
	dashboardPages.AddPage(tabs[2], buildNetworkingTab(app), true, false)
	// Security Tab
	dashboardPages.AddPage(tabs[3], tview.NewBox().SetTitle(" Security (Placeholder) ").SetBorder(true), true, false)

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
		return event
	})

	return dashboardFlex
}

func buildRevisionsTab(app *tview.Application) tview.Primitive {
	revisionsTable = table.New(" Revisions ")
	revisionsTable.SetHeadersWithExpansions(
		[]string{"NAME", "TRAFFIC", "DEPLOYED", "REVISION TAGS"},
		[]int{2, 1, 1, 2},
	)

	revisionsDetail = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	revisionsDetail.SetBorder(true).SetTitle(" Revision Details ")

	revisionsTable.Table.SetSelectionChangedFunc(func(row, column int) {
		updateRevisionDetail(row)
	})

	flex := tview.NewFlex().
		AddItem(revisionsTable.Table, 0, 2, true).
		AddItem(revisionsDetail, 0, 1, false)

	return flex
}

func buildNetworkingTab(app *tview.Application) tview.Primitive {
	networkingDetail = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	networkingDetail.SetBorder(true).SetTitle(" Networking ")
	return networkingDetail
}

func updateNetworkingTab() {
	if dashboardService == nil || dashboardService.Networking == nil {
		networkingDetail.SetText("No networking information available")
		return
	}

	n := dashboardService.Networking

	var sb strings.Builder
	fmt.Fprintln(&sb, "[yellow::b]Ingress[white::-]")
	ingress := n.Ingress
	switch ingress {
	case "INGRESS_TRAFFIC_ALL":
		ingress = "Allow all traffic"
	case "INGRESS_TRAFFIC_INTERNAL_ONLY":
		ingress = "Allow internal traffic only"
	case "INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER":
		ingress = "Allow internal traffic and traffic from Cloud Load Balancing"
	}
	fmt.Fprintf(&sb, "  [lightcyan]Traffic settings:[white] %s\n", ingress)
	fmt.Fprintln(&sb, "")

	fmt.Fprintln(&sb, "[yellow::b]Endpoints[white::-]")
	enabled := "Enabled"
	if n.DefaultUriDisabled {
		enabled = "Disabled"
	}
	fmt.Fprintf(&sb, "  [lightcyan]Default URL:[white] %s (%s)\n", dashboardService.URI, enabled)
	if n.IapEnabled {
		fmt.Fprintln(&sb, "  [lightcyan]IAP:[white] Enabled")
	}
	fmt.Fprintln(&sb, "")

	fmt.Fprintln(&sb, "[yellow::b]VPC[white::-]")
	if n.VpcAccess != nil {
		fmt.Fprintf(&sb, "  [lightcyan]Connector:[white] %s\n", n.VpcAccess.Connector)
		egress := n.VpcAccess.Egress
		switch egress {
		case "ALL_TRAFFIC":
			egress = "Route all traffic through the VPC connector"
		case "PRIVATE_RANGES_ONLY":
			egress = "Route only traffic to private IP addresses through the VPC connector"
		}
		fmt.Fprintf(&sb, "  [lightcyan]Egress:[white] %s\n", egress)
	} else {
		fmt.Fprintln(&sb, "  No VPC connector configured")
	}

	networkingDetail.SetText(sb.String())
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

func updateRevisionDetail(row int) {
	if row < 1 || row > len(dashboardRevisions) {
		revisionsDetail.SetText("")
		return
	}
	rev := dashboardRevisions[row-1]

	var sb strings.Builder
	fmt.Fprintf(&sb, "[lightcyan]Name:[white] %s\n", rev.Name)
	fmt.Fprintf(&sb, "[lightcyan]Created:[white] %s\n", rev.CreateTime.Format("2006-01-02 15:04:05"))
	fmt.Fprintln(&sb, "")
	fmt.Fprintln(&sb, "[yellow::b]General[white::-]")

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
	fmt.Fprintln(&sb, "[yellow::b]Containers[white::-]")
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

	revisionsDetail.SetText(sb.String())
}

// DashboardReload reloads the dashboard for a specific service.
func DashboardReload(app *tview.Application, currentInfo info.Info, service *model_service.Service, onResult func(error)) {
	dashboardService = service
	dashboardHeader.SetText(fmt.Sprintf("[lightcyan]Service: [white]%s", service.Name))
	activeTab = 0
	updateTabs()
	updateNetworkingTab()

	go func() {
		var err error
		dashboardRevisions, err = api_revision.List(currentInfo.Project, service.Region, service.Name)

		app.QueueUpdateDraw(func() {
			revisionsTable.Table.Clear()
			revisionsTable.SetHeadersWithExpansions(
				[]string{"NAME", "TRAFFIC", "DEPLOYED", "REVISION TAGS"},
				[]int{2, 1, 1, 2},
			)

			if err != nil {
				onResult(err)
				return
			}

			for i, rev := range dashboardRevisions {
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

				revisionsTable.Table.SetCell(row, 0, tview.NewTableCell(rev.Name))
				revisionsTable.Table.SetCell(row, 1, tview.NewTableCell(traffic))
				revisionsTable.Table.SetCell(row, 2, tview.NewTableCell(humanize.Time(rev.CreateTime)))
				revisionsTable.Table.SetCell(row, 3, tview.NewTableCell(tags))
			}

			revisionsTable.Table.SetTitle(fmt.Sprintf(" Revisions (%d) ", len(dashboardRevisions)))
			if len(dashboardRevisions) > 0 {
				revisionsTable.Table.Select(1, 0)
				updateRevisionDetail(1)
			}
			onResult(nil)
		})
	}()
}

// DashboardShortcuts sets the shortcuts for the dashboard.
func DashboardShortcuts() {
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<esc> [white]Back
[dodgerblue]<tab> [white]Next Tab
[dodgerblue]<shift-tab> [white]Prev Tab`
	header.ContextShortcutView.SetText(shortcuts)
}
