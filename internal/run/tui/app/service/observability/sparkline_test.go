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

package observability

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderSparkline(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		res := RenderSparkline(nil, 10)
		assert.Equal(t, "[gray]—[-]", res)
	})

	t.Run("ZeroValues", func(t *testing.T) {
		res := RenderSparkline([]float64{0, 0, 0}, 10)
		assert.Equal(t, "   ", res)
	})

	t.Run("ConstantNonZeroValues", func(t *testing.T) {
		res := RenderSparkline([]float64{5, 5, 5}, 10)
		assert.Equal(t, "███", res)
	})

	t.Run("IncreasingValues", func(t *testing.T) {
		res := RenderSparkline([]float64{0, 10, 20, 30, 40, 50, 60, 70}, 10)
		assert.Equal(t, " ▂▃▄▅▆▇█", res)
	})

	t.Run("TruncationWithMaxLen", func(t *testing.T) {
		res := RenderSparkline([]float64{0, 10, 20, 30, 40, 50, 60, 70}, 4)
		assert.Equal(t, "▅▆▇█", res)
	})
}

func TestRenderColoredSparkline(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		res := RenderColoredSparkline(nil, 10, 50, 80)
		assert.Equal(t, "[gray]—[-]", res)
	})

	t.Run("ColoredThresholds", func(t *testing.T) {
		values := []float64{10, 60, 90}
		res := RenderColoredSparkline(values, 10, 50, 80)
		assert.Contains(t, res, "[green]")
		assert.Contains(t, res, "[yellow]")
		assert.Contains(t, res, "[red]")
	})

	t.Run("ZeroRangeWithColor", func(t *testing.T) {
		values := []float64{100, 100}
		res := RenderColoredSparkline(values, 10, 50, 80)
		assert.Contains(t, res, "[red]")
	})
}
