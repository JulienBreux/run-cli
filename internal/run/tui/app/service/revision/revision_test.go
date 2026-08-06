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

package revision

import (
	"testing"
	"time"

	model_container "github.com/JulienBreux/run-cli/internal/run/model/common/container"
	model_resources "github.com/JulienBreux/run-cli/internal/run/model/common/resources"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
	model_revision "github.com/JulienBreux/run-cli/internal/run/model/service/revision"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestNewListComponent(t *testing.T) {
	app := tview.NewApplication()
	comp := NewListComponent(app)
	assert.NotNil(t, comp)
	assert.IsType(t, &ListComponent{}, comp)
}

func TestNewDetailComponent(t *testing.T) {
	comp := NewDetailComponent()
	assert.NotNil(t, comp)
	assert.IsType(t, &DetailComponent{}, comp)
}

func TestListComponent_Update(t *testing.T) {
	app := tview.NewApplication()
	comp := NewListComponent(app)
	svc := &model_service.Service{Name: "s1"}
	revs := []model_revision.Revision{{Name: "rev1"}}

	comp.Update(svc, revs)
	assert.Equal(t, 2, comp.Table.Table.GetRowCount()) // Header + 1 row
}

func TestDetailComponent_Update(t *testing.T) {
	comp := NewDetailComponent()

	rev := model_revision.Revision{
		Name:                          "rev1",
		Author:                        "user@example.com",
		CreateTime:                    time.Now(),
		CpuIdle:                       true,
		StartupCpuBoost:               true,
		MaxInstanceRequestConcurrency: 80,
		Timeout:                       300 * time.Second,
		ExecutionEnvironment:          "EXECUTION_ENVIRONMENT_GEN2",
		Containers: []*model_container.Container{
			{
				Name:  "c1",
				Image: "image1",
				Ports: []*model_container.Port{{ContainerPort: 8080}},
				Resources: &model_resources.Resources{
					Limits: map[string]string{
						"memory": "512Mi",
						"cpu":    "1",
					},
				},
			},
		},
	}

	comp.Update(rev)
	text := comp.GetText(true)
	assert.Contains(t, text, "rev1")
	assert.Contains(t, text, "user@example.com")
	assert.Contains(t, text, "CPU is only allocated")
	assert.Contains(t, text, "Enabled")
	assert.Contains(t, text, "80")
	assert.Contains(t, text, "5m0s")
	assert.Contains(t, text, "Second Generation")
	assert.Contains(t, text, "c1")
	assert.Contains(t, text, "image1")
	assert.Contains(t, text, "8080")
	assert.Contains(t, text, "512Mi Memory, 1 CPU")
}

func TestListComponent_Clear(t *testing.T) {
	app := tview.NewApplication()
	comp := NewListComponent(app)
	svc := &model_service.Service{Name: "s1"}
	revs := []model_revision.Revision{{Name: "rev1"}}

	comp.Update(svc, revs)
	assert.Equal(t, 2, comp.Table.Table.GetRowCount())

	comp.Clear()
	assert.Equal(t, 1, comp.Table.Table.GetRowCount()) // Header only
}

func TestDetailComponent_Clear(t *testing.T) {
	comp := NewDetailComponent()
	rev := model_revision.Revision{Name: "rev1"}
	comp.Update(rev)
	assert.Contains(t, comp.GetText(true), "rev1")

	comp.Clear()
	assert.Equal(t, "", comp.GetText(true))
}
