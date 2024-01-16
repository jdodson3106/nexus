package nexus

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type NewDirectory struct {
	Root          string
	Views         string
	Models        string
	Handlers      string
	Configuration string
}

type ModelCol struct {
	Title    string
	DataType string
}

type Model struct {
	Name string
	Cols []ModelCol
}

type App struct {
	Name      string
	Directory NewDirectory
	Models    []Model
}

func ScaffoldNewApplication(appName string) {
	app := App{Name: appName}

	absolutePath := GetExecutablePath()
	buildAppsDir(absolutePath)
	appDir := createNewAppDir(appName)
	propsFilePath, err := GenerateNewPropsFile(appDir, appName)
	if err != nil {
		panic(err)
	}

	newDirectoryStructure := scaffoldApplicationStructure(appName)
	app.Directory = newDirectoryStructure
	SetProperty("APP_ROOT", newDirectoryStructure.Root, propsFilePath)
	SetProperty("VIEWS_ROOT", newDirectoryStructure.Views, propsFilePath)
	SetProperty("HANDLERS_ROOT", newDirectoryStructure.Handlers, propsFilePath)
	SetProperty("MODELS_ROOT", newDirectoryStructure.Models, propsFilePath)
	SetProperty("CONF_ROOT", newDirectoryStructure.Configuration, propsFilePath)

	generateBaseFiles(&app)
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

	newDirStruct := NewDirectory{
		Root:          newDir,
		Views:         filepath.Join(newDir, "views"),
		Models:        filepath.Join(newDir, "models"),
		Handlers:      filepath.Join(newDir, "handlers"),
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

func generateBaseFiles(app *App) {
	PrintInfo("generating initial application files...\n")

	// 1: Setup Models
	setupModels(app)
	//		- Use Defaults? (SQLite)
	// 2: Create initial views and routes
	// 3: Create handlers to work between views and models

}

func setupModels(app *App) {
	// TODO: If the user is using the default setting - auto install gorm.io
	useDefault := bufio.NewReader(os.Stdin)

	PrintWarning("Setup database model? (y/n): ")
	ans, err := useDefault.ReadString('\n')
	if err != nil {
		panic(err)
	}

	if strings.Trim(strings.ToLower(ans), "\n") == "y" {
		PrintWarning("Initial model name(s): ")
		modelsReader := bufio.NewReader(os.Stdin)
		modelsAns, err := modelsReader.ReadString('\n')
		if err != nil {
			panic(err)
		}

		modelsAns = strings.Trim(modelsAns, "\n")
		userModels := strings.Split(modelsAns, ",")
		models := make([]Model, len(userModels))

		for i, model := range userModels {
			m := Model{Name: strings.TrimSpace(model)}

			// get and parse the column attributes
			PrintWarningInfo(fmt.Sprintf("columns for %s: \n", model))
			usrCols, err := modelsReader.ReadString('\n')
			if err != nil {
				// todo: handle this error
				panic(err)
			}

			cols := strings.Split(usrCols, ",")
			for _, col := range cols {
				data := strings.Split(col, ":")
				if len(data) != 2 {
					PrintWarning(fmt.Sprintf("Invalid Entry: %s\n", col))
					break
				}
				name, attrs := data[0], data[1]
				m.Cols = append(m.Cols, ModelCol{Title: name, DataType: attrs})
			}

			// add the created model to the model list for the app
			models[i] = m
		}

		// add the models to the app object
		app.Models = models

		// geneate the db model files from templates
		if genErr := GenerateFromTemplate(DB_MODEL, app); genErr != nil {
			// todo: handle this error mo' betta
			panic(genErr)
		}
	}

}
