package main

import (
	"log"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
)

var (
	LSBLK_CMD = `lsblk -f -a -p --pairs -n -b -o NAME,FSTYPE,MOUNTPOINT,LABEL,UUID,PARTLABEL,PARTUUID,SIZE,TYPE`
	LSBLK_RX  = regexp.MustCompile(`(\w+)=(?:"(.*?)")`)
)

type Volume struct {
	Name       string
	Fstype     string
	Mountpoint string
	Label      string
	Uuid       string
	Partlabel  string
	Partuuid   string
	Size       string
	Type       string
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
			field.SetString(value)
		}

		volumes = append(volumes, volume)
	}

	return volumes
}
