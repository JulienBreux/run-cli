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

package job

import (
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/container"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_job "github.com/JulienBreux/run-cli/internal/run/model/job"
	model_execution "github.com/JulienBreux/run-cli/internal/run/model/job/execution"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestDashboard(t *testing.T) {
	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	app.SetScreen(screen)

	d := Dashboard(app)
	assert.NotNil(t, d)
	assert.Equal(t, 2, dashboardFlex.GetItemCount())
}

func TestDashboardReload(t *testing.T) {
	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	if err := screen.Init(); err != nil {
		t.Fatal(err)
	}
	app.SetScreen(screen)

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
		time.Sleep(1 * time.Second)
		app.Stop()
	}()
	_ = app.Run()

	// Assertions
	assert.Equal(t, mockJob, dashboardJob)
	assert.Equal(t, mockExecutions, dashboardExecutions)
	assert.Contains(t, dashboardHeader.GetText(true), "test-job")
	assert.Equal(t, 2, executionsTable.Table.GetRowCount()) // Header + 1 row
	assert.Equal(t, "exec-1", executionsTable.Table.GetCell(1, 0).Text)
}

func TestDashboardShortcuts(t *testing.T) {
	_ = footer.New()
	
	assert.NotPanics(t, func() {
		DashboardShortcuts()
	})

	assert.Contains(t, footer.ContextShortcutView.GetText(true), "Back")
}