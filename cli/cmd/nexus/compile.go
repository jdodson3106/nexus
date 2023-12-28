package nexus

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:     "compile",
	Aliases: []string{"comp"},
	Short:   "executes compilation on your .templ files",
	Args:    cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("templ", "generate")
		out, err := c.Output()
		if err != nil {
			fmt.Print(err)
		}
		fmt.Printf("Executing...\n%s\n", out)
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
