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
	"fmt"
	"io"
	"time"

	"github.com/JulienBreux/run-cli/pkg/format"
)

var (
	// Version is the semver release name of this build
	Version = "dev"
	// Commit is the commit hash this build was created from
	Commit = "n/a"
	// RawDate is the time when this build was created in raw string
	RawDate = "n/a"
)

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
