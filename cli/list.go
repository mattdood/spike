package cli

import (
	"errors"
	"flag"

	"github.com/mattdood/spike/run"
)

// Command types, each is required to have a FlagSet
type ListCommand struct {
	fs     *flag.FlagSet
	status string
}

func NewListCommand() *ListCommand {
	lc := &ListCommand{
		fs: flag.NewFlagSet("list", flag.ContinueOnError),
	}
	lc.fs.StringVar(&lc.status, "status", "", "Status to display list of tasks")

	return lc
}

func (lc *ListCommand) ParseFlags(args []string) error {
	err := lc.fs.Parse(args)

	if len(lc.status) == 0 && err != flag.ErrHelp {
		return errors.New("length of -status flag must be >0 characters")
	}

	return err
}

func (lc *ListCommand) Run() int {
	return run.List(
		lc.status,
	)
}
