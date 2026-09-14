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
	"strings"
)

var sparkBlocks = []rune{' ', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

// RenderSparkline generates a Unicode sparkline string for the given slice of values.
// If maxLen > 0 and len(values) > maxLen, it takes the most recent maxLen values.
func RenderSparkline(values []float64, maxLen int) string {
	if len(values) == 0 {
		return "[gray]—[-]"
	}

	if maxLen > 0 && len(values) > maxLen {
		values = values[len(values)-maxLen:]
	}

	minVal := values[0]
	maxVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	// Always anchor min to 0 if all positive, so bar heights reflect real scale
	if minVal >= 0 {
		minVal = 0
	}

	valRange := maxVal - minVal
	var sb strings.Builder

	for _, v := range values {
		if valRange == 0 {
			sb.WriteRune(sparkBlocks[0])
			continue
		}

		idx := int(((v - minVal) / valRange) * float64(len(sparkBlocks)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(sparkBlocks) {
			idx = len(sparkBlocks) - 1
		}
		sb.WriteRune(sparkBlocks[idx])
	}

	return sb.String()
}

// RenderColoredSparkline generates a sparkline with tview color tags according to warning/critical thresholds.
func RenderColoredSparkline(values []float64, maxLen int, warnThreshold, critThreshold float64) string {
	if len(values) == 0 {
		return "[gray]—[-]"
	}

	if maxLen > 0 && len(values) > maxLen {
		values = values[len(values)-maxLen:]
	}

	minVal := values[0]
	maxVal := values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	if minVal >= 0 {
		minVal = 0
	}

	valRange := maxVal - minVal
	var sb strings.Builder
	currentColor := ""

	for _, v := range values {
		targetColor := "green"
		if critThreshold > 0 && v >= critThreshold {
			targetColor = "red"
		} else if warnThreshold > 0 && v >= warnThreshold {
			targetColor = "yellow"
		}

		if targetColor != currentColor {
			sb.WriteString("[" + targetColor + "]")
			currentColor = targetColor
		}

		if valRange == 0 {
			sb.WriteRune(sparkBlocks[0])
			continue
		}

		idx := int(((v - minVal) / valRange) * float64(len(sparkBlocks)-1))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(sparkBlocks) {
			idx = len(sparkBlocks) - 1
		}
		sb.WriteRune(sparkBlocks[idx])
	}

	if currentColor != "" {
		sb.WriteString("[-]")
	}

	return sb.String()
}
