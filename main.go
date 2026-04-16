package main

import (
	"log"
	"os"
)

func main() {
	if err := forge(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
