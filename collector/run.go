package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/shlex"
	"github.com/google/uuid"
	"github.com/ttrnecka/agent_poc/common"
	"github.com/ttrnecka/agent_poc/webapi/ws"
)

// functions handling run process

type CommandResult struct {
	Output []byte
	Code   int
	Err    error
}

func run(mes ws.Message, mh *MessageHandler) {
	logger.Info().Str("source", mes.Source).Str("text", mes.Text).Msg("Running policy")

	result := make(chan CommandResult)

	closeResult := func(cr CommandResult) {
		result <- cr
		close(result)
	}
	go runNotifyLoop(mes, mh, result)

	envs, parts, probeId, err := parseEnvAssignments(mes.Text)
	if err != nil {
		closeResult(CommandResult{Code: 255, Err: err})
		return
	}

	output_folder, err := os.MkdirTemp(*tmpPath, parts[0])
	if err != nil {
		closeResult(CommandResult{Code: 255, Err: err})
		return
	}
	logger.Info().Str("folder", output_folder).Msg("Created temporary upload folder")

	// the rest of the process saves the files to output_folder
	// at the and process the folder and delete the folder
	defer func() {
		processFolder(output_folder, *watchPath, *source, parts[0], probeId)
		logger.Info().Msg("Deleting temporary upload folder")
		err := os.RemoveAll(output_folder)
		if err != nil {
			logger.Error().Err(err).Str("folder", output_folder).Msg("Cannot delete folder")
		}
		logger.Info().Str("folder", output_folder).Msg("Deleted temporary upload folder")
	}()

	parts = append(parts, "--output_folder", output_folder)
	// TODO: obfuscate credentials env variables
	logger.Debug().Str("envs", fmt.Sprintf("%+v", envs)).Msg("Parsed environment variables")
	logger.Debug().Str("parts", fmt.Sprintf("%+v", parts)).Msg("Parsed command parts")

	cmd := exec.Command(fmt.Sprintf("./bin/%s", parts[0]), parts[1:]...)
	cmd.Env = append(os.Environ(), envs...)

	cr := CommandResult{}
	logger.Info().Str("policy", parts[0]).Msg("Running policy")
	output, err := cmd.CombinedOutput()

	cr.Output = output
	time.Sleep(3000 * time.Millisecond) // Simulate some processing delay
	logger.Debug().Str("output", string(output)).Msg("Command output")

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
	logger.Info().Int("exit code", cr.Code).Msg("")
	closeResult(cr)
}

func runNotifyLoop(mes ws.Message, mh *MessageHandler, result chan CommandResult) {
	ticker := time.NewTicker(2000 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case cr := <-result:
			// TODO: return out of this branch needs to handle local message persistence in case we need to resend it
			text := "Request succeeded"
			mc := ws.MSG_FINISHED_OK
			if cr.Code != 0 {
				text = "Request failed"
				mc = ws.MSG_FINISHED_ERR
			}
			m := ws.NewMessage(mc, *source, mes.Source, text)
			m.Session = mes.Session

			logger.Info().Str("raw", fmt.Sprintf("%+v", m)).Msg("Sending FINISHED message")
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

			logger.Info().Str("raw", fmt.Sprintf("%+v", m)).Msg("Sending DATA message")
			err = mh.SendMessage(m)
			if err != nil {
				return
			}

			return
		case <-ticker.C:
			// TODO this will just send a message, it would be nice if we can stream the logs here
			m := ws.NewMessage(ws.MSG_RUNNING, *source, mes.Source, "Request in progress...")
			m.Session = mes.Session

			logger.Info().Str("raw", fmt.Sprintf("%+v", m)).Msg("Sending RUNNING message")
			err := mh.SendMessage(m)
			if err == nil {
				// if the first message update worked we stop the ticker else we try again later
				ticker.Stop()
			}
		}
	}
}

func parseEnvAssignments(input string) ([]string, []string, string, error) {
	tokens, err := shlex.Split(input)

	if err != nil {
		logger.Error().Err(err).Msg("Cannot read command line")
		return nil, nil, "", err
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
	return envVars, rest, probeId, nil
}

func processFolder(src_folder, dest_folder, collector, policy, probeId string) {

	logger.Info().Str("folder", src_folder).Msg("Reading source folder")
	// Read all entries in the source directory
	entries, err := os.ReadDir(src_folder)
	if err != nil {
		logger.Error().Err(err).Str("folder", src_folder).Msg("Failed to read source folder")
	}
	uUID := uuid.New().String()
	for _, entry := range entries {
		if entry.IsDir() {
			// Skip subdirectories (you can recurse if needed)
			continue
		}
		srcPath := filepath.Join(src_folder, entry.Name())
		logger.Info().Str("file", srcPath).Msg("Processing file")

		timestamp, device, endpoint, err := parseFilename(entry.Name())
		if err != nil {
			logger.Error().Err(err).Msg("File name error")
			return
		}

		// Read file content
		content, err := os.ReadFile(srcPath)
		if err != nil {
			logger.Error().Err(err).Str("file", srcPath).Msg("Failed to read file")
		}

		policy, version := splitLast(policy)
		// Prepend namePrefix
		modifiedContent := []byte(
			"---collector:\t" + collector + "\n" +
				"---probe_id:\t" + probeId + "\n" +
				"---collection_id:\t" + uUID + "\n" +
				"---policy:\t" + policy + "\n" +
				"---version:\t" + version + "\n" +
				"---timestamp:\t" + timestamp + "\n" +
				"---device:\t" + stripAfterLast(device, ":") + "\n" +
				"---endpoint:\t" + stripAfterLast(endpoint, ".") + "\n" +
				string(content))

		// Write modified content to destination
		err = os.WriteFile(srcPath, modifiedContent, 0644)
		if err != nil {
			logger.Error().Err(err).Str("file", srcPath).Msg("Failed to tag file")
		}
		logger.Info().Str("file", srcPath).Msg("File tagged")
	}
	// now all files are tagged, zipping the whole folder
	zipFilePath := filepath.Join(src_folder, uUID+".zip")
	err = common.ZipDirFlatAndDelete(src_folder, zipFilePath)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to zip folder")
	}
	logger.Info().Msgf("ZIP created at %s", zipFilePath)

	dstPath := filepath.Join(dest_folder, uUID+".zip")
	err = os.Rename(zipFilePath, dstPath)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to rename file: %s", zipFilePath)
	}
	logger.Info().Msgf("ZIP moved to %s", dstPath)
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

func splitLast(s string) (string, string) {
	parts := strings.Split(s, "_")
	if len(parts) <= 1 {
		return "", s // No underscore or only one part
	}
	first := strings.Join(parts[:len(parts)-1], "_")
	last := parts[len(parts)-1]
	return first, last
}
