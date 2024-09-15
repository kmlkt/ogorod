package main

import (
	"log"
	"os"
	"os/exec"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func StupidHandle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func LocalPath(p Site) string {
	path := "/" + p.Domain + "/" + p.URL + "/"
	return path
}

func BuildSite(p Site) {
	cmd := exec.Command(LocalPath(p) + "build.sh")
	_, err := cmd.Output()
	StupidHandle(err)
}

func CheckBuild(p Site) {
	if _, err := os.Stat(LocalPath(p) + "build.sh"); os.IsNotExist(err) {
		BuildSite(p)
	} else {
		return
	}
}
