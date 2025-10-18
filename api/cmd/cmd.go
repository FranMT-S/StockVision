package cmd

import (
	"os"
)

type Cmd struct {
}

func NewCmd() Cmd {
	return Cmd{}
}

func (c Cmd) Execute() error {
	if len(os.Args) > 1 {
		fillDbCmd.Execute()
	}

	return nil
}
