package main

import (
	"fmt"
	"os"
	"slices"
	"strings"
)

func CliRun() {
	commands := map[string]func(){
		"githubtoken": func() { GetToken() },
		"add":         func() { CliAdd() },
		"apply":       func() { CliApply() },
		"getports":    func() { PortScan() },
		"getfreeport": func() {
			port, err := GetFreePort()
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				fmt.Printf("Свободный порт: %d\n", port)
			}
		},
	}
	for arg, handler := range commands {
		if slices.Contains(os.Args, arg) {
			handler()
			return
		} else {
			os.Exit(0)
		}
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
