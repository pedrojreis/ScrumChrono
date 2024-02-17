package main

import (
	"ScrumChrono/cmd"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		panic(err)
	}
}
