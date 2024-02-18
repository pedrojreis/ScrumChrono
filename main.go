package main

import (
	"github.com/pedrojreis/ScrumChrono/cmd"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		panic(err)
	}
}
