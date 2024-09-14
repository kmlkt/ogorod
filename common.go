package main

import (
	"log"
	"os"
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
