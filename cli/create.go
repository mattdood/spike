package cli

import (
	"errors"
	"flag"

	"github.com/mattdood/spike/run"
)

// Command types, each is required to have a FlagSet
type CreateCommand struct {
	fs          *flag.FlagSet
	name        string
	description string
}

func NewCreateCommand() *CreateCommand {
	cc := &CreateCommand{
		fs: flag.NewFlagSet("create", flag.ContinueOnError),
	}
	cc.fs.StringVar(&cc.name, "name", "", "Name for the task")
	cc.fs.StringVar(&cc.description, "desc", "", "Description of the task")

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

	return err
}

func (cc *CreateCommand) Run() int {
	return run.Create(
		cc.name,
		cc.description,
	)
}
