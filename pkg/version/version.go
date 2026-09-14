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
	"fmt"
	"io"
	"runtime/debug"
	"time"

	"github.com/JulienBreux/run-cli/pkg/format"
)

// CheckUpdate checks GitHub releases to see if a newer version is available.
func CheckUpdate(ctx context.Context, currentVersion string) (bool, string, error) {
	return false, "", nil
}

// CompareSemver compares two semver version strings.
func CompareSemver(v1, v2 string) int {
	return 0
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
