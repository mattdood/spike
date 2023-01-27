package cli

import (
	"fmt"
)

const usage = `
spike [OPTION] [args...]

Options:
    -h, -help, --help display this help

Executes various commands to manage personal tasks in a
'~/spike/tasks.json' file. Optionally handles Git remote
hosting for tasks.

Example:
$ spike create -name "Create new project" -desc "Longer description" -status "O"`

// Name is always the first arg, use to discover
// command to run. Flags are the rest
type CommandArgs struct {
	name string
	args []string
}

// Pass args to parser then run the
// appropriate command, return exit call
// of the given command
func Run(args []string) int {
	// args are:
	// [/tmp/go-build-dir <cmd flag>]
	if len(args) < 1 {
		fmt.Println("Unknown command. Use -h, -help, or --help to display help")
		return 1
	}

	if len(args) > 1 && len(args) < 3 {
		for _, arg := range args {
			if arg == "-h" || arg == "-help" || arg == "--help" {
				fmt.Println(usage)
				return 1
			}
		}
	}

	command := CommandArgs{
		name: args[1],
		args: args[2:],
	}

	return ParseAndRun(command)
}
