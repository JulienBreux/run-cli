package workerpool

import (
	"errors"
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_workerpool "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/workerpool/scaling"
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

	newWorkers := []model_workerpool.WorkerPool{
		{
			Name:        "projects/p/locations/r/workerPools/pool-1",
			DisplayName: "pool-1",
			Region:      "us-central1",
			UpdateTime:  time.Now(),
			Scaling:     &model_scaling.Scaling{ManualInstanceCount: 2},
		},
	}

	Load(newWorkers)

	assert.Equal(t, newWorkers, workers)
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "pool-1", listTable.Table.GetCell(1, 0).Text)
}

func TestListReload(t *testing.T) {
	app := tview.NewApplication()
	List(app)

	// Mock listWorkerPoolsFunc
	originalListWorkerPoolsFunc := listWorkerPoolsFunc
	defer func() { listWorkerPoolsFunc = originalListWorkerPoolsFunc }()

	expectedWorkers := []model_workerpool.WorkerPool{
		{DisplayName: "pool-reloaded"},
	}
	listWorkerPoolsFunc = func(projectID, region string) ([]model_workerpool.WorkerPool, error) {
		return expectedWorkers, nil
	}

	ListReload(app, info.Info{}, func(err error) {
		assert.NoError(t, err)
		app.Stop()
	})

	go func() {
		time.Sleep(2 * time.Second)
		app.Stop()
	}()

	err := app.Run()
	assert.NoError(t, err)

	assert.Equal(t, expectedWorkers, workers)
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "pool-reloaded", listTable.Table.GetCell(1, 0).Text)
}

func TestListReload_Error(t *testing.T) {
	app := tview.NewApplication()
	List(app)

	originalListWorkerPoolsFunc := listWorkerPoolsFunc
	defer func() { listWorkerPoolsFunc = originalListWorkerPoolsFunc }()

	listWorkerPoolsFunc = func(projectID, region string) ([]model_workerpool.WorkerPool, error) {
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

func TestGetSelectedWorkerPool(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	// Manually populate internal state
	workers = []model_workerpool.WorkerPool{
		{
			Name:        "projects/p/locations/r/workerPools/pool-1",
			DisplayName: "pool-1",
			Region:      "us-central1",
			UpdateTime:  time.Now(),
			Scaling:     &model_scaling.Scaling{ManualInstanceCount: 2},
		},
	}

	// Manually populate table
	row := 1
	listTable.Table.SetCell(row, 0, tview.NewTableCell("pool-1"))
	listTable.Table.SetCell(row, 1, tview.NewTableCell("us-central1"))

	// Select Row 1
	listTable.Table.Select(row, 0)

	name, region := GetSelectedWorkerPool()
	assert.Equal(t, "pool-1", name)
	assert.Equal(t, "us-central1", region)

	w := GetSelectedWorkerPoolFull()
	assert.NotNil(t, w)
	assert.Equal(t, "pool-1", w.DisplayName)

	// Test header selection
	listTable.Table.Select(0, 0)
	name, region = GetSelectedWorkerPool()
	assert.Equal(t, "", name)
	assert.Equal(t, "", region)

	wFull := GetSelectedWorkerPoolFull()
	assert.Nil(t, wFull)
}

func TestGetSelectedWorkerPoolFull_Empty(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	workers = []model_workerpool.WorkerPool{} // Ensure empty

	w := GetSelectedWorkerPoolFull()
	assert.Nil(t, w)
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

	testWorkers := []model_workerpool.WorkerPool{
		{
			DisplayName: "pool-1",
			Region:      "us-central1",
			Scaling:     &model_scaling.Scaling{ManualInstanceCount: 5},
			Labels:      map[string]string{"env": "prod"},
		},
	}

	render(testWorkers)

	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "pool-1", listTable.Table.GetCell(1, 0).Text)
	assert.Contains(t, listTable.Table.GetCell(1, 3).Text, "Manual: 5")
	assert.Contains(t, listTable.Table.GetCell(1, 5).Text, "env: prod")
}