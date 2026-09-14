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
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"regexp"
	"runtime/debug"
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

func TestResolveFromBuildInfo(t *testing.T) {
	origVersion := version.Version
	origCommit := version.Commit
	origDate := version.RawDate
	defer func() {
		version.Version = origVersion
		version.Commit = origCommit
		version.RawDate = origDate
	}()

	t.Run("NilBuildInfo", func(t *testing.T) {
		version.Version = "dev"
		version.Commit = "n/a"
		version.RawDate = "n/a"

		version.ResolveFromBuildInfo(nil)

		assert.Equal(t, "dev", version.Version)
		assert.Equal(t, "n/a", version.Commit)
		assert.Equal(t, "n/a", version.RawDate)
	})

	t.Run("PopulateFromBuildInfo", func(t *testing.T) {
		version.Version = "dev"
		version.Commit = "n/a"
		version.RawDate = "n/a"

		info := &debug.BuildInfo{
			Main: debug.Module{
				Version: "v1.2.3",
			},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "0123456789abcdef"},
				{Key: "vcs.time", Value: "2026-09-14T10:00:00Z"},
				{Key: "vcs.modified", Value: "false"},
			},
		}

		version.ResolveFromBuildInfo(info)

		assert.Equal(t, "v1.2.3", version.Version)
		assert.Equal(t, "0123456789abcdef", version.Commit)
		assert.Equal(t, "2026-09-14T10:00:00Z", version.RawDate)
	})

	t.Run("DirtyVCS", func(t *testing.T) {
		version.Version = "dev"
		version.Commit = "n/a"
		version.RawDate = "n/a"

		info := &debug.BuildInfo{
			Main: debug.Module{
				Version: "v1.2.3",
			},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "0123456789abcdef"},
				{Key: "vcs.modified", Value: "true"},
			},
		}

		version.ResolveFromBuildInfo(info)

		assert.Equal(t, "v1.2.3", version.Version)
		assert.Equal(t, "0123456789abcdef-dirty", version.Commit)
	})

	t.Run("PrecedenceOfLdflags", func(t *testing.T) {
		version.Version = "v2.0.0"
		version.Commit = "fedcba9876543210"
		version.RawDate = "2026-01-01T00:00:00Z"

		info := &debug.BuildInfo{
			Main: debug.Module{
				Version: "v1.0.0",
			},
			Settings: []debug.BuildSetting{
				{Key: "vcs.revision", Value: "newrevision"},
				{Key: "vcs.time", Value: "2026-09-14T12:00:00Z"},
			},
		}

		version.ResolveFromBuildInfo(info)

		assert.Equal(t, "v2.0.0", version.Version)
		assert.Equal(t, "fedcba9876543210", version.Commit)
		assert.Equal(t, "2026-01-01T00:00:00Z", version.RawDate)
	})

	t.Run("DevelOrEmptyVersionIgnored", func(t *testing.T) {
		version.Version = "dev"
		version.Commit = "n/a"
		version.RawDate = "n/a"

		info := &debug.BuildInfo{
			Main: debug.Module{
				Version: "(devel)",
			},
		}

		version.ResolveFromBuildInfo(info)

		assert.Equal(t, "dev", version.Version)
	})
}

func TestCompareSemver(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"v0.1.0", "v0.2.0", -1},
		{"v0.2.0", "v0.1.0", 1},
		{"v1.0.0", "v1.0.0", 0},
		{"1.0.0", "v1.0.0", 0},
		{"v1.2.3", "v1.2.4", -1},
		{"v2.0.0", "v1.9.9", 1},
		{"v0.28.1", "v0.29.0", -1},
		{"v0.28.1-0.20260914-commit", "v0.29.0", -1},
		{"v0.29.0", "v0.28.1-0.20260914-commit", 1},
		{"dev", "v0.29.0", 0},
		{"invalid", "v0.29.0", 0},
	}

	for _, tt := range tests {
		actual := version.CompareSemver(tt.v1, tt.v2)
		assert.Equal(t, tt.expected, actual, "CompareSemver(%s, %s)", tt.v1, tt.v2)
	}
}

func TestCheckUpdate(t *testing.T) {
	t.Run("UpdateAvailable", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/repos/JulienBreux/run-cli/releases/latest", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tag_name": "v0.30.0", "html_url": "https://github.com/JulienBreux/run-cli/releases/tag/v0.30.0"}`))
		}))
		defer server.Close()

		origURL := version.GitHubReleaseURL
		version.GitHubReleaseURL = server.URL + "/repos/JulienBreux/run-cli/releases/latest"
		defer func() { version.GitHubReleaseURL = origURL }()

		hasUpdate, latest, err := version.CheckUpdate(context.Background(), "v0.29.0")
		assert.NoError(t, err)
		assert.True(t, hasUpdate)
		assert.Equal(t, "v0.30.0", latest)
	})

	t.Run("AlreadyUpToDate", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tag_name": "v0.29.0"}`))
		}))
		defer server.Close()

		origURL := version.GitHubReleaseURL
		version.GitHubReleaseURL = server.URL + "/repos/JulienBreux/run-cli/releases/latest"
		defer func() { version.GitHubReleaseURL = origURL }()

		hasUpdate, latest, err := version.CheckUpdate(context.Background(), "v0.29.0")
		assert.NoError(t, err)
		assert.False(t, hasUpdate)
		assert.Equal(t, "v0.29.0", latest)
	})

	t.Run("DevVersionNoUpdate", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"tag_name": "v0.30.0"}`))
		}))
		defer server.Close()

		origURL := version.GitHubReleaseURL
		version.GitHubReleaseURL = server.URL + "/repos/JulienBreux/run-cli/releases/latest"
		defer func() { version.GitHubReleaseURL = origURL }()

		hasUpdate, latest, err := version.CheckUpdate(context.Background(), "dev")
		assert.NoError(t, err)
		assert.False(t, hasUpdate)
		assert.Equal(t, "", latest)
	})

	t.Run("ServerErrorGracefullyHandled", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		origURL := version.GitHubReleaseURL
		version.GitHubReleaseURL = server.URL + "/repos/JulienBreux/run-cli/releases/latest"
		defer func() { version.GitHubReleaseURL = origURL }()

		hasUpdate, latest, err := version.CheckUpdate(context.Background(), "v0.29.0")
		assert.Error(t, err)
		assert.False(t, hasUpdate)
		assert.Equal(t, "", latest)
	})
}

