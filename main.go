package main

import (
	"os"

	"github.com/cmrigney/skill-share/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
