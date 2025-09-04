package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"pm/config"
	"pm/packer"
	"pm/sshclient"
	"pm/version"
	"strings"
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
	filePath := os.Args[2]

	switch cmd {
	case "create":
		createPackage(filePath)
	case "update":
		updatePackages(filePath)
	default:
		usage()
	}
}

func usage() {
	fmt.Println("Usage: pm [create|update] <config.json|config.yaml>")
	flag.PrintDefaults()
	os.Exit(1)
}

func createPackage(configFile string) {
	cfg, err := config.LoadPacketConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	files, err := packer.CollectFiles(cfg.Targets)
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
		os.Exit(1)
	}

	zipBuf, err := packer.CreateZip(files)
	if err != nil {
		fmt.Printf("Error creating ZIP: %v\n", err)
		os.Exit(1)
	}

	client, err := sshclient.NewClient(*sshHost, *sshUser, *sshPass)
	if err != nil {
		fmt.Printf("Error connecting to SSH: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	pkgName := fmt.Sprintf("%s-%s.zip", cfg.Name, cfg.Ver)
	remotePath := filepath.Join(*remoteDir, pkgName)
	err = client.UploadFile(remotePath, zipBuf)
	if err != nil {
		fmt.Printf("Error uploading package: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Пакет %s загружен на %s\n", pkgName, remotePath)
}

func updatePackages(configFile string) {
	cfg, err := config.LoadPackagesConfig(configFile)
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}

	client, err := sshclient.NewClient(*sshHost, *sshUser, *sshPass)
	if err != nil {
		fmt.Printf("Error connecting to SSH: %v\n", err)
		os.Exit(1)
	}
	defer client.Close()

	for _, dep := range cfg.Packages {
		files, err := client.ListFiles(*remoteDir)
		if err != nil {
			fmt.Printf("Error listing remote files: %v\n", err)
			os.Exit(1)
		}

		candidates := []string{}
		for _, file := range files {
			if !file.IsDir() && strings.HasPrefix(file.Name(), dep.Name+"-") && strings.HasSuffix(file.Name(), ".zip") {
				ver := strings.TrimSuffix(strings.TrimPrefix(file.Name(), dep.Name+"-"), ".zip")
				if version.MatchesVersion(dep.Ver, ver) {
					candidates = append(candidates, ver)
				}
			}
		}

		if len(candidates) == 0 {
			fmt.Printf("Нет подходящей версии для %s %s\n", dep.Name, dep.Ver)
			continue
		}

		version.SortVersions(candidates)
		highest := candidates[len(candidates)-1]
		pkgName := fmt.Sprintf("%s-%s.zip", dep.Name, highest)
		remotePath := filepath.Join(*remoteDir, pkgName)
		localPath := pkgName

		err = client.DownloadFile(remotePath, localPath)
		if err != nil {
			fmt.Printf("Error downloading %s: %v\n", pkgName, err)
			continue
		}

		err = packer.Unzip(localPath, ".")
		if err != nil {
			fmt.Printf("Error unzipping %s: %v\n", pkgName, err)
			continue
		}

		fmt.Printf("Обновлено %s до %s\n", dep.Name, highest)
	}
}
