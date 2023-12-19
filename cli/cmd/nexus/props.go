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
    loc, err := os.Executable()    
    if err != nil {
        return "", fmt.Errorf("Error loading binary execution location :: %v", err)
    }
    
    // _, err = os.Open(fmt.Sprintf("%s/%s.properties", loc, appName))
    // if err == nil && !overwriteExisting {
    //     return "", fmt.Errorf("App \"%s\" already exists.", appName)
    // }
    _, err = os.Stat(loc) 
    if err != nil {
        return "", err
    }

    err = os.Mkdir("./props", 0755)
    if err != nil {
        fmt.Println("Failed making dir")
        fmt.Println(err.Error())
    }

    fileName := fmt.Sprintf("~%s/%s.%s", loc, appName, PROP_FILE_EXTENSION) 
    f, err := os.Create(fileName)

    if err != nil {
        fmt.Println(err.Error())
        return "", AppExistsError(appName)
    }

    defer f.Close()

    _, err = f.Write([]byte(fmt.Sprintf("%s=%s\n", APP_NAME, appName)))
    if err != nil {
        return "", err
    }
    
    return fileName, nil
}
