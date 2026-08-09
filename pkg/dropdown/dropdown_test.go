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

package dropdown_test

import (
	"testing"

	"github.com/JulienBreux/run-cli/pkg/dropdown"
	"github.com/gdamore/tcell/v2"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	d := dropdown.New()
	assert.NotNil(t, d)
}

func TestDraw(t *testing.T) {
	d := dropdown.New()
	d.SetRect(0, 0, 20, 1) // Set a rect so it has size
	screen := tcell.NewSimulationScreen("")
	err := screen.Init()
	assert.NoError(t, err)

	d.Draw(screen)
	// No panic implies success for now
}
