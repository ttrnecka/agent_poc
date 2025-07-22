package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

// all validation CMDs needs to have retrun code 0
var validationCmds []string = []string{
	"version",
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

	log.Printf("Validating probe %s version %s", NAME, VERSION)
	log.Printf("Connecting to host: %s", viper.GetString("endpoint"))

	client, err := connectToHost()
	if err != nil {
		log.Printf("Validation failed: %v", err)
		return err
	}
	defer client.Close()
	log.Printf("Connected to host: %s", viper.GetString("endpoint"))

	code := 0
	exErr := exitCodeError{}
	for _, cmd := range validationCmds {
		cmd_out := []byte(fmt.Sprintf(">>> %s", cmd))
		log.Printf("Running command: %s", cmd)
		out, err := runCommand(client, cmd)
		if err != nil {
			exErr.Err = err
			log.Printf("Error running command %s: %v", cmd, err)
			if exitErr, ok := err.(*ssh.ExitError); ok {
				exErr.Code = exitErr.ExitStatus()
			} else {
				exErr.Code = 255
			}
		}
		out = append(cmd_out, out...)
		fmt.Println(string(out))

		filename := genearateFilename(cmd)
		log.Printf("Saving output to file: %s", filename)
		err = saveFile(viper.GetString("output_folder"), filename, out)
		if err != nil {
			log.Printf("Validation failed: %v", err)
			return fmt.Errorf("failed to save file %s: %w", filename, err)
		}
	}
	if code != 0 {
		log.Printf("Validation failed with exit code %d\n", code)
		return &exErr
	}
	log.Println("Validation completed successfully")

	return nil
}
