package main

import (
	"fmt"

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
	Volumes Volumes
	HasLUKS bool
}

type Volume struct {
	Name       string
	Fstype     string
	IsLUKS     bool
	Mountpoint string
	IsMounted  bool
	Label      string
	Uuid       string
	Partlabel  string
	Partuuid   string
	Size       string
	SizeH      string
	Type       string
}

type Volumes []Volume

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
				disks[len(disks)-1].HasLUKS = true
				volume.IsLUKS = true
			}
			disks[len(disks)-1].Volumes = append(disks[len(disks)-1].Volumes, volume)
		}
	}

	return disks
}

func lsblk() (volumes Volumes) {
	for _, line := range strings.Split(system(LSBLK_CMD), "\n") {
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
					if size, err := formatByteSize(value); err == nil {
						volume.SizeH = size
					}
				}
			}
		}

		volume.IsMounted = volume.Mountpoint != ""
		volumes = append(volumes, volume)
	}

	return
}

func formatByteSize(value string) (string, error) {
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

func (volumes *Volumes) filter(predicate func(Volume) bool) *Volumes {
	var newVolumes Volumes
	for _, volume := range *volumes {
		if predicate(volume) {
			newVolumes = append(newVolumes, volume)
		}
	}
	return &newVolumes
}

func (volumes *Volumes) luksOnly() *Volumes {
	return volumes.filter(func(v Volume) bool {
		return v.IsLUKS
	})
}

func (volumes *Volumes) mounted() *Volumes {
	return volumes.filter(func(v Volume) bool {
		return v.IsMounted
	})
}

func (volumes *Volumes) unmounted() *Volumes {
	return volumes.filter(func(v Volume) bool {
		return !v.IsMounted
	})
}
