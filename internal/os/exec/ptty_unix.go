//go:build !windows
// +build !windows

package exec

import (
	"io"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"golang.org/x/term"

	"github.com/creack/pty"
	"github.com/gruntwork-io/terragrunt/internal/errors"
	"github.com/gruntwork-io/terragrunt/pkg/log"
)

// runCommandWithPTY will allocate a pseudo-tty to run the subcommand in. This is only necessary when running
// interactive commands, so that terminal features like readline work through the subcommand when stdin, stdout, and
// stderr is being shared.
// NOTE: This is based on the quickstart example from https://github.com/creack/pty
func runCommandWithPTY(logger log.Logger, cmd *exec.Cmd) (err error) {
	// NOTE: in order to ensure we can return errors that occur in cleanup, we use a variable binding for the return
	// value so that it can be updated.
	pseudoTerminal, startErr := pty.Start(cmd)
	defer func() {
		if closeErr := pseudoTerminal.Close(); closeErr != nil {
			logger.Errorf("Error closing pty: %s", closeErr)
			// Only overwrite the previous error if there was no error since this error has lower priority than any
			// errors in the main routine
			if err == nil {
				err = errors.WithStackTrace(closeErr)
			}
		}
	}()

	if startErr != nil {
		return errors.WithStackTrace(startErr)
	}

	// Every time the current terminal size changes, we need to make sure the PTY also updates the size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)

	go func() {
		for range ch {
			if inheritSizeErr := pty.InheritSize(os.Stdin, pseudoTerminal); inheritSizeErr != nil {
				// We don't propagate this error upstream because it does not affect normal operation of the command
				logger.Errorf("error resizing pty: %s", inheritSizeErr)
			}
		}
	}()
	ch <- syscall.SIGWINCH // Make sure the pty matches current size

	// Set stdin in raw mode so that we preserve readline properties
	oldState, setRawErr := term.MakeRaw(int(os.Stdin.Fd()))
	if setRawErr != nil {
		return errors.WithStackTrace(setRawErr)
	}

	defer func() {
		if restoreErr := term.Restore(int(os.Stdin.Fd()), oldState); restoreErr != nil {
			logger.Errorf("Error restoring terminal state: %s", restoreErr)
			// Only overwrite the previous error if there was no error since this error has lower priority than any
			// errors in the main routine
			if err == nil {
				err = errors.WithStackTrace(restoreErr)
			}
		}
	}()

	// Copy stdin to the pty
	go func() {
		_, copyStdinErr := io.Copy(pseudoTerminal, os.Stdin)
		// We don't propagate this error upstream because it does not affect normal operation of the command. A repeat
		// of the same stdin in this case should resolve the issue.
		logger.Errorf("Error forwarding stdin: %s", copyStdinErr)
	}()

	// ... and the pty to stdout.
	_, copyStdoutErr := io.Copy(cmd.Stdout, pseudoTerminal)
	if copyStdoutErr != nil {
		return errors.WithStackTrace(copyStdoutErr)
	}

	return nil
}

func PrepareConsole(logger log.Logger) {
	// No operation function to match windows execution
}