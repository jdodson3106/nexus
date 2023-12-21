package nexus

import (
	"errors"
	"fmt"
	"os"
)


const (
    PROP_FILE_EXTENSION = "properties"
    APP_NAME = "APP_NAME"
)


func ReadProperty(propName, propFilePath string) (string, error) {
    return "", nil
}

func SetProperty(key, val, propFilePath string) bool {
    toWrite := fmt.Sprintf("%s=%s\n", key, val)
    f, err := os.OpenFile(propFilePath, os.O_WRONLY|os.O_APPEND, os.ModeAppend.Perm())  
    if err != nil {
        panic(err)
    }
    defer f.Close()

    _, err = f.Write([]byte(toWrite))
    if err != nil {
        panic(err)
    }

    return true
}

// GenerateNewPropsFile checks if a property file exists for the given
// appName. If it does, && !overwrite then it will throw error
//
// If the file creation succeeds, it will write the app name in and return 
// the property file path.
func GenerateNewPropsFile(appDir, appName string) (string, error) {
    fileName := fmt.Sprintf("%s/%s.%s", appDir, ToSnakeLower(appName), PROP_FILE_EXTENSION)
    propsFileCreated, err := createNewPropsFile(fileName, appName)
    if err != nil {
        return "", err
    }

    if !propsFileCreated {
        return "", errors.New("unexpected error creating properties file")
    }

    return fileName, nil
}

func propsFileExists(path string) bool {
    PrintInfo(fmt.Sprintln("checking if property file exist..."))
    if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
        return false
    }
    return true
}

func createNewPropsFile(fileName, appName string) (bool, error) {
    f, err := os.Create(fileName)

    defer f.Close()

    // write the app name as the first property in the file
    _, err = f.Write([]byte(fmt.Sprintf("%s=%s\n", APP_NAME, ToSnakeUpper(appName))))
    if err != nil {
        return false, err
    }
    PrintCreate(fmt.Sprintf("created property file for app \"%s\"\n", appName))
    return true, nil
}
