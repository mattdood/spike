package cli

import (
	"errors"
	"flag"
	"fmt"
	"strings"

	"github.com/mattdood/spike/run"
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

// Command types, each is required to have a FlagSet
type CreateCommand struct {
	name    string
	description string
	status string
	fs       *flag.FlagSet
}

func NewCreateCommand() *CreateCommand {
	cc := &CreateCommand{
		fs: flag.NewFlagSet("create", flag.ContinueOnError),
	}
	cc.fs.StringVar(&cc.name, "name", "", "Name for the task")
	cc.fs.StringVar(&cc.description, "desc", "", "Description of the task")
	cc.fs.StringVar(&cc.status, "status", "O", "Status code for the task (O - Open, C - Closed)")

	return cc
}

func (cc *CreateCommand) ParseFlags(args []string) error {
	err := cc.fs.Parse(args)

	if len(cc.name) == 0 && err != flag.ErrHelp {
		return errors.New("length of -name flag must be >0 characters")
	}

	if len(cc.description) == 0 && err != flag.ErrHelp {
		return errors.New("length of -desc flag must be >0 characters")
	}

	if len(cc.status) == 0 && err != flag.ErrHelp {
		return errors.New("length of -status flag must be >0 characters")
	}

	return err
}

func (cc *CreateCommand) Run() int {
	run.Create(
		cc.name,
		cc.description,
		cc.status,
	)

	return 0
}

// Git command wrapper for adding
// files to be tracked
// `git add <arg>`
type AddCommand struct {
	fs    *flag.FlagSet
	files arrayFlag
}

func NewAddCommand() *AddCommand {
	ac := &AddCommand{
		fs: flag.NewFlagSet("add", flag.ContinueOnError),
	}
	ac.fs.Var(&ac.files, "file", "Files to add to git tracking. Space separated list.")

	return ac
}

func (ac *AddCommand) ParseFlags(args []string) error {
	err := ac.fs.Parse(args)

	if len(ac.files) == 0 && err != flag.ErrHelp {
		ac.files = append(ac.files, ".")
	}

	return err
}

func (ac *AddCommand) Run() int {
	run.Add(ac.files)
	return 0
}

// Git command wrapper for repo init
type InitCommand struct {
	fs *flag.FlagSet
}

func NewInitCommand() *InitCommand {
	ic := &InitCommand{
		fs: flag.NewFlagSet("init", flag.ContinueOnError),
	}
	CommandUsage(ic.fs, ic.fs.Name())

	return ic
}

func (ic *InitCommand) ParseFlags(args []string) error {
	err := ic.fs.Parse(args)

	if len(args) > 0 && err != flag.ErrHelp {
		return errors.New("this command takes no arguments")
	}

	return err
}

func (ic *InitCommand) Run() int {
	run.Init()
	return 0
}

// Git command wrapper for repo push
type PushCommand struct {
	fs *flag.FlagSet
}

func NewPushCommand() *PushCommand {
	pc := &PushCommand{
		fs: flag.NewFlagSet("push", flag.ContinueOnError),
	}
	CommandUsage(pc.fs, pc.fs.Name())

	return pc
}

func (pc *PushCommand) ParseFlags(args []string) error {
	err := pc.fs.Parse(args)

	if len(args) >= 0 && err != flag.ErrHelp {
		return errors.New("this command takes no arguments")
	}

	return err
}

func (pc *PushCommand) Run() int {
	run.Push()
	return 0
}

// Git command wrapper for repo pull
type PullCommand struct {
	fs *flag.FlagSet
}

func NewPullCommand() *PullCommand {
	pc := &PullCommand{
		fs: flag.NewFlagSet("pull", flag.ContinueOnError),
	}
	CommandUsage(pc.fs, pc.fs.Name())

	return pc
}

func (pc *PullCommand) ParseFlags(args []string) error {
	err := pc.fs.Parse(args)

	if len(args) >= 0 && err != flag.ErrHelp {
		return errors.New("this command takes no arguments")
	}

	return err
}

func (pc *PullCommand) Run() int {
	run.Pull()
	return 0
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
