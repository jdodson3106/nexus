package nexus

import (
	"github.com/spf13/cobra"
)

var newCommand = &cobra.Command{
	Use:     "compile",
	Aliases: []string{"comp"},
	Short:   "executes compilation on your .templ files",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		// generate the entire application layout...
	},
}

func init() {
	rootCmd.AddCommand(newCommand)
}
