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

package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/JulienBreux/run-cli/pkg/format"
)

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

// CheckUpdate checks GitHub releases to see if a newer version is available.
func CheckUpdate(ctx context.Context, currentVersion string) (bool, string, error) {
	if currentVersion == "dev" || currentVersion == "" {
		return false, "", nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, GitHubReleaseURL, nil)
	if err != nil {
		return false, "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "run-cli/"+currentVersion)

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return false, "", fmt.Errorf("unexpected status code from release API: %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, "", err
	}

	if release.TagName == "" {
		return false, "", nil
	}

	if CompareSemver(currentVersion, release.TagName) < 0 {
		return true, release.TagName, nil
	}

	return false, release.TagName, nil
}

// CompareSemver compares two semver version strings.
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
func CompareSemver(v1, v2 string) int {
	maj1, min1, pat1, ok1 := parseSemver(v1)
	maj2, min2, pat2, ok2 := parseSemver(v2)
	if !ok1 || !ok2 {
		return 0
	}
	if maj1 != maj2 {
		if maj1 < maj2 {
			return -1
		}
		return 1
	}
	if min1 != min2 {
		if min1 < min2 {
			return -1
		}
		return 1
	}
	if pat1 != pat2 {
		if pat1 < pat2 {
			return -1
		}
		return 1
	}

	hasPre1 := strings.Contains(strings.TrimPrefix(v1, "v"), "-")
	hasPre2 := strings.Contains(strings.TrimPrefix(v2, "v"), "-")
	if hasPre1 && !hasPre2 {
		return -1
	}
	if !hasPre1 && hasPre2 {
		return 1
	}
	return 0
}

func parseSemver(v string) (major, minor, patch int, ok bool) {
	v = strings.TrimPrefix(v, "v")
	parts := strings.Split(v, "-")
	base := parts[0]
	subParts := strings.Split(base, ".")
	if len(subParts) == 0 {
		return 0, 0, 0, false
	}
	var err error
	major, err = strconv.Atoi(subParts[0])
	if err != nil {
		return 0, 0, 0, false
	}
	if len(subParts) > 1 {
		minor, err = strconv.Atoi(subParts[1])
		if err != nil {
			return 0, 0, 0, false
		}
	}
	if len(subParts) > 2 {
		patch, err = strconv.Atoi(subParts[2])
		if err != nil {
			return 0, 0, 0, false
		}
	}
	return major, minor, patch, true
}

var (
	// GitHubReleaseURL is the endpoint used to check for the latest release.
	GitHubReleaseURL = "https://api.github.com/repos/JulienBreux/run-cli/releases/latest"

	// Version is the semver release name of this build
	Version = "dev"
	// Commit is the commit hash this build was created from
	Commit = "n/a"
	// RawDate is the time when this build was created in raw string
	RawDate = "n/a"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		ResolveFromBuildInfo(info)
	}
}

// ResolveFromBuildInfo extracts version, commit, and date from the provided build info.
func ResolveFromBuildInfo(info *debug.BuildInfo) {
	if info == nil {
		return
	}

	if Version == "dev" && info.Main.Version != "" && info.Main.Version != "(devel)" {
		Version = info.Main.Version
	}

	var revision, vcsTime string
	var modified bool
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			revision = setting.Value
		case "vcs.time":
			vcsTime = setting.Value
		case "vcs.modified":
			modified = setting.Value == "true"
		}
	}

	if Commit == "n/a" && revision != "" {
		if modified {
			Commit = revision + "-dirty"
		} else {
			Commit = revision
		}
	}

	if RawDate == "n/a" && vcsTime != "" {
		RawDate = vcsTime
	}
}

// version represents a version
type version struct {
	Version string `yaml:"version" json:"version"`
	Commit  string `yaml:"commit" json:"commit"`
	Date    string `yaml:"date" json:"date"`
}

// Date returns the version's date
func Date() (time.Time, error) {
	t, err := time.Parse(time.RFC3339, RawDate)
	if err != nil {
		return t, &DateParseError{Date: RawDate, Err: err}
	}

	return t, nil
}

// Print prints the version
func Print(w io.Writer, f string) {
	var c format.Callback = func(w io.Writer) {
		const format = "%-15s %s\n"
		_, _ = fmt.Fprintf(w, format, "Version:", Version)
		_, _ = fmt.Fprintf(w, format, "Commit:", Commit)
		_, _ = fmt.Fprintf(w, format, "Build date:", RawDate)
	}
	var v = version{
		Version: Version,
		Commit:  Commit,
		Date:    RawDate,
	}
	format.Print(w, format.StringToFormat(f), v, c)
}
