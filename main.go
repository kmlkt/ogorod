package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
)

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
`, site.URL, site.LocalPath)
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

type Site struct {
	Repository string
	LocalPath  string
	Domain     string
	URL        string
}

type SitesConfig struct {
	Sites []Site
}

func GitClone(url string) {
	destination := "./repository"
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		log.Fatalf("Error cloning repository: %v", err)
	}
	fmt.Println("Repository cloned successfully!")
}

func CloneSiteFromSettings(config SitesConfig) {
	for _, site := range config.Sites {
		GitClone(site.URL)
	}
}

func AddNewSites(config SitesConfig) {
	var repository string
	var localpath string
	var domain string
	var url string
	fmt.Print("Введите репозиторий: ")
	fmt.Scan(&repository)
	fmt.Print("Введите путь хранения файлов: ")
	fmt.Scan(&localpath)
	fmt.Print("Введите домен: ")
	fmt.Scan(&domain)
	fmt.Print("Введите URL: ")
	fmt.Scan(&url)
	config.Sites = append(config.Sites, Site{repository, localpath, domain, url})
	updatedYAML, err := yaml.Marshal(config)
	if err != nil {
		log.Fatalf("Ошибка при кодировании YAML: %v", err)
	}
	err = os.WriteFile("settings.yaml", updatedYAML, 0666)
	if err != nil {
		log.Fatalf("Ошибка при записи файла: %v", err)
	}
}

func main() {
	yamlFile, err := os.ReadFile("./settings.yaml")
	if err != nil {
		log.Fatalf("Ошибка при чтении файла: %v", err)
	}
	var config SitesConfig
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatalf("Ошибка при декодировании YAML: %v", err)
	}
}
