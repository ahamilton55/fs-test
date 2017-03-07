package main

import (
	"log"
	"os"

	"github.com/ahamilton55/fs-test/pad/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Printf("Error executing deployment: %s", err.Error())
		os.Exit(1)
	}
}
