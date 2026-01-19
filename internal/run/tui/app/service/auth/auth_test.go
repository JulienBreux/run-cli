package auth

import (
	"testing"

	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestModal(t *testing.T) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	service := &model_service.Service{
		Name:    "test-service",
		Project: "test-project",
		Region:  "us-central1",
	}

	modal := Modal(app, service, pages, func() {})
	assert.NotNil(t, modal)
}