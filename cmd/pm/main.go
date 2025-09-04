package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	sshHost   = flag.String("ssh-host", "localhost:22", "SSH host:port")
	sshUser   = flag.String("ssh-user", "user", "SSH username")
	sshPass   = flag.String("ssh-pass", "", "SSH password")
	remoteDir = flag.String("remote-dir", "/packages", "Remote directory for packages")
)

func main() {
	flag.Parse()

	if len(os.Args) < 3 {
		usage()
		return
	}

	cmd := os.Args[1]
	//filePath := os.Args[2]

	switch cmd {
	case "create":
		fmt.Println("Create command not implemented yet")
	case "update":
		fmt.Println("Update command not implemented yet")
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: pm [create|update] <config.json|config.yaml>")
	flag.PrintDefaults()
}
