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
		fmt.Fprintf(w, format, "Version:", Version)
		fmt.Fprintf(w, format, "Commit:", Commit)
		fmt.Fprintf(w, format, "Build date:", RawDate)
	}
	var v = version{
		Version: Version,
		Commit:  Commit,
		Date:    RawDate,
	}
	format.Print(w, format.StringToFormat(f), v, c)
}
