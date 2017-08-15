package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
)

func init() {
	RootCmd.AddCommand(functionsCmd)
}

var functionsCmd = &cobra.Command{
	Use:   "functions",
	Short: "Print the documentation for custom added functions",
	Long:  `Chew has many custom functions that can be used inside templates. The documentation for each function can be viewed with this command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Functions")
		// TODO
	},
}
