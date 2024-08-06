package main

import (
	"os"

	"github.com/patricka3125/pickle-bot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
