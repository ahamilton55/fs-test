package utils

import (
	"fmt"
	"os"
)

// Print out an error and then quit with the given exit code.
func ErrorAndQuit(msg string, err error, exitCode int) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", msg, err)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}
	os.Exit(exitCode)
}
