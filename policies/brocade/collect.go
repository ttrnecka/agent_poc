package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

var commandsV100 []string = []string{
	"version",
}

var commandsV101 []string = []string{
	"version",
	"switchshow",
}

var commandsV102 []string = []string{
	"version",
	"switchshow",
	"fabricshow",
	"licenseshow",
}

var commandsV103 []string = []string{
	"version",
	"switchshow",
	"fabricshow",
	"license --show",
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

	log.Printf("Collecting data for %s version %s", NAME, VERSION)
	log.Printf("Connecting to host: %s", viper.GetString("endpoint"))
	client, err := connectToHost()
	if err != nil {
		return err
	}
	defer client.Close()
	log.Printf("Connected to host: %s", viper.GetString("endpoint"))

	var commands []string
	switch VERSION {
	case "1.0.0":
		commands = commandsV100
	case "1.0.1":
		commands = commandsV101
	case "1.0.2":
		commands = commandsV102
	case "1.0.3":
		commands = commandsV103
	default:
		return fmt.Errorf("unknown build %s", VERSION)
	}

	exErr := exitCodeError{}
	// calls each command and saves each output to tagged file
	for _, cmd := range commands {
		cmd_out := []byte(fmt.Sprintf(">>> %s", cmd))
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
			return fmt.Errorf("failed to save file %s: %w", filename, err)
		}
	}
	if exErr.Code != 0 {
		return &exErr
	}

	return nil
}
