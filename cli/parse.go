package cli

import (
	"flag"
	"fmt"
	"strings"
)

// Custom array implementation of flag variable
type arrayFlag []string

func (fs *arrayFlag) String() string {
	return strings.Join(*fs, " ")
}

// Custom Usage() override for flags with no arguments
func CommandUsage(fs *flag.FlagSet, name string) {
	// If more custom usage is needed we can parse
	// the command name to preset messages
	fs.Usage = func() {
		fmt.Printf(
			"Wraps `git %s`, takes no arguments. Only operates on `cook` directory.\n",
			name,
		)
	}
}

// Accepts space separated list of values
func (fs *arrayFlag) Set(value string) error {
	for _, file := range strings.Split(value, " ") {
		*fs = append(*fs, file)
	}

	return nil
}

// Runner interface that passes all
// command functions
type Runner interface {
	ParseFlags([]string) error
	Run() int
}

func ParseAndRun(command CommandArgs) int {
	// Register commands
	cmds := map[string]Runner{
		"create": NewCreateCommand(),
		"update": NewUpdateCommand(),
		"list":   NewListCommand(),

		"commit": NewCommitCommand(),
		"add":    NewAddCommand(),
		"init":   NewInitCommand(),
		"pull":   NewPullCommand(),
		"push":   NewPushCommand(),
	}

	// Determine cmd that was passed, init,
	// then run
	cmd := cmds[command.name]
	err := cmd.ParseFlags(command.args)

	switch {
	// Usage information for flags is enabled by default
	// if we pass on the `flag.ErrHelp` during arg parsing
	case err == flag.ErrHelp:
		return 0
	case err != nil:
		fmt.Println(err.Error())
		return 1
	}

	return cmd.Run()
}
