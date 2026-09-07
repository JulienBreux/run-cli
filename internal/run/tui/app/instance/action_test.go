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

package instance

import (
	"errors"
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestStartAction_Success(t *testing.T) {
	origStart := StartInstanceFunc
	defer func() { StartInstanceFunc = origStart }()

	called := false
	StartInstanceFunc = func(project, region, instanceID string) error {
		called = true
		return nil
	}

	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)

	done := make(chan struct{})
	go func() { _ = app.Run() }()
	defer app.Stop()

	StartAction(app, "proj-1", "us-central1", "inst-1", func(err error) {
		assert.NoError(t, err)
		close(done)
	})

	select {
	case <-done:
		assert.True(t, called)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for StartAction")
	}
}

func TestStartAction_Error(t *testing.T) {
	origStart := StartInstanceFunc
	defer func() { StartInstanceFunc = origStart }()

	StartInstanceFunc = func(project, region, instanceID string) error {
		return errors.New("failed to start")
	}

	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)

	done := make(chan struct{})
	go func() { _ = app.Run() }()
	defer app.Stop()

	StartAction(app, "proj-1", "us-central1", "inst-1", func(err error) {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to start")
		close(done)
	})

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for StartAction error")
	}
}

func TestStopAction_Success(t *testing.T) {
	origStop := StopInstanceFunc
	defer func() { StopInstanceFunc = origStop }()

	called := false
	StopInstanceFunc = func(project, region, instanceID string) error {
		called = true
		return nil
	}

	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)

	done := make(chan struct{})
	go func() { _ = app.Run() }()
	defer app.Stop()

	StopAction(app, "proj-1", "us-central1", "inst-1", func(err error) {
		assert.NoError(t, err)
		close(done)
	})

	select {
	case <-done:
		assert.True(t, called)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for StopAction")
	}
}

func TestStopAction_Error(t *testing.T) {
	origStop := StopInstanceFunc
	defer func() { StopInstanceFunc = origStop }()

	StopInstanceFunc = func(project, region, instanceID string) error {
		return errors.New("failed to stop")
	}

	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)

	done := make(chan struct{})
	go func() { _ = app.Run() }()
	defer app.Stop()

	StopAction(app, "proj-1", "us-central1", "inst-1", func(err error) {
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to stop")
		close(done)
	})

	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for StopAction error")
	}
}

func TestDeleteModal_CreationAndCancel(t *testing.T) {
	app := tview.NewApplication()
	cancelled := false

	modal := DeleteModal(app, "proj-1", "us-central1", "inst-1", func(deleted bool, err error) {
		assert.False(t, deleted)
		assert.NoError(t, err)
		cancelled = true
	})

	assert.NotNil(t, modal)

	// Simulate Esc key capture
	handler := modal.GetInputCapture()
	assert.NotNil(t, handler)
	eventEsc := tcell.NewEventKey(tcell.KeyEscape, 0, tcell.ModNone)
	handler(eventEsc)

	assert.True(t, cancelled)
}

func TestDeleteModal_Confirm(t *testing.T) {
	origDelete := DeleteInstanceFunc
	defer func() { DeleteInstanceFunc = origDelete }()

	called := false
	DeleteInstanceFunc = func(project, region, instanceID string) error {
		called = true
		return nil
	}

	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)

	done := make(chan struct{})
	go func() { _ = app.Run() }()
	defer app.Stop()

	modal := DeleteModal(app, "proj-1", "us-central1", "inst-1", func(deleted bool, err error) {
		assert.True(t, deleted)
		assert.NoError(t, err)
		close(done)
	})

	assert.NotNil(t, modal)

	// Trigger confirmation action directly
	TriggerDeleteConfirmation(app, "proj-1", "us-central1", "inst-1", func(deleted bool, err error) {
		assert.True(t, deleted)
		assert.NoError(t, err)
		close(done)
	})

	select {
	case <-done:
		assert.True(t, called)
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for delete confirmation")
	}
}
