package main

import (
	"fmt"
	"os"
	"slices"
)

func CliRun() {
	if slices.Contains(os.Args, "add") {
		CliAdd()
	} else {
		CliApply()
	}
}

func CliAdd() {
	var repository string
	var branch string
	var localpath string
	var domain string
	var url string
	fmt.Println("repository: ")
	fmt.Scanln(&repository)
	fmt.Println("branch: ")
	fmt.Scanln(&branch)
	fmt.Println("local path: ")
	fmt.Scanln(&localpath)
	fmt.Println("domain: ")
	fmt.Scanln(&domain)
	fmt.Println("url: ")
	fmt.Scanln(&url)
	config := ReadConfig()
	config.Sites = append(config.Sites, Site{repository, branch, localpath, domain, url})
	config.Apply()
	SaveConfig(config)
}

func CliApply() {
	config := ReadConfig()
	config.Apply()
}
