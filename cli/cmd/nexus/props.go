package nexus

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)


const (
    PROP_FILE_EXTENSION = "properties"
    APP_NAME = "APP_NAME"
)

func AppExistsError(appName string) error {
    return errors.New(fmt.Sprintf("a nexus application called \"%s\" already exists.", appName))
}


func ReadProperty(propName, propFilePath string) (string, error) {
    return "", nil
}

func SetProperty(key, val, propFilePath string) bool {
    return true
}

// GenerateNewPropsFile checks if a property file exists for the given
// appName. If it does, && !overwrite then it will throw error
//
// If the file creation succeeds, it will write the app name in and return 
// the property file path.
func GenerateNewPropsFile(appName string, overwriteExisting bool) (string, error) {
    // get the executable with the entire path
    loc, err := os.Executable()    

    // split off the actual nexus executable from the path to get the base path
    exBase, _ := filepath.Split(loc) 
    if err != nil {
        return "", fmt.Errorf("error loading binary execution location :: %v", err)
    }

    // setupt to check and create a props directory
    propsDir := filepath.Join(exBase, "props")

    created, err := createNewDir(propsDir)
    if err != nil {
        return "", fmt.Errorf("error creating properties director at %s :: %v", propsDir, err)
    }

    if created {
        fileName := fmt.Sprintf("%s/%s.%s", propsDir, appName, PROP_FILE_EXTENSION)
        exists := propsFileExists(fileName)
        if !exists || (exists && overwriteExisting) {
            propsFileCreated, err := createNewPropsFile(fileName, appName)
            if err != nil {
                return "", err
            }

            if propsFileCreated {
                return propsDir, nil
            } else {
                return "", fmt.Errorf("unexpected error checking stat for properties file at %s", fileName)
            }
        } 

        // the file already exists (a la an app already exists with that name) so return the err
        return "", AppExistsError(appName)
    }
    
    // shouldn't make it here... but nevertheless, we return a generic err 
    return "", fmt.Errorf("unexpected error creating properties director at %s", propsDir)
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
    _, err = f.Write([]byte(fmt.Sprintf("%s=%s\n", APP_NAME, appName)))
    if err != nil {
        return false, err
    }
    PrintCreate(fmt.Sprintf("Created property file for app %s\n", appName))
    return true, nil
}

func createNewDir(path string) (bool, error) {
    PrintInfo(fmt.Sprintln("checking for nexus directory..."))
    _, err := os.Stat(path) 
    if err != nil {
        if os.IsNotExist(err) {
            if err = os.Mkdir(path, os.ModePerm); err != nil {
                return true, err
            } else {
                PrintCreate(fmt.Sprintf("created property directory at %s\n", path))
            }
        } else {
            return true, err
        }
    }
    return true, nil
}
