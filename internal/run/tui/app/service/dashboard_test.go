package service

import (
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/internal/run/model/common/info"
	model_container "github.com/JulienBreux/run-cli/internal/run/model/common/container"
	model_resources "github.com/JulienBreux/run-cli/internal/run/model/common/resources"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_networking "github.com/JulienBreux/run-cli/internal/run/model/service/networking"
	model_revision "github.com/JulienBreux/run-cli/internal/run/model/service/revision"
	model_security "github.com/JulienBreux/run-cli/internal/run/model/service/security"
	"github.com/JulienBreux/run-cli/internal/run/tui/component/header"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestDashboard(t *testing.T) {
	app := tview.NewApplication()
	d := Dashboard(app)
	assert.NotNil(t, d)
	// Header, Tabs, Pages
	assert.Equal(t, 3, d.GetItemCount())
}

func TestDashboardShortcuts(t *testing.T) {
	// Initialize header global view
	_ = header.New(info.Info{})
	
	assert.NotPanics(t, func() {
		DashboardShortcuts()
	})
	
	assert.Contains(t, header.ContextShortcutView.GetText(true), "Back")
}

func TestUpdateTabs(t *testing.T) {
	app := tview.NewApplication()
	_ = Dashboard(app)
	
	// Setup dummy service
	dashboardService = &model_service.Service{
		Name: "s1",
		Networking: &model_networking.Networking{
			Ingress: "INGRESS_TRAFFIC_ALL",
		},
		Security: &model_security.Security{
			InvokerIAMDisabled: true,
		},
	}
	
	assert.NotPanics(t, func() {
		updateNetworkingTab()
		assert.Contains(t, networkingDetail.GetText(true), "Allow all traffic")
		
		updateSecurityTab()
		assert.Contains(t, securityDetail.GetText(true), "Allow unauthenticated invocations")
	})
}

func TestDashboardReload(t *testing.T) {
	// Setup Mocks
	origList := listRevisionsFunc
	defer func() { listRevisionsFunc = origList }()
	
	called := false
	listRevisionsFunc = func(project, region, service string) ([]model_revision.Revision, error) {
		called = true
		return []model_revision.Revision{
			{Name: "rev1", CreateTime: time.Now()},
		}, nil
	}
	
	// Init
	app := tview.NewApplication()
	// Use SimulationScreen to allow QueueUpdateDraw to work if we ran app
	// However, QueueUpdateDraw blocks if app not running or just queues it.
	// We can manually trigger the callback passed to QueueUpdateDraw if we mock app?
	// Tview Application struct is hard to mock internal logic.
	
	// But we can just run the function and check if it started the goroutine/called mock.
	// Since it's async, we use channels or wait.
	
	svc := &model_service.Service{Name: "s1", Region: "r1"}
	done := make(chan struct{})
	
	DashboardReload(app, info.Info{Project: "p"}, svc, func(err error) {
		assert.NoError(t, err)
		close(done)
	})
	
	// Since QueueUpdateDraw might not execute without running App, 
	// we might need to run app or rely on the fact that we mocked listRevisionsFunc.
	// Wait, listRevisionsFunc is called synchronously? No, inside goroutine.
	
	// To test this properly without race conditions or hanging, we should use the SimulationScreen pattern 
	// similar to Spinner test if we really want to execute the callback.
	
	// Let's try to just wait for 'called' to be true? No, that's racy.
	
	// Ideally we use SimulationScreen and run app.
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)
	
	go func() {
		// Run app loop to process updates
		_ = app.Run()
	}()
	defer app.Stop()
	
	select {
	case <-done:
		assert.True(t, called)
		// The callback runs AFTER update. So table should be populated.
		// Header row + 1 item = 2 rows
		assert.Equal(t, 2, revisionsTable.Table.GetRowCount())
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for DashboardReload")
	}
}

func TestDashboardInputCapture(t *testing.T) {
	app := tview.NewApplication()
	d := Dashboard(app)
	
	handler := d.GetInputCapture()
	assert.NotNil(t, handler)
	
	// Initial tab is 0
	activeTab = 0
	
	// Test Tab (Right)
	eventTab := tcell.NewEventKey(tcell.KeyTab, 0, tcell.ModNone)
	handler(eventTab)
	assert.Equal(t, 1, activeTab)
	
	// Test Backtab (Left)
	eventBack := tcell.NewEventKey(tcell.KeyBacktab, 0, tcell.ModNone)
	handler(eventBack)
	assert.Equal(t, 0, activeTab)
}

func TestDashboardReload_Error(t *testing.T) {
	origList := listRevisionsFunc
	defer func() { listRevisionsFunc = origList }()
	
	listRevisionsFunc = func(project, region, service string) ([]model_revision.Revision, error) {
		return nil, assert.AnError
	}
	
	app := tview.NewApplication()
	screen := tcell.NewSimulationScreen("UTF-8")
	_ = screen.Init()
	app.SetScreen(screen)
	
	go func() { _ = app.Run() }()
	defer app.Stop()
	
	svc := &model_service.Service{Name: "s1"}
	done := make(chan struct{})
	
	DashboardReload(app, info.Info{}, svc, func(err error) {
		assert.Error(t, err)
		close(done)
	})
	
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout")
	}
}

func TestUpdateRevisionDetail(t *testing.T) {
	// Initialize global variables
	app := tview.NewApplication()
	_ = Dashboard(app)
	
	// Setup data
	dashboardRevisions = []model_revision.Revision{
		{
			Name: "rev1",
			ExecutionEnvironment: "EXECUTION_ENVIRONMENT_GEN1",
			Containers: []*model_container.Container{
				{
					Name: "c1",
					Resources: &model_resources.Resources{
						Limits: map[string]string{
							"memory": "512Mi",
							"cpu": "1",
						},
					},
				},
			},
		},
		{
			Name: "rev2",
			ExecutionEnvironment: "EXECUTION_ENVIRONMENT_GEN2",
			Accelerator: "nvidia-tesla-t4",
			Containers: []*model_container.Container{
				{
					Resources: &model_resources.Resources{
						Limits: map[string]string{
							"memory": "1Gi",
							"cpu": "2",
							"nvidia.com/gpu": "1",
						},
					},
				},
			},
		},
	}
	
	assert.NotPanics(t, func() {
		// Row 1 (rev1)
		updateRevisionDetail(1)
		text := revisionsDetail.GetText(true)
		assert.Contains(t, text, "First Generation")
		assert.Contains(t, text, "512Mi Memory, 1 CPU")
		
		// Row 2 (rev2)
		updateRevisionDetail(2)
		text2 := revisionsDetail.GetText(true)
		assert.Contains(t, text2, "Second Generation")
		assert.Contains(t, text2, "1 GPU (nvidia-tesla-t4)")
		
		// Invalid Row
		updateRevisionDetail(0)
		assert.Equal(t, "", revisionsDetail.GetText(true))
	})
}
