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

package scale

import (
	"testing"

	model_workerpool "github.com/JulienBreux/run-cli/internal/run/model/workerpool"
	model_scaling "github.com/JulienBreux/run-cli/internal/run/model/workerpool/scaling"
	"github.com/rivo/tview"
	"github.com/stretchr/testify/assert"
)

func TestModal(t *testing.T) {
	app := tview.NewApplication()
	pages := tview.NewPages()
	
	pool := &model_workerpool.WorkerPool{
		DisplayName: "pool-1",
		Scaling: &model_scaling.Scaling{
			ManualInstanceCount: 2,
		},
	}
	
	onCompletion := func() {}
	
	modal := Modal(app, pool, pages, onCompletion)
	
	assert.NotNil(t, modal)
	_, ok := modal.(*tview.Grid)
	assert.True(t, ok, "Expected Modal to return a Grid")
}

func TestValidateScaleParams(t *testing.T) {
	tests := []struct {
		input     string
		want      int64
		wantErr   bool
	}{
		{"5", 5, false},
		{"0", 0, false},
		{"-1", 0, true},
		{"abc", 0, true},
		{"", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := validateScaleParams(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
