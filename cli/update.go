package cli

import (
	"errors"
	"flag"
	"fmt"
	"strconv"

	"github.com/mattdood/spike/run"
)

// Command types, each is required to have a FlagSet
type UpdateCommand struct {
	fs     *flag.FlagSet
	id     string
	key    string
	value  string
	status string
}

func NewUpdateCommand() *UpdateCommand {
	uc := &UpdateCommand{
		fs: flag.NewFlagSet("update", flag.ContinueOnError),
	}
	uc.fs.StringVar(&uc.id, "id", "", "ID of the task")
	uc.fs.StringVar(&uc.key, "key", "", "Key value to update, can be one of 'name', 'description'")
	uc.fs.StringVar(&uc.value, "value", "", "Value of corresponding key")
	uc.fs.StringVar(&uc.status, "status", "", "Status update to handle, can be one of 'C' or 'O', to close or open a task")

	return uc
}

// Update has mutually exclusive arguments:
// key, value go together
// status goes alone
func (uc *UpdateCommand) ParseFlags(args []string) error {
	err := uc.fs.Parse(args)

	if len(uc.id) == 0 && err != flag.ErrHelp {
		return errors.New("length of -id flag must be >0 characters")
	}

	if len(uc.status) == 0 {
		if len(uc.key) == 0 || len(uc.value) == 0 && err != flag.ErrHelp {
			return errors.New("length of both -key and -value flag must be >0 characters when not using -status")
		}
	}

	if len(uc.status) != 0 {
		if len(uc.key) != 0 || len(uc.value) != 0 && err != flag.ErrHelp {
			return errors.New("length of both -key and -value flag must be 0 characters if using -status")
		}
	}

	return err
}

func (uc *UpdateCommand) Run() int {
	id, err := strconv.Atoi(uc.id)
	if err != nil {
		fmt.Println("error converting ID from str to int: ", err)
		return 1
	}

	return run.Update(
		id,
		uc.key,
		uc.value,
		uc.status,
	)
}
