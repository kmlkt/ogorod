package main

import (
	"fmt"
	"os"

	"github.com/kmlkt/ogorod"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {
	fmt.Println(os.Getwd())
	p := ogorod.Postgres{}
	p.Version = "18.0"
	err := p.Install()
	return err
}
