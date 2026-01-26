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
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/service/scaling"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
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
			Name:         "service-1",
			Region:       "us-central1",
			URI:          "https://s1.example.com",
			LastModifier: "user@example.com",
			UpdateTime:   time.Now(),
			Scaling: &model_scaling.Scaling{
				ScalingMode:  "AUTOMATIC",
				MinInstances: 1,
				MaxInstances: 5,
			},
		},
		{
			Name:         "service-2",
			Region:       "europe-west1",
			URI:          "https://s2.example.com",
			LastModifier: "user@example.com",
			UpdateTime:   time.Now(),
			Scaling: &model_scaling.Scaling{
				ScalingMode:         "MANUAL",
				ManualInstanceCount: 2,
			},
		},
	}

	Load(testServices)

	// Verify Table Content
	// Header is row 0.
	assert.Equal(t, 3, tbl.Table.GetRowCount()) // 1 header + 2 rows

	// Row 1 (service-1)
	assert.Equal(t, "service-1", tbl.Table.GetCell(1, 1).Text)
	assert.Equal(t, "us-central1", tbl.Table.GetCell(1, 2).Text)
	assert.Contains(t, tbl.Table.GetCell(1, 3).Text, "Auto: min 1, max 5")

	// Row 2 (service-2)
	assert.Equal(t, "service-2", tbl.Table.GetCell(2, 1).Text)
	assert.Contains(t, tbl.Table.GetCell(2, 3).Text, "Manual: 2")
}

func TestGetSelectedService(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	testServices := []model_service.Service{
		{
			Name:   "service-1",
			Region: "us-central1",
			URI:    "https://s1.example.com",
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
	name, _ = GetSelectedService()
	assert.Equal(t, "", name)

	s = GetSelectedServiceFull()
	assert.Nil(t, s)
}

func TestShortcuts(t *testing.T) {
	_ = footer.New()

	assert.NotPanics(t, func() {
		Shortcuts()
	})

	assert.Contains(t, footer.ContextShortcutView.GetText(true), "Refresh")
	assert.Contains(t, footer.ContextShortcutView.GetText(true), "Proxy")
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

func TestHandleShortcuts_ProxyOpenURL(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	testServices := []model_service.Service{
		{
			Name: "s1",
			URI:  "http://public",
			Proxy: &model_service.ProxyStatus{
				Enabled: true,
				Port:    8080,
				URL:     "http://127.0.0.1:8080",
			},
		},
	}
	Load(testServices)
	listTable.Table.Select(1, 0)

	// Test 'o' shortcut
	ev := tcell.NewEventKey(tcell.KeyRune, 'o', tcell.ModNone)

	assert.NotPanics(t, func() {
		HandleShortcuts(ev)
		// We expect this to try opening the PROXY URL.
		// Since we can't intercept the browser.OpenURL call easily in this unit test format,
		// we mainly ensure that the logic doesn't crash and follows the proxy path if possible.
		// Ideally we would mock browser.OpenURL but it is a package level function.
		// For now, this ensures checking Proxy field doesn't crash.
	})
}

func TestShortcuts_Proxy(t *testing.T) {
	_ = footer.New()
	app := tview.NewApplication()
	_ = List(app)

	testServices := []model_service.Service{
		{
			Name: "s1",
			Proxy: &model_service.ProxyStatus{
				Enabled: true,
				Port:    1234,
				URL:     "http://local",
			},
		},
	}
	Load(testServices)
	listTable.Table.Select(1, 0)

	Shortcuts()

	text := footer.ContextShortcutView.GetText(true)
	assert.Contains(t, text, "Proxy (127.0.0.1:1234)")
	assert.Contains(t, text, "Open URL (proxy)")
}

func TestHandleShortcuts_Proxy(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	testServices := []model_service.Service{
		{Name: "s1", URI: "http://test"},
	}
	Load(testServices)
	listTable.Table.Select(1, 0)

	// Mock proxies manually since we can't easily mock net.Listen in integration test without refactor
	// However, we start with no proxy.
	// Press 'p'.
	// It will try to start a proxy. Starting a proxy involves net.Listen on localhost:0 which should work in test env.
	// It involves auth.GetIDToken which might fail or panic if not mocked?
	// The auth.GetIDToken calls idtoken.NewTokenSource which makes network calls.
	// This might fail in test environment without creds.
	// To avoid failure, we might need to skip deep verification or handle the error gracefully in toggleProxy.
	// In toggleProxy: if err != nil { return }
	// So if start fails, nothing happens to model.

	// Let's try to trigger it.
	ev := tcell.NewEventKey(tcell.KeyRune, 'p', tcell.ModNone)

	assert.NotPanics(t, func() {
		HandleShortcuts(ev)
	})

	// If it failed (likely due to auth), the proxy status remains nil.
	// If it succeeded (unlikely without creds), it would be enabled.
	// We just want to ensure it doesn't panic and coverage is hit.
}

func TestRender(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	svcs := []model_service.Service{
		{
			Name:         "s1",
			Region:       "r1",
			URI:          "u1",
			LastModifier: "me",
			UpdateTime:   time.Now(),
			Scaling:      &model_scaling.Scaling{ScalingMode: "AUTOMATIC", MinInstances: 1},
		},
	}

	render(svcs)

	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "s1", listTable.Table.GetCell(1, 1).Text)
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
		assert.Equal(t, "s1", listTable.Table.GetCell(1, 1).Text)
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
			Name:    "s2",
			Scaling: &model_scaling.Scaling{ScalingMode: "MANUAL", ManualInstanceCount: 5},
		},
	}

	render(svcs)

	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Contains(t, listTable.Table.GetCell(1, 3).Text, "Manual: 5")
}
