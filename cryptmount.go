package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
		usage()
		os.Exit(1)
	}
}

func usage() {
	usage := fmt.Sprintf("usage: %s <path>\n", filepath.Base(os.Args[0]))
	fmt.Printf(usage)
}
