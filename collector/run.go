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
	"github.com/google/uuid"
	"github.com/ttrnecka/agent_poc/webapi/ws"
)

// functions handling run process

type CommandResult struct {
	Output []byte
	Code   int
	Err    error
}

func run(mes ws.Message, mh *MessageHandler) {
	log.Printf("Running policy for source %s with message: %s", mes.Source, mes.Text)

	envs, parts, probeId := parseEnvAssignments(mes.Text)

	output_folder, err := os.MkdirTemp(*tmpPath, parts[0])
	if err != nil {
		panic(err)
	}
	log.Printf("Created temporary upload folder: %s", output_folder)

	// the rest of the process saves the files to output_folder
	// at the and process the folder and delete the folder
	defer func() {
		processFolder(output_folder, *watchPath, *source, parts[0], probeId)
		log.Println("Deleting temproary upload folder")
		err := os.RemoveAll(output_folder)
		if err != nil {
			log.Println(fmt.Errorf("cannot delete %s: %w", output_folder, err))
		}
		log.Printf("Deleted temporary upload folder: %s", output_folder)
	}()

	parts = append(parts, "--output_folder", output_folder)
	// TODO: obfuscate credentials env variables
	log.Printf("Parsed environment variables: %v", envs)
	log.Printf("Parsed command parts: %v", parts)
	cmd := exec.Command(fmt.Sprintf("./bin/%s", parts[0]), parts[1:]...)
	cmd.Env = append(os.Environ(), envs...)

	ticker := time.NewTicker(2000 * time.Millisecond)
	defer ticker.Stop()

	result := make(chan CommandResult)

	// execute the command in goroutine and pass results to result channel
	go func() {
		cr := CommandResult{}
		log.Printf("Running plugin %s", parts[0])
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
				// If it's another kind of error (e.g., command not found), just set dummy non-0 code
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
			// TODO: return out of this branch needs to handle local message persistence in cas we need to resend it
			text := "Request succeeded"
			mc := ws.MSG_FINISHED_OK
			if cr.Code != 0 {
				text = "Request failed"
				mc = ws.MSG_FINISHED_ERR
			}
			m := ws.NewMessage(mc, *source, mes.Source, text)
			m.Session = mes.Session

			log.Printf("Sending FINISHED message: %v", m)
			err := mh.SendMessage(m)
			if err != nil {
				return
			}

			// processing simulator
			time.Sleep(2000 * time.Millisecond)

			// prepare DATA message
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

			log.Printf("Sending DATA message: %v", m)
			err = mh.SendMessage(m)
			if err != nil {
				return
			}

			return
		case <-ticker.C:
			// TODO this will just send a message, it would be nice if we can stream here the logs
			m := ws.NewMessage(ws.MSG_RUNNING, *source, mes.Source, "Request in progress...")
			m.Session = mes.Session

			log.Printf("Sending RUNNING message: %v", m)
			err = mh.SendMessage(m)
			if err != nil {
				return
			}
		}
	}
}

func parseEnvAssignments(input string) ([]string, []string, string) {
	tokens, err := shlex.Split(input)

	if err != nil {
		log.Fatal(err)
	}

	var envVars []string
	var rest []string
	var probeId string

	for i, token := range tokens {
		if strings.Contains(token, "=") && !strings.HasPrefix(token, "=") {
			// hack to get probe id while the stuff is in POC
			if strings.Contains(token, "PROBE_ID") {
				parts := strings.SplitN(token, "=", 2)
				probeId = parts[1]
			} else {
				envVars = append(envVars, token)
			}
		} else {
			rest = tokens[i:]
			break
		}
	}
	return envVars, rest, probeId
}

func processFolder(src_folder, dest_folder, collector, policy, probeId string) {

	log.Printf("Reading source folder %s", src_folder)
	// Read all entries in the source directory
	entries, err := os.ReadDir(src_folder)
	if err != nil {
		log.Println(fmt.Errorf("failed to read source folder: %w", err))
	}
	uUID := uuid.New().String()
	for _, entry := range entries {
		if entry.IsDir() {
			// Skip subdirectories (you can recurse if needed)
			continue
		}
		srcPath := filepath.Join(src_folder, entry.Name())
		destPath := filepath.Join(dest_folder, entry.Name())

		timestamp, device, endpoint, err := parseFilename(entry.Name())
		log.Printf("Processing file: %s", srcPath)
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
		modifiedContent := []byte(
			"---collector:\t" + collector + "\n" +
				"---probe_id:\t" + probeId + "\n" +
				"---collection_id:\t" + uUID + "\n" +
				"---policy:\t" + policy + "\n" +
				"---timestamp:\t" + timestamp + "\n" +
				"---device:\t" + stripAfterLast(device, ":") + "\n" +
				"---endpoint:\t" + stripAfterLast(endpoint, ".") + "\n" +
				string(content))

		// Write modified content to destination
		err = os.WriteFile(destPath, modifiedContent, 0644)
		if err != nil {
			log.Println(fmt.Errorf("failed to write file %s: %w", destPath, err))
		}
		log.Printf("New file written: %s", destPath)
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
