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
