package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func CliRun() {
	if slices.Contains(os.Args, "githubtoken") {
		GetToken()
	}
	if slices.Contains(os.Args, "add") {
		CliAdd()
	}
	if slices.Contains(os.Args, "apply") {
		CliApply()
	}
	if slices.Contains(os.Args, "getports") {
		PortScan()
	}
	if slices.Contains(os.Args, "getfreeport") {
		port, err := GetFreePort()
		if err != nil {
			fmt.Println("Ошибка:", err)
		} else {
			fmt.Printf("Свободный порт: %d\n", port)
		}
	} else {
		os.Exit(1)
	}
}

func CliAdd() {
	var repository string
	var branch string
	var domain string
	var url string
	var siteaddress string
	fmt.Println("repository: ")
	fmt.Scanln(&repository)
	fmt.Println("branch: ")
	fmt.Scanln(&branch)
	fmt.Println("site address: ")
	fmt.Scanln(&siteaddress)
	slicesiteaddress := strings.Split(siteaddress, "/")
	domain = slicesiteaddress[0]
	for i := 1; i < len(slicesiteaddress); i++ {
		if i == len(slicesiteaddress)-1 {
			url = url + slicesiteaddress[i]
		} else {
			url = url + slicesiteaddress[i] + "/"
		}
	}
	if !strings.HasPrefix(repository, "https://") && !strings.HasPrefix(repository, "http://") {
		repository = "https://" + repository
	}
	config := ReadConfig()
	config.Sites = append(config.Sites, Site{repository, branch, domain, url})
	config.Apply()
	SaveConfig(config)
}

func CliApply() {
	config := ReadConfig()
	config.Apply()
}
