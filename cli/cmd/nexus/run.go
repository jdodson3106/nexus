package nexus

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var executePath string
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Executes the Nexus server",
	Run: func(cmd *cobra.Command, args []string) {
		f := fmt.Sprintf("%smain.go", executePath)
		c := exec.Command("go", "run", f)
		stdout, err := c.StdoutPipe()
		err = c.Start()
		if err != nil {
			panic(err)
		}
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)
		}
		c.Wait()
	},
}

func init() {
	runCmd.Flags().StringVarP(&executePath, "path", "p", "./", "Provide a path to your main.go")
	rootCmd.AddCommand(runCmd)
}
