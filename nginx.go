package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
)

func (c Config) ConfigureNginx() {
	if runtime.GOOS != "linux" {
		log.Println("NGINX can be restarted automatically only on linux")
		return
	}
	servers := map[string][]Site{}
	for _, site := range c.Sites {
		servers[site.Domain] = append(servers[site.Domain], site)
	}
	file, err := os.OpenFile(sitesEnabled, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	StupidHandle(err)
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
`, site.URL, LocalPath(site))
		}
		fmt.Fprint(file, "}")
	}
}

func RestartNginx() {
	if runtime.GOOS != "linux" {
		log.Println("NGINX can be restarted automatically only on linux")
		return
	}
	cmd := exec.Command("systemctl", "restart", "nginx")
	StupidHandle(cmd.Err)
	err := cmd.Run()
	StupidHandle(err)
}
