package nexus

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ToSnakeLower(appName string) string {
	lowerName := strings.ToLower(appName)
	return strings.Join(strings.Split(lowerName, " "), "_")
}

func ToSnakeUpper(appName string) string {
	upperName := strings.ToUpper(appName)
	return strings.Join(strings.Split(upperName, " "), "_")
}

func GetExecutablePath() string {
	// get the executable with the entire path
	loc, err := os.Executable()
	if err != nil {
		panic(fmt.Errorf("error loading binary execution location :: %v", err))
	}

	// split off the actual nexus executable from the path to get the base path
	exBase, _ := filepath.Split(loc)
	return exBase
}
