package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	log.Printf("Running policy for source %s with message: %s", mes.Source, mes.Text)
	// Here you would typically call the function that executes the policy.

	envs, parts := parseEnvAssignments(mes.Text)

	output_folder, err := os.MkdirTemp(*tmpPath, parts[0])
	if err != nil {
		panic(err)
	}

	defer func() {
		processFolder(output_folder, *watchPath, *source, parts[0])
		os.RemoveAll(output_folder)
	}()

	parts = append(parts, "--output_folder", output_folder)
	log.Printf("Parsed environment variables: %v", envs)
	log.Printf("Parsed command parts: %v", parts)
	cmd := exec.Command(fmt.Sprintf("./bin/%s", parts[0]), parts[1:]...)
	cmd.Env = append(os.Environ(), envs...)

	ticker := time.NewTicker(2000 * time.Millisecond)
	defer ticker.Stop()

	type CommandResult struct {
		Output []byte
		Code   int
		Err    error
	}

	result := make(chan CommandResult)

	go func() {
		cr := CommandResult{}
		output, err := cmd.CombinedOutput()

		cr.Output = output
		time.Sleep(3000 * time.Millisecond) // Simulate some processing delay
		log.Printf("Command output: %s", output)

		// Check if there was an error (non-zero exit or command failure)
		if err != nil {
			// If it's an ExitError, we can get the exit code
			cr.Err = err
			if exitErr, ok := err.(*exec.ExitError); ok {
				cr.Code = exitErr.ExitCode()
			} else {
				// If it's another kind of error (e.g., command not found), just print it
				log.Printf("Command execution failed: %v\n", err)
				cr.Code = 255
			}
		}
		log.Printf("Exit Code: %d\n", cr.Code)
		result <- cr
		close(result)
	}()

	for {
		select {
		case cr := <-result:
			text := "Request succeeded"
			mc := ws.MSG_FINISHED_OK
			if cr.Code != 0 {
				text = "Request failed"
				mc = ws.MSG_FINISHED_ERR
			}
			m := ws.NewMessage(mc, *source, mes.Source, text)
			m.Session = mes.Session
			mu.Lock()
			err := wsConn.WriteJSON(m)
			mu.Unlock()
			if err != nil {
				log.Printf("Sending FINISHED message failed: %v\n", err)
				return
			}
			log.Printf("Sending FINISHED message suceeded: %v\n", err)
			time.Sleep(2000 * time.Millisecond)
			var sb strings.Builder
			sb.Write(cr.Output)
			sb.WriteString("\n")

			if cr.Err != nil {
				sb.WriteString(cr.Err.Error())
				sb.WriteString("\n")
			}

			// sb.WriteString(fmt.Sprintf("Exit Code: %d", cr.Code))
			m = ws.NewMessage(ws.MSG_DATA, *source, mes.Source, sb.String())
			m.Session = mes.Session
			mu.Lock()
			err = wsConn.WriteJSON(m)
			mu.Unlock()
			if err != nil {
				log.Printf("Sending DATA message failed: %v\n", err)
				return
			}
			log.Printf("Sending DATA message suceeded: %v\n", err)

			return
		case <-ticker.C:
			m := ws.NewMessage(ws.MSG_RUNNING, *source, mes.Source, "Request in progress...")
			m.Session = mes.Session
			mu.Lock()
			err := wsConn.WriteJSON(m)
			mu.Unlock()
			if err != nil {
				log.Printf("Sending RUNNING message failed: %v\n", err)
				return
			}
			log.Printf("Sending RUNNING message suceeded: %v\n", err)
		}
	}
}

func processFolder(src_folder, dest_folder, collector, probe string) {

	// Read all entries in the source directory
	entries, err := os.ReadDir(src_folder)
	if err != nil {
		log.Println(fmt.Errorf("failed to read source directory: %w", err))
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Skip subdirectories (you can recurse if needed)
			continue
		}

		srcPath := filepath.Join(src_folder, entry.Name())
		destPath := filepath.Join(dest_folder, entry.Name())

		timestamp, device, endpoint, err := parseFilename(entry.Name())
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Read file content
		content, err := os.ReadFile(srcPath)
		if err != nil {
			log.Println(fmt.Errorf("failed to read file %s: %w", srcPath, err))
		}

		// Prepend namePrefix
		modifiedContent := []byte("---collector:\t" + collector + "\n" +
			"---probe:\t" + probe + "\n" +
			"---timestamp:\t" + timestamp + "\n" +
			"---device:\t" + stripAfterLast(device, ":") + "\n" +
			"---endpoint:\t" + stripAfterLast(endpoint, ".") + "\n" +
			string(content))

		// Write modified content to destination
		err = os.WriteFile(destPath, modifiedContent, 0644)
		if err != nil {
			log.Println(fmt.Errorf("failed to write file %s: %w", destPath, err))
		}
	}
}

func parseFilename(input string) (timestamp, device, rest string, err error) {
	parts := strings.SplitN(input, "_", 3)
	if len(parts) != 3 {
		return "", "", "", fmt.Errorf("input string does not match expected format")
	}
	return parts[0], parts[1], parts[2], nil
}

func stripAfterLast(s, subs string) string {
	if idx := strings.LastIndex(s, subs); idx != -1 {
		return s[:idx]
	}
	return s
}
