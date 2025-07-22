package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/shlex"
	"github.com/gorilla/websocket"
	"github.com/ttrnecka/agent_poc/ws"
)

func parseEnvAssignments(input string) ([]string, []string) {
	tokens, err := shlex.Split(input)

	if err != nil {
		log.Fatal(err)
	}

	var envVars []string
	var rest []string

	for i, token := range tokens {
		if strings.Contains(token, "=") && !strings.HasPrefix(token, "=") {
			envVars = append(envVars, token)
		} else {
			rest = tokens[i:]
			break
		}
	}
	return envVars, rest
}

func run(mes ws.Message, wsConn *websocket.Conn) {
	// This function should implement the logic to run the policy
	// specified in the message. For now, we will just log the action.
	log.Printf("Running policy for collector %s with message: %s", mes.Source, mes.Text)
	// Here you would typically call the function that executes the policy.

	envs, parts := parseEnvAssignments(mes.Text)
	log.Printf("Parsed environment variables: %v", envs)
	log.Printf("Parsed command parts: %v", parts)
	cmd := exec.Command(fmt.Sprintf("./bin/%s", parts[0]), parts[1:]...)
	cmd.Env = append(os.Environ(), envs...)

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	done := make(chan bool)

	go func() {
		defer close(done)
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
	}()

	for {
		select {
		case <-done:
			err := wsConn.WriteJSON(ws.NewMessage(ws.MSG_FINISHED, *source, "hub", "Collector is going offline"))
			if err != nil {
				log.Println("write:", err)
				return
			}
			return
		case t := <-ticker.C:
			fmt.Println("Tick at", t)
		}
	}
}
