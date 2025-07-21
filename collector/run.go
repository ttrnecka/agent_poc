package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/ttrnecka/agent_poc/ws"
)

func run(mes ws.Message) {
	// This function should implement the logic to run the policy
	// specified in the message. For now, we will just log the action.
	log.Printf("Running policy for collector %s with message: %s", mes.Source, mes.Text)
	// Here you would typically call the function that executes the policy.

	parts := strings.Fields(mes.Text)
	cmd := exec.Command(fmt.Sprintf("./bin/%s", parts[0]), parts[1:]...)
	output, err := cmd.CombinedOutput()

	log.Printf("Command output: %s", output)
	// Get the exit code
	// Default exit code
	exitCode := 0

	// Check if there was an error (non-zero exit or command failure)
	if err != nil {
		// If it's an ExitError, we can get the exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			// If it's another kind of error (e.g., command not found), just print it
			log.Printf("Command execution failed: %v\n", err)
			return
		}
	}
	log.Printf("Exit Code: %d\n", exitCode)
}
