package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var commandsV100 []string = []string{
	"switchshow",
	"version",
}

var commandsV101 []string = []string{
	"switchshow",
	"version",
	"fabricshow",
	"agshow",
}

// runCmd represents the run command
func NewCollectCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "collect",
		Short:        "collect data",
		Long:         `collectdata`,
		RunE:         collect,
		SilenceUsage: true,
	}
	return cmd
}

func collect(cmd *cobra.Command, args []string) error {

	client, err := connectToHost()
	if err != nil {
		return err
	}
	defer client.Close()

	var commands []string
	switch VERSION {
	case "1.0.0":
		commands = commandsV100
	case "1.0.1":
		commands = commandsV101
	default:
		return fmt.Errorf("unknown build %s", VERSION)
	}
	code := 0
	for _, cmd := range commands {
		out, err := runCommand(client, cmd)
		if err != nil {
			log.Printf("Error running command %s: %v", cmd, err)
			if exitErr, ok := err.(*ssh.ExitError); ok {
				code = exitErr.ExitStatus()
			} else {
				code = 255
			}
		}
		fmt.Println(string(out))
	}
	if code != 0 {
		return &exitCodeError{Code: code, Err: err}
	}

	return nil
}
