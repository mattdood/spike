package cli

import (
	"errors"
	"flag"

	"github.com/mattdood/spike/run"
)

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

// Git command wrapper for commiting
// files to be tracked
// `git commit <arg>`
type CommitCommand struct {
	fs    *flag.FlagSet
	message string
}

func NewCommitCommand() *CommitCommand {
	cc := &CommitCommand{
		fs: flag.NewFlagSet("commit", flag.ContinueOnError),
	}
	cc.fs.StringVar(&cc.message, "m", "", "Message for the git commit.")

	return cc
}

func (cc *CommitCommand) ParseFlags(args []string) error {
	err := cc.fs.Parse(args)
	return err
}

func (cc *CommitCommand) Run() int {
	run.Commit(cc.message)
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
