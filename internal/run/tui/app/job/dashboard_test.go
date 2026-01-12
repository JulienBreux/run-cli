package job

import (
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/container"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_job "github.com/JulienBreux/run-cli/internal/run/model/job"
	model_execution "github.com/JulienBreux/run-cli/internal/run/model/job/execution"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestDashboard(t *testing.T) {
	app := tview.NewApplication()
	d := Dashboard(app)
	assert.NotNil(t, d)
	assert.Equal(t, 2, dashboardFlex.GetItemCount())
}

func TestDashboardReload(t *testing.T) {
	app := tview.NewApplication()
	Dashboard(app) // Initialize

	// Mock Data
	mockJob := &model_job.Job{
		Name:   "test-job",
		Region: "us-central1",
		Template: &model_job.ExecutionTemplate{
			TaskCount:   5,
			Parallelism: 2,
			Template: &model_job.TaskTemplate{
				Containers: []*container.Container{
					{
						Name:  "c1",
						Image: "img:latest",
					},
				},
				MaxRetries: 3,
			},
		},
	}

	mockExecutions := []model_execution.Execution{
		{
			Name:              "exec-1",
			CreateTime:        time.Now(),
			StartTime:         time.Now(),
			CompletionTime:    time.Now().Add(10 * time.Second),
			TaskCount:         5,
			SucceededCount:    5,
			TerminalCondition: &condition.Condition{State: "Succeeded"},
		},
	}

	// Mock API
	originalListExecutionsFunc := listExecutionsFunc
	defer func() { listExecutionsFunc = originalListExecutionsFunc }()
	listExecutionsFunc = func(project, region, jobName string) ([]model_execution.Execution, error) {
		return mockExecutions, nil
	}

	// Call Reload
	DashboardReload(app, info.Info{Project: "p"}, mockJob, func(err error) {
		assert.NoError(t, err)
		app.Stop()
	})

	// Wait for async update
	go func() {
		time.Sleep(100 * time.Millisecond)
		app.Stop()
	}()
	app.Run()

	// Assertions
	assert.Equal(t, mockJob, dashboardJob)
	assert.Equal(t, mockExecutions, dashboardExecutions)
	assert.Contains(t, dashboardHeader.GetText(true), "test-job")
	assert.Equal(t, 2, executionsTable.Table.GetRowCount()) // Header + 1 row
	assert.Equal(t, "exec-1", executionsTable.Table.GetCell(1, 0).Text)
}

func TestDashboardShortcuts(t *testing.T) {
	_ = header.New(info.Info{})
	
	assert.NotPanics(t, func() {
		DashboardShortcuts()
	})
}