package main

import (
	"fmt"

	"github.com/go-hep/fwk"
)

func main() {
	fmt.Printf("::: fwk-app...\n")
	app := fwk.NewApp()
	err := app.Run()
	if err != nil {
		panic(err)
	}

	fmt.Printf("::: fwk-app... [done]\n")
}
