package service

import (
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/service/scaling"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestListAndLoad(t *testing.T) {
	app := tview.NewApplication()
	
	// 1. Initialize List
	tbl := List(app)
	assert.NotNil(t, tbl)
	assert.Equal(t, LIST_PAGE_TITLE, tbl.Title)

	// 2. Load Data
	testServices := []model_service.Service{
		{
			Name: "service-1",
			Region: "us-central1",
			URI: "https://s1.example.com",
			LastModifier: "user@example.com",
			UpdateTime: time.Now(),
			Scaling: &model_scaling.Scaling{
				ScalingMode: "AUTOMATIC",
				MinInstances: 1,
				MaxInstances: 5,
			},
		},
		{
			Name: "service-2",
			Region: "europe-west1",
			URI: "https://s2.example.com",
			LastModifier: "user@example.com",
			UpdateTime: time.Now(),
			Scaling: &model_scaling.Scaling{
				ScalingMode: "MANUAL",
				ManualInstanceCount: 2,
			},
		},
	}
	
	Load(testServices)
	
	// Verify Table Content
	// Header is row 0.
	assert.Equal(t, 3, tbl.Table.GetRowCount()) // 1 header + 2 rows
	
	// Row 1 (service-1)
	assert.Equal(t, "service-1", tbl.Table.GetCell(1, 0).Text)
	assert.Equal(t, "us-central1", tbl.Table.GetCell(1, 1).Text)
	assert.Contains(t, tbl.Table.GetCell(1, 2).Text, "Auto: min 1, max 5")
	
	// Row 2 (service-2)
	assert.Equal(t, "service-2", tbl.Table.GetCell(2, 0).Text)
	assert.Contains(t, tbl.Table.GetCell(2, 2).Text, "Manual: 2")
}

func TestGetSelectedService(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	
	testServices := []model_service.Service{
		{
			Name: "service-1",
			Region: "us-central1",
			URI: "https://s1.example.com",
		},
	}
	Load(testServices)
	
	// Select Row 1
	listTable.Table.Select(1, 0)
	
	name, region := GetSelectedService()
	assert.Equal(t, "service-1", name)
	assert.Equal(t, "us-central1", region)
	
	url := GetSelectedServiceURL()
	assert.Equal(t, "https://s1.example.com", url)
	
	s := GetSelectedServiceFull()
	assert.NotNil(t, s)
	assert.Equal(t, "service-1", s.Name)
	
	// Test Header selection (Row 0)
	listTable.Table.Select(0, 0)
	name, region = GetSelectedService()
	assert.Equal(t, "", name)
	
	s = GetSelectedServiceFull()
	assert.Nil(t, s)
}

func TestShortcuts(t *testing.T) {
	_ = header.New(info.Info{})
	
	assert.NotPanics(t, func() {
		Shortcuts()
	})
	
	assert.Contains(t, header.ContextShortcutView.GetText(true), "Refresh")
}

func TestHandleShortcuts(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	testServices := []model_service.Service{
		{Name: "s1", URI: "http://test"},
	}
	Load(testServices)
	listTable.Table.Select(1, 0)
	
	// Test 'o' shortcut
	ev := tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone)
	
	// browser.OpenURL might fail or do nothing in test env, but we check if event is consumed (returns nil)
	// Actually HandleShortcuts calls browser.OpenURL which might panic or log if no browser.
	// Since we can't easily mock browser package here, we just check if it runs without panic.
	
	assert.NotPanics(t, func() {
		ret := HandleShortcuts(ev)
		// If URL is present, it returns nil (consumed) or event (if failed/empty).
		// Our dummy URL is valid string but browser open might fail.
		_ = ret
	})
	
	// Test unknown key
	ev2 := tcell.NewEventKey(tcell.KeyRune, 'z', tcell.ModNone)
	ret := HandleShortcuts(ev2)
	assert.Equal(t, ev2, ret)
}

func TestRender(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	
	svcs := []model_service.Service{
		{
			Name: "s1",
			Region: "r1",
			URI: "u1",
			LastModifier: "me",
			UpdateTime: time.Now(),
			Scaling: &model_scaling.Scaling{ScalingMode: "AUTOMATIC", MinInstances: 1},
		},
	}
	
	render(svcs)
	
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "s1", listTable.Table.GetCell(1, 0).Text)
}

func TestFetch(t *testing.T) {
	origList := listServicesFunc
	defer func() { listServicesFunc = origList }()
	
	listServicesFunc = func(projectID, region string) ([]model_service.Service, error) {
		return []model_service.Service{{Name: "s1"}}, nil
	}
	
	svcs, err := Fetch("p", "r")
	assert.NoError(t, err)
	assert.Len(t, svcs, 1)
	assert.Equal(t, "s1", svcs[0].Name)
}

func TestListReload(t *testing.T) {
	// Mock
	origList := listServicesFunc
	defer func() { listServicesFunc = origList }()
	
	listServicesFunc = func(projectID, region string) ([]model_service.Service, error) {
		return []model_service.Service{{Name: "s1"}}, nil
	}
	
	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)
	
	// Init table
	_ = List(app)
	
	go func() {
		_ = app.Run()
	}()
	defer app.Stop()
	
	done := make(chan struct{})
	ListReload(app, info.Info{}, func(err error) {
		assert.NoError(t, err)
		close(done)
	})
	
	select {
	case <-done:
		// Verify Render was called (Table should have data)
		// Header + 1 Item = 2 Rows
		assert.Equal(t, 2, listTable.Table.GetRowCount())
		assert.Equal(t, "s1", listTable.Table.GetCell(1, 0).Text)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for ListReload")
	}
}

func TestListReload_Error(t *testing.T) {
	// Mock Error
	origList := listServicesFunc
	defer func() { listServicesFunc = origList }()
	
	listServicesFunc = func(projectID, region string) ([]model_service.Service, error) {
		return nil, assert.AnError
	}
	
	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)
	
	// Init table
	_ = List(app)
	
	go func() {
		_ = app.Run()
	}()
	defer app.Stop()
	
	done := make(chan struct{})
	ListReload(app, info.Info{}, func(err error) {
		assert.Error(t, err)
		close(done)
	})
	
	select {
	case <-done:
		// Passed
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for ListReload Error")
	}
}

func TestRender_ScalingManual(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	
	svcs := []model_service.Service{
		{
			Name: "s2",
			Scaling: &model_scaling.Scaling{ScalingMode: "MANUAL", ManualInstanceCount: 5},
		},
	}
	
	render(svcs)
	
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Contains(t, listTable.Table.GetCell(1, 2).Text, "Manual: 5")
}
