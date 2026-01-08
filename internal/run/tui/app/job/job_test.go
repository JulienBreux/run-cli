package job

import (
	"errors"
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_job "github.com/JulienBreux/run-cli/internal/run/model/job"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	app := tview.NewApplication()
	tbl := List(app)
	assert.NotNil(t, tbl)
	assert.Equal(t, LIST_PAGE_TITLE, tbl.Title)
}

func TestLoad(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	newJobs := []model_job.Job{
		{
			Name:              "projects/p/locations/r/jobs/job-1",
			Region:            "us-central1",
			TerminalCondition: &condition.Condition{State: "Succeeded"},
		},
	}

	Load(newJobs)

	assert.Equal(t, newJobs, jobs)
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "job-1", listTable.Table.GetCell(1, 0).Text)
}

func TestListReload(t *testing.T) {
	app := tview.NewApplication()
	List(app)

	// Mock listJobsFunc
	originalListJobsFunc := listJobsFunc
	defer func() { listJobsFunc = originalListJobsFunc }()

	expectedJobs := []model_job.Job{
		{Name: "projects/p/locations/r/jobs/job-reloaded"},
	}
	listJobsFunc = func(projectID, region string) ([]model_job.Job, error) {
		return expectedJobs, nil
	}

	// ListReload calls api_job.List which we mocked.
	// It then calls app.QueueUpdateDraw.
	// We need to run the app to process the queue.
	// The callback passed to ListReload calls app.Stop() to exit the Run loop.

	ListReload(app, info.Info{}, func(err error) {
		assert.NoError(t, err)
		app.Stop()
	})

	// Timeout safety
	go func() {
		time.Sleep(2 * time.Second)
		app.Stop()
	}()

	err := app.Run()
	assert.NoError(t, err)

	assert.Equal(t, expectedJobs, jobs)
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "job-reloaded", listTable.Table.GetCell(1, 0).Text)
}

func TestListReload_Error(t *testing.T) {
	app := tview.NewApplication()
	List(app)

	originalListJobsFunc := listJobsFunc
	defer func() { listJobsFunc = originalListJobsFunc }()

	listJobsFunc = func(projectID, region string) ([]model_job.Job, error) {
		return nil, errors.New("fetch error")
	}

	ListReload(app, info.Info{}, func(err error) {
		assert.Error(t, err)
		app.Stop()
	})

	go func() {
		time.Sleep(2 * time.Second)
		app.Stop()
	}()

	err := app.Run()
	assert.NoError(t, err)
}

func TestGetSelectedJob(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	// Manually populate internal state
	jobs = []model_job.Job{
		{
			Name:                   "projects/p/locations/us-central1/jobs/job-1",
			Region:                 "us-central1",
			Creator:                "user@example.com",
			TerminalCondition:      &condition.Condition{State: "Succeeded"},
			LatestCreatedExecution: &model_job.ExecutionReference{CreateTime: time.Now()},
		},
	}

	// Manually populate table (simplified, just what GetSelectedJob needs)
	// 0: Name, 3: Region
	row := 1
	listTable.Table.SetCell(row, 0, tview.NewTableCell("job-1"))
	listTable.Table.SetCell(row, 3, tview.NewTableCell("us-central1"))

	// Select Row 1
	listTable.Table.Select(row, 0)

	name, region := GetSelectedJob()
	assert.Equal(t, "job-1", name)
	assert.Equal(t, "us-central1", region)

	j := GetSelectedJobFull()
	assert.NotNil(t, j)
	assert.Equal(t, "projects/p/locations/us-central1/jobs/job-1", j.Name)

	// Test header selection
	listTable.Table.Select(0, 0)
	name, region = GetSelectedJob()
	assert.Equal(t, "", name)

	jFull := GetSelectedJobFull()
	assert.Nil(t, jFull)
}

func TestGetSelectedJobFull_Empty(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	jobs = []model_job.Job{} // Ensure empty

	j := GetSelectedJobFull()
	assert.Nil(t, j)
}

func TestShortcuts(t *testing.T) {
	_ = header.New(info.Info{})

	assert.NotPanics(t, func() {
		Shortcuts()
	})

	assert.Contains(t, header.ContextShortcutView.GetText(true), "Refresh")
}

func TestRender(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	testJobs := []model_job.Job{
		{
			Name:              "projects/p/locations/r/jobs/job-1",
			Region:            "us-central1",
			TerminalCondition: &condition.Condition{State: "Succeeded"},
		},
	}

	render(testJobs)

	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "job-1", listTable.Table.GetCell(1, 0).Text)
	assert.Equal(t, "Succeeded", listTable.Table.GetCell(1, 1).Text)
}