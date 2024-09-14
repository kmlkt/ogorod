package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-git/go-git/v5"
	"gopkg.in/yaml.v3"
)

type Site struct {
	Repository string
	LocalPath  string
	Domain     string
	URL        string
}

type SitesConfig struct {
	Sites []Site
}

func gitclone(url string) {
	destination := "./repository"
	_, err := git.PlainClone(destination, false, &git.CloneOptions{
		URL: url,
	})
	if err != nil {
		log.Fatalf("Error cloning repository: %v", err)
	}
	fmt.Println("Repository cloned successfully!")
}

func clonesitefromsettings(config SitesConfig) {
	for _, site := range config.Sites {
		gitclone(site.URL)
	}
}

func addnewsites(config SitesConfig) {
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
