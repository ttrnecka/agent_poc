package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/ssh"
)

var (
	Version = "1.0.0"
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

func main() {
	if len(os.Args) != 4 {
		log.Fatalf("Usage: %s <user> <password> <host:port>", os.Args[0])
	}

	client, err := connectToHost(os.Args[1], os.Args[2], os.Args[3])
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to host: %v", err))
	}

	var commands []string
	switch Version {
	case "1.0.0":
		commands = commandsV100
	case "1.0.1":
		commands = commandsV101
	default:
		log.Fatal(fmt.Errorf("unknown build %s", Version))
	}
	for _, cmd := range commands {
		out, err := runCommand(client, cmd)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
	}
	client.Close()
}

func connectToHost(user, pass, host string) (*ssh.Client, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{ssh.Password(pass)},
	}
	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	client, err := ssh.Dial("tcp", host, sshConfig)
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
		client.Close()
		return "", err
	}
	defer session.Close()

	out, err := session.CombinedOutput(command)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
