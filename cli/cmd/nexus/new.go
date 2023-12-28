package nexus

import (
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new Nexus Application",
	Run: func(cmd *cobra.Command, args []string) {
		var appName string
		if validateArgs(args) {
			appName = args[0]
			startAppCreate(appName)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

func validateArgs(args []string) bool {
	if len(args) < 1 {
		cobra.CompErrorln("Must provide a name for the application")
		return false
	}

	if len(args) > 1 {
		cobra.CompErrorln("Too many arguments provided.")
		return false
	}

	return true
}

func startAppCreate(appName string) {
	ScaffoldNewApplication(appName)
}
