package nexus

import (
    "bufio"
    "os"
    "fmt"
    "strings"
	"path/filepath"
)

type NewDirectory struct {
    Root string 
    Views string 
    Models string
    Handlers  string 
    Configuration string
}

func ScaffoldNewApplication(appName string) {
    absolutePath := GetExecutablePath()
    buildAppsDir(absolutePath)
    appDir := createNewAppDir(appName)
    propsFilePath, err := GenerateNewPropsFile(appDir, appName)
    if err != nil {
        panic(err)
    }
    
    newDirectoryStructure := scaffoldApplicationStructure(appName)
    SetProperty("APP_ROOT", newDirectoryStructure.Root, propsFilePath)
    SetProperty("VIEWS_ROOT", newDirectoryStructure.Views, propsFilePath)
    SetProperty("HANDLERS_ROOT", newDirectoryStructure.Handlers, propsFilePath)
    SetProperty("MODELS_ROOT", newDirectoryStructure.Models, propsFilePath)
    SetProperty("CONF_ROOT", newDirectoryStructure.Configuration, propsFilePath)

    generateBaseFiles(newDirectoryStructure)
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
    appsPath := filepath.Join(GetExecutablePath(), "apps")
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

func scaffoldApplicationStructure(appName string) NewDirectory {
    wd, err := os.Getwd()
    if err != nil {
        panic(err)
    }
    folderName := ToSnakeLower(appName)
    newDir := filepath.Join(wd, folderName)
    
    _, err = os.Stat(newDir)
    if err != nil && os.IsNotExist(err) {
        if err = os.Mkdir(newDir, os.ModePerm); err != nil {
            panic(err)
        }
    } else {
        PrintWarningInfo(fmt.Sprintf("A directory already exists at %s\n", newDir))
        os.Exit(1)
    }

    newDirStruct := NewDirectory {
        Root: newDir,
        Views: filepath.Join(newDir, "views"),
        Models: filepath.Join(newDir, "models"),
        Handlers: filepath.Join(newDir, "handlers"),
        Configuration: filepath.Join(newDir, "conf"),
    }

    if err = os.Mkdir(newDirStruct.Views, os.ModePerm); err != nil {
        panic(err)
    }  

    if err = os.Mkdir(newDirStruct.Models, os.ModePerm); err != nil {
        panic(err)
    }  

    if err = os.Mkdir(newDirStruct.Handlers, os.ModePerm); err != nil {
        panic(err)
    }  

    if err = os.Mkdir(newDirStruct.Configuration, os.ModePerm); err != nil {
        panic(err)
    }  
    
    PrintCreate(fmt.Sprintf("created new application directory %s\n", newDir))
    return newDirStruct
}

func generateBaseFiles(directory NewDirectory) {
    PrintInfo("generating core files...")
}
