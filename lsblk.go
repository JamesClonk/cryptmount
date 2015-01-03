package main

import (
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	LSBLK_CMD = `lsblk -f -a -p --pairs -n -b -o NAME,FSTYPE,MOUNTPOINT,LABEL,UUID,PARTLABEL,PARTUUID,SIZE,TYPE`
	LSBLK_RX  = regexp.MustCompile(`(\w+)=(?:"(.*?)")`)
)

const (
	_             = iota
	KBYTE float64 = 1 << (10 * iota)
	MBYTE
	GBYTE
	TBYTE
)

type Disk struct {
	Name    string
	Size    string
	SizeH   string
	Volumes []Volume
	HasLUKS bool
}

type Volume struct {
	Name       string
	Fstype     string
	Mountpoint string
	Label      string
	Uuid       string
	Partlabel  string
	Partuuid   string
	Size       string
	SizeH      string
	Type       string
}

func lsdsk() []Disk {
	disks := make([]Disk, 0)

	var disk Disk
	for _, volume := range lsblk() {
		if volume.Type == "disk" {
			// add new disk
			disk = Disk{
				Name:    volume.Name,
				Size:    volume.Size,
				SizeH:   volume.SizeH,
				Volumes: make([]Volume, 0),
			}
			disks = append(disks, disk)
		} else {
			// add volume to disk
			if volume.Fstype == "crypto_LUKS" {
				disks[len(disks)-1].Volumes = append(disks[len(disks)-1].Volumes, volume)
				disks[len(disks)-1].HasLUKS = true
			}
		}
	}

	return disks
}

func lsblk() []Volume {
	cmd := strings.Split(LSBLK_CMD, " ")
	out, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		log.Fatal(err)
	}

	volumes := make([]Volume, 0)
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.Trim(line, "\t\n\r ")
		if line == "" {
			continue
		}

		volume := Volume{}
		vol := reflect.ValueOf(&volume).Elem()

		pairs := LSBLK_RX.FindAllStringSubmatch(line, -1)
		for _, pair := range pairs {
			if len(pair) != 3 {
				continue
			}

			key, value := strings.Title(strings.ToLower(pair[1])), pair[2]
			field := vol.FieldByName(key)
			if field.IsValid() {
				field.SetString(value)
				if key == "Size" {
					if size, err := fmtByteSize(value); err == nil {
						volume.SizeH = size
					}
				}
			}
		}

		volumes = append(volumes, volume)
	}

	return volumes
}

func fmtByteSize(value string) (string, error) {
	size, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return "", err
	}

	switch {
	case size >= TBYTE:
		return fmt.Sprintf("%.1fT", size/TBYTE), nil
	case size >= GBYTE:
		return fmt.Sprintf("%.1fG", size/GBYTE), nil
	case size >= MBYTE:
		return fmt.Sprintf("%.1fM", size/MBYTE), nil
	case size >= KBYTE:
		return fmt.Sprintf("%.1fK", size/KBYTE), nil
	}

	return fmt.Sprintf("%.1fB", size), nil
}
