package exec

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

func ExternalCommand(command string) func(input *RunContext, args ...string) *RunContext {
	return func(input *RunContext, args ...string) *RunContext {

		// Create a context if it doesn't exist
		if input.Ctx == nil {
			ctx := context.Background()
			input.Ctx = &ctx
		}

		cmd := exec.CommandContext(*input.Ctx, command, args...)

		// If there is input from a previous command, connect it to this command's stdin
		if input != nil && input.Stdout != nil {
			stdin, err := cmd.StdinPipe()
			if err != nil {
				return &RunContext{
					Err: fmt.Errorf("failed to get stdin pipe: %v", err),
				}
			}

			go func() {
				defer stdin.Close()
				io.Copy(stdin, input.Stdout)
			}()
		}

		// Get stdout and stderr pipes
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return &RunContext{
				Err: fmt.Errorf("failed to get stdout pipe: %v", err),
			}
		}

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return &RunContext{
				Err: fmt.Errorf("failed to get stderr pipe: %v", err),
			}
		}

		// Start the command
		if err := cmd.Start(); err != nil {
			return &RunContext{
				Err: fmt.Errorf("failed to start command: %v", err),
			}
		}

		// Return the RunContext with wait function
		return &RunContext{
			Stdout: stdout,
			Stderr: stderr,
			Ctx:    input.Ctx,
		}
	}
}
