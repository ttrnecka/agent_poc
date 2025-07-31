package core

import (
	"fmt"
	"io"

	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type ExitCodeError struct {
	Code int
	Err  error
}

func (e *ExitCodeError) Error() string {
	return e.Err.Error()
}

func (e *ExitCodeError) Unwrap() error {
	return e.Err
}

type SshRunner struct {
	client *ssh.Client
}

func (s *SshRunner) Connect() (io.Closer, error) {
	client, err := connectToHost()
	if err != nil {
		return nil, err
	}
	s.client = client
	return s.client, nil
}

func (s *SshRunner) Run(cmd string) ([]byte, *ExitCodeError) {
	exErr := ExitCodeError{}
	out, err := runCommand(s.client, cmd)
	if err != nil {
		exErr.Err = err
		logger.Error().Err(err).Str("command", cmd).Msg("Error running command")
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exErr.Code = exitErr.ExitStatus()
		} else {
			exErr.Code = 255
		}
	}
	logger.Debug().Str("output", string(out)).Msg("Command output")

	filename := generateFilename(cmd)
	logger.Info().Str("file", filename).Msg("Saving output to file")
	err = saveFile(viper.GetString("output_folder"), filename, out)
	if err != nil {
		logger.Error().Err(err).Msg("Validation failed")
		exErr.Code = 255
		exErr.Err = fmt.Errorf("failed to save file %s: %w", filename, err)
	}
	return out, &exErr
}

func connectToHost() (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: viper.GetString("user"),
		Auth: []ssh.AuthMethod{ssh.Password(viper.GetString("password"))},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	logger.Info().Str("host", viper.GetString("endpoint")).Msg("Connecting to host")

	client, err := ssh.Dial("tcp", viper.GetString("endpoint"), sshConfig)
	if err != nil {
		return nil, err
	}
	logger.Info().Str("host", viper.GetString("endpoint")).Msg("Connected to host")
	return client, nil
}
