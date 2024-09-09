package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

type Site struct {
	Repository string
	LocalPath  string
	Domain     string
	Url        string
}

var sitesEnabled = "/etc/nginx/sites-enabled/ogorod"

func CreateNginxConfig(sites []Site) {
	servers := map[string][]Site{}
	for _, site := range sites {
		servers[site.Domain] = append(servers[site.Domain], site)
	}
	file, err := os.OpenFile(sitesEnabled, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for serverName, serverSites := range servers {
		fmt.Fprintf(file, `
server {
	listen 80;
	server_name %s;
	index index.html;
		`, serverName)
		for _, site := range serverSites {
			fmt.Fprintf(file, `
	location %s {
		alias %s;
		try_files $uri $uri/ =404;
	}
`, site.Url, site.LocalPath)
		}
		fmt.Fprint(file, "}")
	}
}

func RestartNginx() {
	cmd := exec.Command("systemctl", "restart", "nginx")
	fmt.Print(cmd)
	if cmd.Err != nil {
		log.Fatal(cmd.Err)
	}
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	CreateNginxConfig([]Site{
		{Domain: "example.org", Url: "/", LocalPath: "/dist/"},
		{Domain: "www.example.org", Url: "/naher/", LocalPath: "/dist/"},
	})
	RestartNginx()
}
