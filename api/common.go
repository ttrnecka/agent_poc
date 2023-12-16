package api

import (
	"fmt"
	"os"
)

func output(file string) string {
	cmd_file := fmt.Sprintf("data/%s", file)
	if _, err := os.Stat(cmd_file); os.IsNotExist(err) {
		return fmt.Sprintf("error: file %s does not exist", cmd_file)
	}

	b, err := os.ReadFile(cmd_file) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return string(b)
}
