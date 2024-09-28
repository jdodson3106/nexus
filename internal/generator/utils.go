package generator

import (
	"fmt"
	"os"
)

var (
	viewsPath      string
	controllerPath string
	modelsPath     string
	buildDirPath   string
)

const (
	DEFAULT_APP_NAME = "app"
)

func getPathVar() string {
	return fmt.Sprintf("/%s", os.Getenv("NEXUS_APP_EXECUTION_PATH"))
}

func getRelativeAppPath() string {
	dir, err := os.Getwd()
	if err != nil {
		dir = fmt.Sprintf("./%s", DEFAULT_APP_NAME)
	}
	p := getPathVar()
	dir += p
	return dir
}

func getRelativeViewsPath() string {
	if viewsPath == "" {
		viewsPath = getRelativeAppPath() + "views"
		return viewsPath
	}
	return viewsPath
}

func getRelativeControllerPath() string {
	if controllerPath == "" {
		controllerPath = getRelativeAppPath() + "controllers"
		return controllerPath
	}
	return controllerPath
}

func getRelativeModelsPath() string {
	if modelsPath == "" {
		modelsPath = getRelativeAppPath() + "models"
		return modelsPath
	}
	return modelsPath
}

func generateRoutePlugin() {

}
