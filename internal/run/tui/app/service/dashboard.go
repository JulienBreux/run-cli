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

package service

import (
	"fmt"
	"strings"

	api_revision "github.com/JulienBreux/run-cli/internal/run/api/service/revision"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_revision "github.com/JulienBreux/run-cli/internal/run/model/service/revision"
	"github.com/JulienBreux/run-cli/internal/run/tui/app/service/revision"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
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
	revisionsList   *revision.ListComponent
	revisionsDetail *revision.DetailComponent

	// Networking tab components
	networkingDetail *tview.TextView

	// Security tab components
	securityDetail *tview.TextView

	activeTab = 0
	tabs      = []string{"Revisions", "Observability", "Networking", "Security"}
)

var listRevisionsFunc = api_revision.List

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
	dashboardPages.AddPage(tabs[2], buildNetworkingTab(), true, false)
	// Security Tab
	dashboardPages.AddPage(tabs[3], buildSecurityTab(), true, false)

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
	revisionsList = revision.NewListComponent(app)
	revisionsDetail = revision.NewDetailComponent()

	revisionsList.Table.Table.SetSelectionChangedFunc(func(row, column int) {
		updateRevisionDetail(row)
	})

	flex := tview.NewFlex().
		AddItem(revisionsList.Table.View, 0, 2, true).
		AddItem(revisionsDetail.TextView, 0, 1, false)

	return flex
}

func buildNetworkingTab() tview.Primitive {
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

func buildSecurityTab() tview.Primitive {
	securityDetail = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	securityDetail.SetBorder(true).SetTitle(" Security ")
	return securityDetail
}

func updateSecurityTab() {
	if dashboardService == nil || dashboardService.Security == nil {
		securityDetail.SetText("No security information available")
		return
	}

	s := dashboardService.Security

	var sb strings.Builder

	// Authentication
	fmt.Fprintln(&sb, "[yellow::b]Authentication[white::-]")
	auth := "Require authentication"
	if s.InvokerIAMDisabled {
		auth = "Allow unauthenticated invocations"
	}
	fmt.Fprintf(&sb, "  [lightcyan]Access:[white] %s\n", auth)
	fmt.Fprintln(&sb, "")

	// Service Account
	fmt.Fprintln(&sb, "[yellow::b]Service Account[white::-]")
	sa := "Default compute service account"
	if s.ServiceAccount != "" {
		sa = s.ServiceAccount
	}
	fmt.Fprintf(&sb, "  [lightcyan]Identity:[white] %s\n", sa)
	fmt.Fprintln(&sb, "")

	// Encryption
	fmt.Fprintln(&sb, "[yellow::b]Encryption[white::-]")
	enc := "Google-managed key"
	if s.EncryptionKey != "" {
		enc = s.EncryptionKey
	}
	fmt.Fprintf(&sb, "  [lightcyan]Key:[white] %s\n", enc)
	fmt.Fprintln(&sb, "")

	// Binary Authorization
	fmt.Fprintln(&sb, "[yellow::b]Binary Authorization[white::-]")
	binAuth := "Disabled"
	if s.BinaryAuthorization != "" {
		binAuth = fmt.Sprintf("Enabled (Policy: %s)", s.BinaryAuthorization)
		if s.BreakglassJustification != "" {
			binAuth += fmt.Sprintf("\n  [red]Breakglass used:[white] %s", s.BreakglassJustification)
		}
	}
	fmt.Fprintf(&sb, "  [lightcyan]Status:[white] %s\n", binAuth)

	securityDetail.SetText(sb.String())
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
	revisionsDetail.Update(dashboardRevisions[row-1])
}

// GetSelectedRevisions returns the selected revisions from the dashboard.
func GetSelectedRevisions() []model_revision.Revision {
	return revisionsList.GetSelectedRevisions(dashboardRevisions)
}

// DashboardReload reloads the dashboard for a specific service.
func DashboardReload(app *tview.Application, currentInfo info.Info, service *model_service.Service, onResult func(error)) {
	dashboardService = service
	dashboardHeader.SetText(fmt.Sprintf("[lightcyan]Service: [white]%s", service.Name))
	activeTab = 0
	updateTabs()
	updateNetworkingTab()
	updateSecurityTab()

	go func() {
		var err error
		dashboardRevisions, err = listRevisionsFunc(currentInfo.Project, service.Region, service.Name)

		app.QueueUpdateDraw(func() {
			if err != nil {
				onResult(err)
				return
			}

			revisionsList.Update(service, dashboardRevisions)

			if len(dashboardRevisions) > 0 {
				revisionsList.Table.Table.Select(1, 0)
				updateRevisionDetail(1)
			}
			onResult(nil)
		})
	}()
}

// DashboardShortcuts sets the shortcuts for the dashboard.
func DashboardShortcuts() {
	footer.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<esc> [white]Back  [dodgerblue]<tab> [white]Next Tab  [dodgerblue]<shift-tab> [white]Prev Tab  [dodgerblue]t [white]Traffic Split  [dodgerblue]<space> [white]Select`
	footer.ContextShortcutView.SetText(shortcuts)
}
