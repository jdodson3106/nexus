package nexus

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

func startAppCreate (appName string) {
    createPropertyFile(appName)
}

func createPropertyFile(appName string) bool {
    _, err := GenerateNewPropsFile(appName, false)
    if err != nil {
        if err.Error() == AppExistsError(appName).Error() {
            reader := bufio.NewReader(os.Stdin)
            PrintWarningInfo(fmt.Sprintf("An app with name \"%s\" already exists. \n", appName))
            PrintNote(fmt.Sprintf("Continuing will overwrite the properties for the existing application and may result in conflicts.\n"))
            PrintWarning("Create anyway? (y/n): ")
            ans, err := reader.ReadString('\n')
            if err != nil {
                panic(err)
            }
            a := strings.Trim(strings.ToLower(ans), "\n")
            if a == "y" || a == "yes" {
                _, err := GenerateNewPropsFile(appName, true)
                if err != nil {
                    cobra.CompErrorln(err.Error())
                    return false
                }
                return true
            }
            // if not a yes response then just exit entire workflow
            PrintWarningInfo("Cancelled App Creation...\n")
            os.Exit(1)
        } else {
            cobra.CompErrorln(err.Error())
            return false
        }
    }
    return true
}
