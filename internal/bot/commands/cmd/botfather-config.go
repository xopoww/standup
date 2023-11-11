//nolint:forbidigo // it's ok to use fmt.Print* in this utility
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xopoww/standup/internal/bot/commands"
)

func main() {
	cmds, err := commands.LoadDescriptions()
	if err != nil {
		log.Fatalf("Load descriptions: %s.", err)
	}
	fmt.Fprintf(os.Stderr, "Use the following in /setcommands for your bot:\n\n")
	if len(cmds) == 0 {
		fmt.Println("/empty")
	}
	for _, cmd := range cmds {
		fmt.Printf("%s - %s\n", cmd.Name, cmd.Short)
	}
}
