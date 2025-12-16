package command

import (
	"io"

	"github.com/JulienBreux/run-cli/internal/run/ui/app"
	"github.com/spf13/cobra"
)

// New returns the root of CLI.
func New(in io.Reader, out, err io.Writer) (cmd *cobra.Command) {
	// Create CLI
	cmd = &cobra.Command{
		Use:   "run",
		Short: "run is a CLI to play with Google Cloud Run interactively.",
		// Version: ver.Version,
		RunE: run,
	}

	// cmd.AddCommand(version.NewCmdVersion(in, out, err))

	return
}

// run runs the root command.
func run(cmd *cobra.Command, args []string) error {
	return app.Run()
}
