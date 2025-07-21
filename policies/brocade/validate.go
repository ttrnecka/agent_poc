package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// all validation CMDs needs to have retrun code 0
var validationCmds []string = []string{
	"switchshow",
}

// runCmd represents the run command
func NewValidateCmd() *cobra.Command {

	var cmd = &cobra.Command{
		Use:          "validate",
		Short:        "validate access",
		Long:         `validate access`,
		RunE:         validate,
		SilenceUsage: true,
	}

	return cmd
}

func validate(cmd *cobra.Command, args []string) error {

	client, err := connectToHost()
	if err != nil {
		return err
	}
	defer client.Close()

	code := 0
	for _, cmd := range validationCmds {
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
