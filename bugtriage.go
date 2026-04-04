package main

import (
	"log"
	"os"
)

func bugTriage(inputFile string) {
	file, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}
	os.Stdout.Write(file)
}
