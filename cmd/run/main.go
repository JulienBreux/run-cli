package main

import (
	"os"

	"github.com/JulienBreux/run-cli/internal/run/command"
)

func main() {
	cmd := command.New(os.Stdin, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		_ = command.PrintError(os.Stderr, err)
	}
}
