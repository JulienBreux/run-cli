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

package domainmapping

import (
	"errors"
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/condition"
	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_domainmapping "github.com/JulienBreux/run-cli/internal/run/model/domainmapping"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/footer"
	"github.com/gdamore/tcell/v2"
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

	newDMs := []model_domainmapping.DomainMapping{
		{
			Name:      "example.com",
			Region:    "us-central1",
			RouteName: "service-1",
		},
	}

	Load(newDMs)

	assert.Equal(t, newDMs, domainMappings)
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "example.com", listTable.Table.GetCell(1, 0).Text)
}

func TestListReload(t *testing.T) {
	app := tview.NewApplication()
	simScreen := tcell.NewSimulationScreen("UTF-8")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init sim screen: %v", err)
	}
	app.SetScreen(simScreen)

	List(app)

	// Mock listDomainMappingsFunc
	originalListDomainMappingsFunc := listDomainMappingsFunc
	defer func() { listDomainMappingsFunc = originalListDomainMappingsFunc }()

	expectedDMs := []model_domainmapping.DomainMapping{
		{Name: "reloaded.example.com"},
	}
	listDomainMappingsFunc = func(projectID, region string) ([]model_domainmapping.DomainMapping, error) {
		return expectedDMs, nil
	}

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

	assert.Equal(t, expectedDMs, domainMappings)
	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "reloaded.example.com", listTable.Table.GetCell(1, 0).Text)
}

func TestListReload_Error(t *testing.T) {
	app := tview.NewApplication()
	simScreen := tcell.NewSimulationScreen("UTF-8")
	if err := simScreen.Init(); err != nil {
		t.Fatalf("failed to init sim screen: %v", err)
	}
	app.SetScreen(simScreen)

	List(app)

	originalListDomainMappingsFunc := listDomainMappingsFunc
	defer func() { listDomainMappingsFunc = originalListDomainMappingsFunc }()

	listDomainMappingsFunc = func(projectID, region string) ([]model_domainmapping.DomainMapping, error) {
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

func TestGetSelectedDomainMappingFull(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	// Manually populate internal state
	domainMappings = []model_domainmapping.DomainMapping{
		{
			Name:   "example.com",
			Region: "us-central1",
		},
	}

	// Manually populate table
	row := 1
	listTable.Table.SetCell(row, 0, tview.NewTableCell("example.com"))

	// Select Row 1
	listTable.Table.Select(row, 0)

	dm := GetSelectedDomainMappingFull()
	assert.NotNil(t, dm)
	assert.Equal(t, "example.com", dm.Name)

	// Test header selection
	listTable.Table.Select(0, 0)
	dmFull := GetSelectedDomainMappingFull()
	assert.Nil(t, dmFull)
}

func TestGetSelectedDomainURL(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	// Manually populate table
	row := 1
	listTable.Table.SetCell(row, 0, tview.NewTableCell("example.com"))

	// Select Row 1
	listTable.Table.Select(row, 0)

	url := GetSelectedDomainURL()
	assert.Equal(t, "https://example.com", url)

	// Test header selection
	listTable.Table.Select(0, 0)
	url = GetSelectedDomainURL()
	assert.Equal(t, "", url)
}

func TestGetSelectedDomainMappingFull_Empty(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)
	domainMappings = []model_domainmapping.DomainMapping{} // Ensure empty

	dm := GetSelectedDomainMappingFull()
	assert.Nil(t, dm)
}

func TestShortcuts(t *testing.T) {
	_ = footer.New()

	// Ensure domainMappings list is not empty
	domainMappings = []model_domainmapping.DomainMapping{{Name: "example.com"}}

	assert.NotPanics(t, func() {
		Shortcuts()
	})

	assert.Contains(t, footer.ContextShortcutView.GetText(true), "Refresh")
	assert.Contains(t, footer.ContextShortcutView.GetText(true), "Open URL")
	assert.Contains(t, footer.ContextShortcutView.GetText(true), "Info")
}

func TestRender(t *testing.T) {
	app := tview.NewApplication()
	_ = List(app)

	testDMs := []model_domainmapping.DomainMapping{
		{
			Name:       "example.com",
			Region:     "us-central1",
			Creator:    "user@example.com",
			Conditions: []*condition.Condition{{Type: "Ready", State: "True"}},
		},
	}

	render(testDMs)

	assert.Equal(t, 2, listTable.Table.GetRowCount())
	assert.Equal(t, "example.com", listTable.Table.GetCell(1, 0).Text)
	assert.Equal(t, "us-central1", listTable.Table.GetCell(1, 2).Text)
	assert.Equal(t, "user@example.com", listTable.Table.GetCell(1, 3).Text)
}

func TestDomainMappingInfoModal(t *testing.T) {
	app := tview.NewApplication()
	dm := &model_domainmapping.DomainMapping{
		Name: "example.com",
		Records: []model_domainmapping.ResourceRecord{
			{Type: "A", Name: "@", RRData: "1.2.3.4"},
		},
		Conditions: []*condition.Condition{
			{Type: "Ready", State: "True"},
		},
	}

	modal := DomainMappingInfoModal(app, dm, func() {})
	assert.NotNil(t, modal)

	// We can assert type is Grid
	_, ok := modal.(*tview.Grid)
	assert.True(t, ok)
}
