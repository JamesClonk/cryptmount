package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	LUKSOPEN_CMD  = `cryptsetup luksOpen %v %v`
	LUKSCLOSE_CMD = `cryptsetup luksClose %v`
	MOUNT_CMD     = `mount /dev/mapper/%v %v`
	UNMOUNT_CMD   = `umount %v`
)

func mountVolume(device string, mountpoint string) {
	fmt.Printf("\nMount [%v] to [%v]\n", magenta(device), magenta(mountpoint))

	name := mapperName(mountpoint)
	luksOpen(device, name)
	mount(name, mountpoint)
}

func unmountVolume(device string, mountpoint string) {
	fmt.Printf("\nUnmount device [%v] from [%v]\n", magenta(device), magenta(mountpoint))

	unmount(mountpoint)
	luksClose(mapperName(mountpoint))
}

func luksOpen(device string, name string) {
	// check if mapper name already exists, if so then close it
	if _, err := os.Stat("/dev/mapper/" + name); err != nil {
		if !os.IsNotExist(err) {
			log.Fatal(err)
		}
	} else {
		luksClose(name)
	}

	interactive(fmt.Sprintf(LUKSOPEN_CMD, device, name))
}

func luksClose(name string) {
	interactive(fmt.Sprintf(LUKSCLOSE_CMD, name))
}

func mount(name string, mountpoint string) {
	// check if mountpoint directory exists, if not then create it
	if _, err := os.Stat(mountpoint); os.IsNotExist(err) {
		if err := os.Mkdir(mountpoint, 0750); err != nil {
			log.Fatal(err)
		}
	}

	interactive(fmt.Sprintf(MOUNT_CMD, name, mountpoint))
}

func unmount(mountpoint string) {
	interactive(fmt.Sprintf(UNMOUNT_CMD, mountpoint))
}

func mapperName(path string) string {
	return strings.Replace(strings.TrimLeft(path, "/"), "/", "_", -1)
}
