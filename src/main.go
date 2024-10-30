package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

type CommandStruct struct {
	command        string
	isDescoveryCmd bool
	commandArray   []string
}

func executeCmd(s *IBMStorageFS5K, command string) string {
	config := &ssh.ClientConfig{
		User:            s.user,
		Auth:            []ssh.AuthMethod{ssh.Password(s.password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	conn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%s", s.hostname, s.port), config)
	if err != nil {
		log.Fatalf("Failed to create Dial: %v", err)
	}
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()
	var stdoutBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Run(command)

	// resultChan <- fmt.Sprintf("%s -> %s", s.hostname, stdoutBuf.String())
	return stdoutBuf.String()
}

func main() {

	// Validate input arguments
	if len(os.Args) < 5 {
		log.Fatalln("Usage: ./ibm-storage-collect <hosts> <user> <password> <command>")
		return
	}

	host := os.Args[1]
	user := os.Args[2]
	port := "22"
	password := os.Args[3]
	commandArray := os.Args[4:]
	commandString := strings.Join(commandArray, " ")

	// Create a channel to store the results
	// results := make(chan string, 10)
	// Create a timeout channel
	// timeout := time.After(5 * time.Second)
	// Create a channel to store the command

	// Create a command object
	cmd := CommandStruct{
		command:        commandString,
		isDescoveryCmd: (len(commandArray) == 1),
		commandArray:   commandArray,
	}

	// Create IBMStorageFS5K object
	IBMStorageFS5K := IBMStorageFS5K{
		hostname: host,
		port:     port,
		user:     user,
		password: password,
	}

	// Switch command to execute the correct command
	IBMStorageFS5K.handleCommand(cmd)

}
