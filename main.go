package main

import (
	"fmt"

	"github.com/alexey-sderzhikov/regent/cli"
	"github.com/alexey-sderzhikov/regent/restapi"
)

func main() {
	is, err := restapi.GetIssues("86")
	if err != nil {
		fmt.Print(err)
	} else {
		fmt.Print(is)
	}
	cli.Start()
}
