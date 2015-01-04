package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

func system(command string) string {
	cmd := strings.Split(command, " ")
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		log.Fatal(string(out), err)
	}
	return string(out)
}

func interactive(command string) {
	cmds := strings.Split(command, " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
