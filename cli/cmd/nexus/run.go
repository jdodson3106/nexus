package nexus

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var executePath string
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Executes the Nexus server",
	Run: func(cmd *cobra.Command, args []string) {
		f := fmt.Sprintf("%smain.go", executePath)
        if executePath != "" {
            if e := os.Setenv("NEXUS_APP_EXECUTION_PATH", getAppFolderName()); e != nil {
                panic(e)
            }
        }
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

func getAppFolderName() string {
    parts := strings.Split(executePath, "/")
    return parts[len(parts)-2]
} 
