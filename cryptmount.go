package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

const VERSION = "0.0.1"

func main() {
	app := cli.NewApp()
	app.Name = "cryptmount"
	app.Author = "JamesClonk"
	app.Email = "jamesclonk@jamesclonk.ch"
	app.Version = VERSION
	app.Usage = "a tool for easy mounting of LUKS encrypted volumes"

	app.Commands = []cli.Command{{
		Name:        "mount",
		ShortName:   "m",
		Usage:       "mount an encrypted volumes",
		Description: ".....", // TODO: add description
		Action: func(c *cli.Context) {
			println("mount: ", c.Args().First())
		},
	}, {
		Name:        "unmount",
		ShortName:   "u",
		Usage:       "unmount an encrypted volumes",
		Description: ".....", // TODO: add description
		Action: func(c *cli.Context) {
			println("unmount: ", c.Args().First())
		},
	}, {
		Name:        "unmount-all",
		ShortName:   "ua",
		Usage:       "unmount all encrypted volumes",
		Description: ".....", // TODO: add description
		Action: func(c *cli.Context) {
			println("unmount-all: ", c.Args().First())
		},
	}, {
		Name:        "list",
		ShortName:   "l",
		Usage:       "list all encrypted volumes",
		Description: ".....", // TODO: add description
		Action: func(c *cli.Context) {
			fmt.Printf("%q\n", lsblk())
		},
	}}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
	app.RunAndExitOnError()
}
