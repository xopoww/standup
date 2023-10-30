package main

import (
	"fmt"
	"os"

	"github.com/xopoww/standup/internal/standupctl/commands"
)

func main() {
	cmd := commands.Root()
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s.\n", err)
		os.Exit(1)
	}
}
