package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/siaal/gator/internal/cli"
)

func main() {
	state, err := cli.NewState()
	if err != nil {
		slog.Error("Failed to initialise cli", "err", err)
		os.Exit(1)
	}
	cmds := cli.DefaultCommands()
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Printf("requires cmd arg\n")
		os.Exit(1)
	}
	cmd := cli.Command{}
	cmd.Name = args[0]
	cmd.Args = args[1:]
	if err = cmds.Run(&state, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
