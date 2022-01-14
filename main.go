package main

import (
	"fmt"

	"github.com/alexey-sderzhikov/regent/cli"
)

func main() {

	err := cli.Start()
	if err != nil {
		fmt.Print(err)
	}
}
