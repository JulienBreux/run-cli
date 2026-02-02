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

	model_revision "github.com/JulienBreux/run-cli/internal/run/model/service/revision"
	model_service "github.com/JulienBreux/run-cli/internal/run/model/service"
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
	rev := model_revision.Revision{Name: "rev1", Author: "user@example.com"}

	comp.Update(rev)
	assert.Contains(t, comp.GetText(true), "rev1")
	assert.Contains(t, comp.GetText(true), "user@example.com")
}
