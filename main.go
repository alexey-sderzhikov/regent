package main

import (
	"log"

	"github.com/alexey-sderzhikov/regent/cli"
)

func main() {
	if err := cli.Start(); err != nil {
		log.Fatal(err)
	}
}
