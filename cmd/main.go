package main

import (
	"fmt"
	"os"

	"github.com/kmlkt/ogorod/postgres"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println(os.Getwd())
	p := postgres.Postgres{}
	p.Version = "18.0"
	fmt.Println(p.State())
	err := p.Install()
	return err
}
