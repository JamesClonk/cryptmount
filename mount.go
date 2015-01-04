package main

import (
	"fmt"

	"strings"
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
	out := system(fmt.Sprintf(LUKSOPEN_CMD, device, name))
	fmt.Println(out)
}

func mount(name string, mountpoint string) {
	out := system(fmt.Sprintf(MOUNT_CMD, name, mountpoint))
	fmt.Println(out)
}

func mapperName(path string) string {
	return strings.Replace(strings.TrimLeft(path, "/"), "/", "_", -1)
}
