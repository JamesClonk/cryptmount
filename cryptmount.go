package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

const VERSION = "0.0.1"

var (
	blink   = color.New(color.BlinkSlow).SprintfFunc()
	bold    = color.New(color.Bold).SprintfFunc()
	magenta = color.New(color.FgMagenta).SprintfFunc()
)

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
			w := new(tabwriter.Writer)
			w.Init(os.Stdout, 0, 8, 2, '\t', 0)

			for _, disk := range lsdsk() {
				fmt.Fprintf(w, "Disk: %v\tSize: %v\n", bold(disk.Name), bold(disk.SizeH))
				w.Flush()

				if !disk.HasLUKS {
					fmt.Printf("   %v\n", blink("No LUKS partition found"))
				} else {
					for idx, volume := range disk.Volumes {
						if idx != len(disk.Volumes)-1 {
							fmt.Print("  ├── ")
						} else {
							fmt.Print("  └── ")
						}

						mountpoint := magenta(volume.Mountpoint)
						if volume.Mountpoint == "" {
							mountpoint = blink(magenta("not mounted"))
						}
						fmt.Fprintf(w, "Partition: %v\tSize: %v\tMountpoint: %v\tUUID: %v\tLabel: %v\n",
							magenta(volume.Name),
							magenta(volume.SizeH),
							mountpoint,
							magenta(volume.Uuid),
							magenta(volume.Label))
						w.Flush()
					}
				}
				fmt.Println()
			}
		},
	}}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
	app.RunAndExitOnError()
}
