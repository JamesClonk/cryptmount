package main_test

import (
	. "github.com/JamesClonk/cryptmount"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Lsblk", func() {

	var (
		expectedOutput = []Disk{
			Disk{
				Name:  "/dev/sda",
				Size:  "340054409536",
				SizeH: "316.7G",
				Volumes: Volumes{
					Volume{
						Name:          "/dev/sda1",
						Fstype:        "ntfs",
						IsLUKS:        false,
						MappedName:    "",
						IsMapped:      false,
						MappedVolumes: nil,
						Mountpoint:    "",
						IsMounted:     false,
						Label:         "System Reserved",
						Uuid:          "B89A79380A19B935",
						Size:          "272188000",
						SizeH:         "259.6M",
						Type:          "part",
					},
					Volume{
						Name:          "/dev/sda2",
						Fstype:        "ext4",
						IsLUKS:        false,
						MappedName:    "",
						IsMapped:      false,
						MappedVolumes: nil,
						Mountpoint:    "/",
						IsMounted:     true,
						Label:         "LINUX",
						Uuid:          "5673c8b2-ef3c-4cbd-bc31-ef045fc12731",
						Size:          "249794218880",
						SizeH:         "232.6G",
						Type:          "part",
					},
				},
				HasLUKS: false,
			},
			Disk{
				Name:  "/dev/sdb",
				Size:  "340054409536",
				SizeH: "316.7G",
				Volumes: Volumes{
					Volume{
						Name:          "/dev/sdb1",
						Fstype:        "ntfs",
						IsLUKS:        false,
						MappedName:    "",
						IsMapped:      false,
						MappedVolumes: nil,
						Mountpoint:    "",
						IsMounted:     false,
						Label:         "",
						Uuid:          "123CB6822FE24C99",
						Size:          "280054697777",
						SizeH:         "260.8G",
						Type:          "part",
					},
				},
				HasLUKS: false,
			},
			Disk{
				Name:  "/dev/sdc",
				Size:  "5000111111888",
				SizeH: "4.5T",
				Volumes: Volumes{
					Volume{
						Name:          "/dev/sdc1",
						Fstype:        "ntfs",
						IsLUKS:        false,
						MappedName:    "",
						IsMapped:      false,
						MappedVolumes: nil,
						Mountpoint:    "",
						IsMounted:     false,
						Label:         "DATA",
						Uuid:          "BD34CDFFFFCDBFFF",
						Size:          "29999999999",
						SizeH:         "27.9G",
						Type:          "part",
					},
					Volume{
						Name:       "/dev/sdc2",
						Fstype:     "crypto_LUKS",
						IsLUKS:     true,
						MappedName: "/dev/mapper/sdc2",
						IsMapped:   true,
						MappedVolumes: []MappedVolume{
							MappedVolume{
								Name:       "/dev/mapper/sdc2",
								Fstype:     "ext4",
								Mountpoint: "/mnt/sdc2",
								IsMounted:  true,
								Label:      "",
								Uuid:       "99999999-9999-1111-4444-444455556666",
								Size:       "100000888888",
								SizeH:      "93.1G",
								Type:       "crypt",
							},
						},
						Mountpoint: "",
						IsMounted:  true,
						Label:      "",
						Uuid:       "99f8879f-9999-1111-2222-ff699a9f9b98",
						Size:       "100000099000",
						SizeH:      "93.1G",
						Type:       "part",
					},
					Volume{
						Name:          "/dev/sdc3",
						Fstype:        "crypto_LUKS",
						IsLUKS:        true,
						MappedName:    "",
						IsMapped:      false,
						MappedVolumes: nil,
						Mountpoint:    "",
						IsMounted:     false,
						Label:         "",
						Uuid:          "99f887f-9999-2222-3333-ff699a9f9b98",
						Size:          "200000099000",
						SizeH:         "186.3G",
						Type:          "part",
					},
					Volume{
						Name:          "/dev/sdc4",
						Fstype:        "ext4",
						IsLUKS:        false,
						MappedName:    "",
						IsMapped:      false,
						MappedVolumes: nil,
						Mountpoint:    "/media/external_one",
						IsMounted:     true,
						Label:         "DATA1",
						Uuid:          "999c899-9999-1111-2222-96c444493333",
						Size:          "700204248576",
						SizeH:         "652.1G",
						Type:          "part",
					},
				},
				HasLUKS: true,
			},
			Disk{
				Name:  "/dev/sdd",
				Size:  "1000999999999",
				SizeH: "932.3G",
				Volumes: Volumes{
					Volume{
						Name:          "/dev/sdd1",
						Fstype:        "crypto_LUKS",
						IsLUKS:        true,
						MappedName:    "",
						IsMapped:      false,
						MappedVolumes: nil,
						Mountpoint:    "",
						IsMounted:     false,
						Label:         "",
						Uuid:          "9999a999-9999-2222-1111-7777000777ef",
						Size:          "100000000001",
						SizeH:         "93.1G",
						Type:          "part",
					},
					Volume{
						Name:          "/dev/sdd2",
						Fstype:        "ext4",
						IsLUKS:        false,
						MappedName:    "",
						IsMapped:      false,
						MappedVolumes: nil,
						Mountpoint:    "/media/external_two",
						IsMounted:     true,
						Label:         "DATA2",
						Uuid:          "9999a999-1111-1111-1111-7777000777ef",
						Size:          "90099999",
						SizeH:         "85.9M",
						Type:          "part",
					},
				},
				HasLUKS: true,
			},
		}
	)

	BeforeEach(func() {
		LSBLK_CMD = `cat fixtures/lsblk.out`
	})

	Describe("Calling Lsdsk()", func() {
		Context("With all possible volume combinations", func() {
			It("should return the expected output", func() {
				Expect(Lsdsk()).To(Equal(expectedOutput))
			})
		})
	})

})
