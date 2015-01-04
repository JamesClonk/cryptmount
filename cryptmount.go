package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/codegangsta/cli"
	"github.com/fatih/color"
)

const VERSION = "0.0.1"

var (
	blink           = color.New(color.BlinkSlow).SprintfFunc()
	bold            = color.New(color.Bold).SprintfFunc()
	magenta         = color.New(color.FgMagenta).SprintfFunc()
	red             = color.New(color.FgRed).SprintfFunc()
	boldRedBlinking = color.New(color.BlinkSlow).Add(color.Bold).Add(color.FgRed).SprintfFunc()
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
			mountDevice(c)
		},
	}, {
		Name:        "mount-all",
		ShortName:   "ma",
		Usage:       "mount all encrypted volumes",
		Description: ".....", // TODO: add description
		Action: func(c *cli.Context) {
			//mountDevices(c)
		},
	}, {
		Name:        "unmount",
		ShortName:   "u",
		Usage:       "unmount an encrypted volumes",
		Description: ".....", // TODO: add description
		Action: func(c *cli.Context) {
			unmountDevice(c)
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
			listDevices()
		},
	}}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
	app.RunAndExitOnError()
}

func mountDevice(c *cli.Context) {
	if !c.Args().Present() {
		volumes := listDevices()
		unmountedVolumes := *volumes.unmounted()

		if len(unmountedVolumes) > 0 {
			fmt.Printf("\n******************************************************\n")
			fmt.Printf("Which volume to mount?\n\n")
			for idx, volume := range unmountedVolumes {
				fmt.Printf("%v: %v\n", red(strconv.Itoa(idx)), bold(volume.Name))
			}

			choice := chooseVolume()
			fmt.Println(choice)

			volume := unmountedVolumes[choice]
			mountpoint := "/mnt/" + strings.Replace(strings.TrimLeft(volume.Name, "/dev/"), "/", "_", -1)
			mountVolume(volume.Name, mountpoint)
		} else {
			fmt.Printf("%v\n", boldRedBlinking("No volumes to mount were found"))
			os.Exit(7)
		}
	} else {
		fmt.Println("mount directly", c.Args()) // TODO: mount directly with given arguments
	}
}

func unmountDevice(c *cli.Context) {
	if !c.Args().Present() {
		volumes := listDevices()
		mountedVolumes := *volumes.mounted()

		if len(mountedVolumes) > 0 {
			fmt.Printf("\n******************************************************\n")
			fmt.Printf("Which volume to unmount?\n\n")
			for idx, volume := range mountedVolumes {
				fmt.Printf("%v: %v\n", red(strconv.Itoa(idx)), bold(volume.Name))
			}

			choice := chooseVolume()
			fmt.Println(choice)

			volume := mountedVolumes[choice]
			mountpoint := "/mnt/" + strings.Replace(strings.TrimLeft(volume.Name, "/dev/"), "/", "_", -1)
			unmountVolume(volume.Name, mountpoint)
		} else {
			fmt.Printf("%v\n", boldRedBlinking("No mounted volumes were found"))
			os.Exit(8)
		}
	} else {
		fmt.Println("unmount directly", c.Args()) // TODO: unmount directly with given arguments
	}
}

func listDevices() (result Volumes) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	for _, disk := range lsdsk() {
		fmt.Fprintf(w, "Disk: %v\tSize: %v\n", bold(disk.Name), bold(disk.SizeH))
		w.Flush()

		if !disk.HasLUKS {
			fmt.Printf("   %v\n", blink("No LUKS partition found"))
		} else {
			volumes := disk.Volumes.luksOnly()
			for idx, volume := range *volumes {
				if idx != len(*volumes)-1 {
					fmt.Print("  ├── ")
				} else {
					fmt.Print("  └── ")
				}

				mountpoint := magenta(volume.Mountpoint)
				if !volume.IsMounted {
					mountpoint = blink(magenta("not mounted"))
				}
				fmt.Fprintf(w, "Partition: %v\tSize: %v\tMountpoint: %v\tUUID: %v\tLabel: %v\n",
					magenta(volume.Name),
					magenta(volume.SizeH),
					mountpoint,
					magenta(volume.Uuid),
					magenta(volume.Label))
				w.Flush()

				result = append(result, volume)
			}
		}
		fmt.Println()
	}

	return
}

func chooseVolume() (result int64) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nVolume? [0]: ")
	choice, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	choice = strings.Trim(choice, "\t\n\r ")
	if choice != "" {
		result, err = strconv.ParseInt(choice, 10, 64)
		if err != nil {
			fmt.Printf("\n%v: [%v]\n\n", red("Not a valid volume choice"), boldRedBlinking(choice))
			os.Exit(5)
		}
	}

	return
}
