package main

import (
	"log"
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
