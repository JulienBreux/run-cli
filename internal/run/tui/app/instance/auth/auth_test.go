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

package auth

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	model "github.com/JulienBreux/run-cli/internal/run/model/instance"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestModal(t *testing.T) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	inst := &model.Instance{
		Name:               "projects/test-project/locations/us-central1/instances/test-inst",
		Project:            "test-project",
		Region:             "us-central1",
		InvokerIamDisabled: false,
	}

	completed := false
	refreshVal := false
	modal := Modal(app, inst, pages, func(refresh bool) {
		completed = true
		refreshVal = refresh
	})
	assert.NotNil(t, modal)

	grid, ok := modal.(*tview.Grid)
	assert.True(t, ok)

	// Test Escape key
	handler := grid.GetInputCapture()
	assert.NotNil(t, handler)
	ev := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	res := handler(ev)
	assert.Nil(t, res)
	assert.True(t, completed)
	assert.False(t, refreshVal)

	instPublic := &model.Instance{
		Name:               "projects/test-project/locations/us-central1/instances/test-inst-public",
		Project:            "test-project",
		Region:             "us-central1",
		InvokerIamDisabled: true,
	}
	modalPublic := Modal(app, instPublic, pages, func(refresh bool) {})
	assert.NotNil(t, modalPublic)
}

func TestModal_NonEscapeKey(t *testing.T) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	inst := &model.Instance{
		Name:               "projects/test-project/locations/us-central1/instances/test-inst",
		Project:            "test-project",
		Region:             "us-central1",
		InvokerIamDisabled: false,
	}

	canceled := false
	modal := Modal(app, inst, pages, func(refresh bool) {
		if !refresh {
			canceled = true
		}
	})

	grid := modal.(*tview.Grid)
	handler := grid.GetInputCapture()
	assert.NotNil(t, handler)
	otherKey := tcell.NewEventKey(tcell.KeyRune, 'a', tcell.ModNone)
	ret := handler(otherKey)
	assert.Equal(t, otherKey, ret)
	assert.False(t, canceled)
}

func TestTriggerSave_Success_WithoutApp(t *testing.T) {
	origFunc := UpdateAuthenticationFunc
	defer func() { UpdateAuthenticationFunc = origFunc }()

	called := false
	UpdateAuthenticationFunc = func(ctx context.Context, project, region, name string, allowUnauthenticated bool) (*model.Instance, error) {
		called = true
		assert.Equal(t, "test-project", project)
		assert.Equal(t, "us-central1", region)
		assert.Equal(t, "test-inst", name)
		assert.True(t, allowUnauthenticated)
		return &model.Instance{}, nil
	}

	inst := &model.Instance{
		Name:    "test-inst",
		Project: "test-project",
		Region:  "us-central1",
	}

	var wg sync.WaitGroup
	wg.Add(1)

	var completedVal bool
	var statusVal string

	TriggerSave(nil, inst, true, func(refresh bool) {
		completedVal = refresh
		wg.Done()
	}, func(status string) {
		statusVal = status
	})

	wg.Wait()
	assert.True(t, called)
	assert.True(t, completedVal)
	assert.Equal(t, "", statusVal)
}

func TestTriggerSave_Error_WithoutApp(t *testing.T) {
	origFunc := UpdateAuthenticationFunc
	defer func() { UpdateAuthenticationFunc = origFunc }()

	UpdateAuthenticationFunc = func(ctx context.Context, project, region, name string, allowUnauthenticated bool) (*model.Instance, error) {
		return nil, errors.New("IAM update failed")
	}

	inst := &model.Instance{
		Name:    "test-inst",
		Project: "test-project",
		Region:  "us-central1",
	}

	var wg sync.WaitGroup
	wg.Add(1)

	var statusVal string
	TriggerSave(nil, inst, false, func(refresh bool) {
		t.Fatal("onCompletion should not be called on error")
	}, func(status string) {
		statusVal = status
		wg.Done()
	})

	wg.Wait()
	assert.Contains(t, statusVal, "IAM update failed")
}

func TestTriggerSave_Success_WithApp(t *testing.T) {
	origFunc := UpdateAuthenticationFunc
	defer func() { UpdateAuthenticationFunc = origFunc }()

	UpdateAuthenticationFunc = func(ctx context.Context, project, region, name string, allowUnauthenticated bool) (*model.Instance, error) {
		return &model.Instance{}, nil
	}

	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)

	done := make(chan struct{})
	go func() { _ = app.Run() }()
	defer app.Stop()

	inst := &model.Instance{
		Name:    "test-inst",
		Project: "test-project",
		Region:  "us-central1",
	}

	TriggerSave(app, inst, true, func(refresh bool) {
		assert.True(t, refresh)
		close(done)
	}, nil)

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for TriggerSave")
	}
}

func TestTriggerSave_Error_WithApp(t *testing.T) {
	origFunc := UpdateAuthenticationFunc
	defer func() { UpdateAuthenticationFunc = origFunc }()

	UpdateAuthenticationFunc = func(ctx context.Context, project, region, name string, allowUnauthenticated bool) (*model.Instance, error) {
		return nil, errors.New("network error")
	}

	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)

	done := make(chan struct{})
	go func() { _ = app.Run() }()
	defer app.Stop()

	inst := &model.Instance{
		Name:    "test-inst",
		Project: "test-project",
		Region:  "us-central1",
	}

	TriggerSave(app, inst, true, func(refresh bool) {
		t.Fatal("onCompletion should not be called on error")
	}, func(status string) {
		assert.Contains(t, status, "network error")
		close(done)
	})

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for TriggerSave error")
	}
}
