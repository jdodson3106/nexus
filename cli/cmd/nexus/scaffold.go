package nexus

import (
    "bufio"
    "os"
    "fmt"
    "strings"
	"path/filepath"
)

func ScaffoldNewApplication(appName string) {
    absolutePath := getExecutablePath()
    buildAppsDir(absolutePath)
    appDir := createNewAppDir(appName)
    _, err := GenerateNewPropsFile(appDir, appName)
    if err != nil {
        panic(err)
    }
     
}

func getExecutablePath() string {
    // get the executable with the entire path
    loc, err := os.Executable()    
    if err != nil {
        panic(fmt.Errorf("error loading binary execution location :: %v", err))
    }

    // split off the actual nexus executable from the path to get the base path
    exBase, _ := filepath.Split(loc) 
    return exBase
}

func buildAppsDir(absPath string) {
    PrintInfo(fmt.Sprintln("checking for nexus apps directory..."))
    appsDir := filepath.Join(absPath, "apps")
    _, err := os.Stat(appsDir) 
    if err != nil {
        if os.IsNotExist(err) {
            if err = os.Mkdir(appsDir, os.ModePerm); err != nil {
                panic(err)
            } else {
                PrintCreate(fmt.Sprintf("created nexus apps directory at %s\n", appsDir))
            }
        } else {
           panic(err) 
        }
    }
}

func createNewAppDir(appName string) string {
    appsPath := filepath.Join(getExecutablePath(), "apps")
    cleanedAppName := ToSnakeLower(appName)
    appDir := filepath.Join(appsPath, cleanedAppName)
    _, err := os.Stat(appDir)

    if err != nil {
        // that dir doesn't exist. Let's create one
        if os.IsNotExist(err) {
            if mkErr := os.Mkdir(appDir, os.ModePerm); mkErr != nil {
                panic(mkErr) // something terrible happened...
            }
            PrintCreate(fmt.Sprintf("Created new internal app directory at %s\n", appDir))
            return appDir
        } else {
            // unexpected error...
            panic(err)
        }
    }

    // app directory already exists
    kontiue := handleAppExists(appName)

    if !kontiue {
        PrintWarningInfo("Cancelled App Creation...\n")
        os.Exit(1)
    }

    return appDir
}

func handleAppExists(appName string) bool {
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
        return true
    }
    return false
}
