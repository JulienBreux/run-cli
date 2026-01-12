package job

import (
	"fmt"
	"strings"
	"time"

	api_execution "github.com/JulienBreux/run-cli/internal/run/api/job/execution"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_job "github.com/JulienBreux/run-cli/internal/run/model/job"
	model_execution "github.com/JulienBreux/run-cli/internal/run/model/job/execution"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/table"
	"github.com/dustin/go-humanize"
	"github.com/rivo/tview"
)

const (
	DASHBOARD_PAGE_ID = "job-dashboard"
)

var (
	dashboardFlex       *tview.Flex
	dashboardHeader     *tview.TextView
	dashboardPages      *tview.Pages
	dashboardJob        *model_job.Job
	dashboardExecutions []model_execution.Execution

	// Executions tab components
	executionsTable  *table.Table
	executionsDetail *tview.TextView
)

var listExecutionsFunc = api_execution.List

// Dashboard returns the dashboard primitive.
func Dashboard(app *tview.Application) *tview.Flex {
	dashboardHeader = tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	dashboardPages = tview.NewPages()

	// Executions View
	dashboardPages.AddPage("Executions", buildExecutionsTab(app), true, true)

	dashboardFlex = tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(dashboardHeader, 1, 0, false).
		AddItem(dashboardPages, 0, 1, true)

	return dashboardFlex
}

func buildExecutionsTab(app *tview.Application) tview.Primitive {
	executionsTable = table.New(" Executions ")
	executionsTable.SetHeadersWithExpansions(
		[]string{"NAME", "STATUS", "CREATED", "DURATION", "TASKS (S/F)"},
		[]int{2, 1, 1, 1, 1},
	)

	executionsDetail = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetWrap(true)
	executionsDetail.SetBorder(true).SetTitle(" Execution Details ")

	executionsTable.Table.SetSelectionChangedFunc(func(row, column int) {
		updateExecutionDetail(row)
	})

	flex := tview.NewFlex().
		AddItem(executionsTable.Table, 0, 2, true).
		AddItem(executionsDetail, 0, 1, false)

	return flex
}

func updateExecutionDetail(row int) {
	if row < 1 || row > len(dashboardExecutions) {
		executionsDetail.SetText("")
		return
	}
	exec := dashboardExecutions[row-1]

	var sb strings.Builder
	fmt.Fprintf(&sb, "[lightcyan]Name:[white] %s\n", shortName(exec.Name))
	fmt.Fprintf(&sb, "[lightcyan]Created:[white] %s\n", exec.CreateTime.Format("2006-01-02 15:04:05"))
	
	status := "Unknown"
	if exec.TerminalCondition != nil {
		status = exec.TerminalCondition.State
		if exec.TerminalCondition.Message != "" {
			status += fmt.Sprintf(" (%s)", exec.TerminalCondition.Message)
		}
	}
	fmt.Fprintf(&sb, "[lightcyan]Status:[white] %s\n", status)
	
	duration := "-"
	if !exec.CompletionTime.IsZero() {
		duration = exec.CompletionTime.Sub(exec.StartTime).String()
	}
	fmt.Fprintf(&sb, "[lightcyan]Duration:[white] %s\n", duration)
	
	fmt.Fprintln(&sb, "")
	fmt.Fprintln(&sb, "[yellow::b]Tasks[white::-]")
	fmt.Fprintf(&sb, "  [lightcyan]Total:[white] %d\n", exec.TaskCount)
	fmt.Fprintf(&sb, "  [green]Succeeded:[white] %d\n", exec.SucceededCount)
	fmt.Fprintf(&sb, "  [red]Failed:[white] %d\n", exec.FailedCount)
	fmt.Fprintf(&sb, "  [blue]Running:[white] %d\n", exec.RunningCount)
	fmt.Fprintf(&sb, "  [yellow]Retried:[white] %d\n", exec.RetriedCount)
	fmt.Fprintf(&sb, "  [gray]Cancelled:[white] %d\n", exec.CancelledCount)

	executionsDetail.SetText(sb.String())
}

// DashboardReload reloads the dashboard for a specific job.
func DashboardReload(app *tview.Application, currentInfo info.Info, job *model_job.Job, onResult func(error)) {
	dashboardJob = job
	dashboardHeader.SetText(fmt.Sprintf("[lightcyan]Job: [white]%s", shortName(job.Name)))

	go func() {
		var err error
		dashboardExecutions, err = listExecutionsFunc(currentInfo.Project, job.Region, job.Name)

		app.QueueUpdateDraw(func() {
			executionsTable.Table.Clear()
			executionsTable.SetHeadersWithExpansions(
				[]string{"NAME", "STATUS", "CREATED", "DURATION", "TASKS (S/F)"},
				[]int{2, 1, 1, 1, 1},
			)

			if err != nil {
				onResult(err)
				return
			}

			for i, exec := range dashboardExecutions {
				row := i + 1

				status := "-"
				if exec.TerminalCondition != nil {
					status = exec.TerminalCondition.State
				}

				duration := "-"
				if !exec.CompletionTime.IsZero() {
					duration = exec.CompletionTime.Sub(exec.StartTime).Round(time.Second).String()
				}
				
			tasks := fmt.Sprintf("%d (%d/%d)", exec.TaskCount, exec.SucceededCount, exec.FailedCount)

				executionsTable.Table.SetCell(row, 0, tview.NewTableCell(shortName(exec.Name)))
				executionsTable.Table.SetCell(row, 1, tview.NewTableCell(status))
				executionsTable.Table.SetCell(row, 2, tview.NewTableCell(humanize.Time(exec.CreateTime)))
				executionsTable.Table.SetCell(row, 3, tview.NewTableCell(duration))
				executionsTable.Table.SetCell(row, 4, tview.NewTableCell(tasks))
			}

			executionsTable.Table.SetTitle(fmt.Sprintf(" Executions (%d) ", len(dashboardExecutions)))
			if len(dashboardExecutions) > 0 {
				executionsTable.Table.Select(1, 0)
				updateExecutionDetail(1)
			}
			onResult(nil)
		})
	}()
}

// DashboardShortcuts sets the shortcuts for the dashboard.
func DashboardShortcuts() {
	header.ContextShortcutView.Clear()
	shortcuts := `[dodgerblue]<esc> [white]Back`
	header.ContextShortcutView.SetText(shortcuts)
}

func shortName(name string) string {
	parts := strings.Split(name, "/")
	return parts[len(parts)-1]
}