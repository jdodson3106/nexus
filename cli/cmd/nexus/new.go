package nexus

import (

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Creates a new Nexus Application",
	Run: func(cmd *cobra.Command, args []string) {
        var appName string 

        if len(args) < 1 {
            cobra.CompErrorln("Must provide a name for the application")
            return 
        }

        if len(args) > 1 {
            cobra.CompErrorln("Too many arguments provided.")
            return 
        }

        appName = args[0]
        _, err := GenerateNewPropsFile(appName, false)
        if err != nil {
            if err == AppExistsError(appName) {
              // request an overwrite?  
            } else {
                cobra.CompErrorln(err.Error())
                return
            }
        }
        


	},
}

func init() {
	rootCmd.AddCommand(newCmd)
}

