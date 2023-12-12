package nexus

import (
	"fmt"

	"github.com/spf13/cobra"
)

const helpText = `
+--------------------------------------------------------------------------+
|       ________    _______       ___    ___  ___  ___   ________          |
|      |\   ___  \ |\  ___ \     |\  \  /  /||\  \|\  \ |\   ____\         |
|      \ \  \\ \  \\ \   __/|    \ \  \/  / /\ \  \\\  \\ \  \___|_        |
|       \ \  \\ \  \\ \  \_|/__   \ \    / /  \ \  \\\  \\ \_____  \       |
|        \ \  \\ \  \\ \  \_|\ \   /     \/    \ \  \\\  \\|____|\  \      |
|         \ \__\\ \__\\ \_______\ /  /\   \     \ \_______\ ____\_\  \     |
|          \|__| \|__| \|_______|/__/ /\ __\     \|_______||\_________\    |
|                                |__|/ \|__|               \|_________|    |
|                                                                          |
|                     An opinionated Web Framework in Go                   |
+--------------------------------------------------------------------------+
`

var rootCmd = &cobra.Command{
	Use:   "nexus",
	Short: "nexus - an opinionated web framework written in Go",
	Long:  helpText,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			fmt.Printf(helpText)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
