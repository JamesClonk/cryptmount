package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	LUKSOPEN_CMD = `cryptsetup luksOpen %v %v`
	MOUNT_CMD    = `mount /dev/mapper/%v %v`
)

func mountVolume(device string, mountpoint string) {
	fmt.Printf("\nMount [%v] to [%v]\n", magenta(device), magenta(mountpoint))

	name := mapperName(mountpoint)
	luksOpen(device, name)
	mount(name, mountpoint)
}

func luksOpen(device string, name string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nEnter passphrase for %v:", device)
	passphrase, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	passphrase = strings.Trim(passphrase, "\t\n\r ")
	if passphrase == "" {
		log.Fatalf("%v", boldRedBlinking("No passphrase given!"))
	}

	cmds := strings.Split(fmt.Sprintf(LUKSOPEN_CMD, device, name), " ")
	cmd := exec.Command(cmds[0], cmds[1:]...)
	//cmd := exec.Command("cryptsetup", "luksOpen", "/dev/sdc2", "sdc2")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	stdin, _ := cmd.StdinPipe()
	stdin.Write([]byte(passphrase + "\n"))

	time.AfterFunc(time.Duration(2)*time.Second, func() {
		if !cmd.ProcessState.Exited() {
			log.Fatalf("\n%v\n", boldRedBlinking("Incorrect passphrase!"))
		}
	})

	cmd.Run()
}

func mount(name string, mountpoint string) {
	out := system(fmt.Sprintf(MOUNT_CMD, name, mountpoint))
	fmt.Println(out)
}

func mapperName(path string) string {
	return strings.Replace(strings.TrimLeft(path, "/"), "/", "_", -1)
}
