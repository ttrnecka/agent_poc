package main

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
)

type exitCodeError struct {
	Code int
	Err  error
}

func (e *exitCodeError) Error() string {
	return e.Error()
}

func (e *exitCodeError) Unwrap() error {
	return e.Err
}

func connectToHost() (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: viper.GetString("user"),
		Auth: []ssh.AuthMethod{ssh.Password(viper.GetString("password"))},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", viper.GetString("endpoint"), sshConfig)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// =================================================================================
// Function: RunCommand
//
//	Execute command and return only error status
func runCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	out, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
