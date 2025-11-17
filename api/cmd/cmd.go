package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "myapp",
	Short: "Cli for StockVision",
}

type Cmd struct {
}

func NewCmd() Cmd {
	return Cmd{}
}

var jsonPath string

func init() {
	rootCmd.AddCommand(fillDbCmd)
	fillDbCmd.Flags().StringVar(&jsonPath, "json", "", "Path to the JSON file (optional)")

}

// Execute runs the command
// fill-db: fills the database with initial data
func (c Cmd) Execute() error {
	if len(os.Args) > 1 {
		err := rootCmd.Execute()
		if err != nil {
			return err
		}
	}

	return nil
}
