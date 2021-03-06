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

const VERSION = "1.0.0"

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
			mountAllDevices(c)
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
			unmountAllDevices(c)
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
	volumes := listDevices()
	unmountedVolumes := *volumes.unmounted()

	if len(unmountedVolumes) > 0 {
		fmt.Printf("\n******************************************************\n")
		fmt.Printf("Which volume to mount?\n\n")
		for idx, volume := range unmountedVolumes {
			fmt.Printf("%v: %v\n", red(strconv.Itoa(idx)), bold(volume.Name))
		}

		choice := chooseVolume()
		mountVolume(unmountedVolumes[choice])
	} else {
		fmt.Printf("%v\n", boldRedBlinking("No volumes to mount were found"))
		os.Exit(7)
	}
}

func mountAllDevices(c *cli.Context) {
	volumes := listDevices()
	for _, volume := range *volumes.unmounted() {
		mountVolume(volume)
	}
}

func unmountDevice(c *cli.Context) {
	volumes := listDevices()
	mountedVolumes := *volumes.mounted()

	if len(mountedVolumes) > 0 {
		fmt.Printf("\n******************************************************\n")
		fmt.Printf("Which volume to unmount?\n\n")
		for idx, volume := range mountedVolumes {
			fmt.Printf("%v: %v\n", red(strconv.Itoa(idx)), bold(volume.Name))
		}

		choice := chooseVolume()
		unmountVolume(mountedVolumes[choice])
	} else {
		fmt.Printf("%v\n", boldRedBlinking("No mounted volumes were found"))
		os.Exit(8)
	}
}

func unmountAllDevices(c *cli.Context) {
	volumes := listDevices()
	for _, volume := range *volumes.mounted() {
		unmountVolume(volume)
	}
}

func listDevices() (result Volumes) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 2, '\t', 0)

	for _, disk := range Lsdsk() {
		fmt.Fprintf(w, "Disk: %v\tSize: %v\n", bold(disk.Name), bold(disk.SizeH))
		w.Flush()

		if !disk.HasLUKS {
			fmt.Printf("   %v\n", "No LUKS partition found")
		} else {
			volumes := disk.Volumes.luksOnly()
			for idx, volume := range *volumes {
				if idx != len(*volumes)-1 {
					fmt.Print("  ├── ")
				} else {
					fmt.Print("  └── ")
				}

				mappedName := magenta(volume.MappedName)
				if !volume.IsMapped {
					mappedName = blink(magenta("not mapped"))
				}
				fmt.Fprintf(w, "Partition: %v\tSize: %v\tMountpoint: %v\tUUID: %v\tLabel: %v\n",
					magenta(volume.Name),
					magenta(volume.SizeH),
					mappedName,
					magenta(volume.Uuid),
					magenta(volume.Label))
				w.Flush()

				// check for mapped volumes
				if len(volume.MappedVolumes) > 0 {
					for midx, mappedVolume := range volume.MappedVolumes {
						if midx != len(volume.MappedVolumes)-1 {
							fmt.Print("   ├─ ")
						} else {
							fmt.Print("   └─ ")
						}

						mountpoint := magenta(mappedVolume.Mountpoint)
						if !mappedVolume.IsMounted {
							mountpoint = blink(magenta("not mounted"))
						}
						fmt.Fprintf(w, "Partition: %v\tSize: %v\tMountpoint: %v\tUUID: %v\tLabel: %v\n",
							magenta(mappedVolume.Name),
							magenta(mappedVolume.SizeH),
							mountpoint,
							magenta(mappedVolume.Uuid),
							magenta(mappedVolume.Label))
						w.Flush()
					}
				}

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
