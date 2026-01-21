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
	"io"

	"github.com/JulienBreux/run-cli/pkg/version"
	"github.com/spf13/cobra"
)

var output = ""

// NewCmdVersion returns a command to print version.
func NewCmdVersion(in io.Reader, out, err io.Writer) (cmd *cobra.Command) {
	cmd = &cobra.Command{
		Use:   "version",
		Short: "Print the Run CLI version",
		Long:  "Print the Run CLI version",
		Run:   run(out),
	}

	cmd.Flags().StringVarP(&output, "output", "o", "", "One of '', 'yaml' or 'json'.")

	return
}

// run returns the command.
func run(out io.Writer) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		version.Print(out, output)
	}
}
