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

package version_test

import (
	"bytes"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/JulienBreux/run-cli/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestVersionDateFailed(t *testing.T) {
	version.RawDate = "n/a"

	expectedErr := "unable to parse date: n/a"
	_, err := version.Date()

	assert.Error(t, err, expectedErr)
	assert.Equal(t, expectedErr, err.Error())

	var expectedErrType = err.(*version.DateParseError)
	assert.True(t, errors.As(err, &expectedErrType))
	assert.Equal(
		t,
		expectedErrType.Unwrap().Error(),
		"parsing time \"n/a\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"n/a\" as \"2006\"",
	)
}

func TestVersionDateSuccess(t *testing.T) {
	version.RawDate = "1987-01-16T09:00:00Z"

	d, err := version.Date()

	assert.NoError(t, err)
	assert.Equal(t, d.Year(), 1987)
	assert.Equal(t, d.Month(), time.January)
	assert.Equal(t, d.Day(), 16)
}

func TestPrintVersionJSON(t *testing.T) {
	var r *regexp.Regexp
	w := &bytes.Buffer{}

	version.Print(w, "json")
	r = regexp.MustCompile(`{"version":"dev","commit":"n/a","date":"[0-9T:Z-]+"}`)
	assert.Regexp(t, r, w.String())
}

func TestPrintVersionYAML(t *testing.T) {
	var r *regexp.Regexp
	w := &bytes.Buffer{}

	version.Print(w, "yaml")
	r = regexp.MustCompile(`version: dev\ncommit: n/a\ndate: "[0-9T:Z-]+"\n`)
	assert.Regexp(t, r, w.String())
}

func TestPrintVersionText(t *testing.T) {
	var r *regexp.Regexp
	w := &bytes.Buffer{}

	version.Print(w, "")
	r = regexp.MustCompile(`Version:\s+dev\nCommit:\s+n/a\nBuild date:\s+[0-9T:Z-]+\n`)
	assert.Regexp(t, r, w.String())
}
