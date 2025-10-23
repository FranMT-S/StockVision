package cmd

import (
	"os"
)

type Cmd struct {
}

func NewCmd() Cmd {
	return Cmd{}
}

// Execute runs the command
// fill-db: fills the database with initial data
func (c Cmd) Execute() error {
	if len(os.Args) > 1 {
		err := fillDbCmd.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
