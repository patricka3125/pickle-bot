package main

import (
	"os"

	"code.byted.org/patrick.liao/pickle-bot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
