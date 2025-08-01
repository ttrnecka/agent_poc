package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

// =================================================================================
// Function: RunCommand
//
//	Execute command and return only error status
func runCommand(client *ssh.Client, command string) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	logger.Info().Msgf("Running command: %s", command)
	out, err := session.CombinedOutput(command)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func checkFolder(folder string) error {
	info, err := os.Stat(folder)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("path %s is not a directory", folder)
	}
	return nil
}

// TODO if this collector is enhanced so it pulls from more that just that on switch, the endpoint needs to be properly pointin to corect device
func generateFilename(command string) string {
	return fmt.Sprintf("%d_%s_%s.txt", time.Now().UnixMicro(), strings.SplitN(viper.GetString("endpoint"), ":", 2)[0], sanitizeCommand(command))
}

func sanitizeCommand(input string) string {
	// Replace spaces with underscores
	s := strings.ReplaceAll(input, " ", "_")

	// Build a new string keeping only valid filename characters
	var builder strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			builder.WriteRune(r)
		}
		// else drop the character
	}
	return builder.String()
}

func saveFile(folder, filename string, data []byte) error {
	filePath := filepath.Join(folder, filename)
	err := os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}
	return nil
}
