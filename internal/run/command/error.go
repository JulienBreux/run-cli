package command

import (
	"fmt"
	"io"
)

// PrintError prints error properly to the writer.
func PrintError(w io.Writer, err error) error {
	_, printErr := fmt.Fprintf(w, "%s\n", err)
	return printErr
}
