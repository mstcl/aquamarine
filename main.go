package main

import (
	"log"

	"github.com/mstcl/aquamarine/internal/cli"
)

func main() {
	if err := cli.Parse(); err != nil {
		log.Fatal(err)
	}
}
