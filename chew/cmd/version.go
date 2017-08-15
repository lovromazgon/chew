package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"bitbucket.org/lovromazgon/chew"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Chew",
	Long:  `All software has versions. This is Chew's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Chew %s (%s)\n", chew.Version, chew.VersionDate)
	},
}
