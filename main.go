package main

import (
	"fmt"
	"github.com/ishustava/migrator/cmd"
	"github.com/jessevdk/go-flags"
	"os"
	"runtime/debug"
)

func main() {
	debug.SetTraceback("all")
	parser := flags.NewParser(&cmd.Migrator, flags.HelpFlag)
	parser.SubcommandsOptional = true
	parser.CommandHandler = func(command flags.Commander, args []string) error {
		if command == nil {
			parser.WriteHelp(os.Stderr)
			os.Exit(1)
		}

		return command.Execute(args)
	}

	_, err := parser.Parse()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
