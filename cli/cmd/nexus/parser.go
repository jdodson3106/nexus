package nexus

import (
	"fmt"
	"os"
)

const (
	DB_MODEL = "dbModel.nexus"
	HANDLER  = "handler.nexus"
	MODEL    = "model.nexus"
	VIEW     = "view.nexus"
)

func GenerateFromTemplate(templatType string, app *App) error {

	switch templatType {
	case DB_MODEL:
		if err := genDbModel(app); err != nil {
			return err
		}
		return nil
	case HANDLER:
		if err := genHandler(); err != nil {
			return err
		}
		return nil
	case MODEL:
		if err := genModel(); err != nil {
			return err
		}
		return nil
	case VIEW:
		if err := genView(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("invalid template type provided: %s", templatType)
	}
}

func genDbModel(app *App) error {
	template, err := getFileForType(DB_MODEL)
	if err != nil {
		return err
	}

	for _, model := range app.Models {
		if err := parseTemplate(template, model); err != nil {
			return err
		}
	}
	return nil
}

func genModel() error {
	template, err := getFileForType(MODEL)
	if err != nil {
		return err
	}

	if err := parseTemplate(template); err != nil {
		return err
	}
	return nil
}

func genHandler() error {
	template, err := getFileForType(HANDLER)
	if err != nil {
		return err
	}

	if err := parseTemplate(template); err != nil {
		return err
	}
	return nil
}

func genView() error {
	template, err := getFileForType(VIEW)
	if err != nil {
		return err
	}

	if err := parseTemplate(template); err != nil {
		return err
	}
	return nil
}

func parseTemplate(templContents []byte, args ...interface{}) error {
	// ACTIONS:
	// 1. check for interfaces to determine the types required and how to inject into the templates
	// 2. get the new app dir
	// 3. parse the template tokens and inject the appropriate information into a new file
	// 4. save the new file into the relevant directory location
	return nil
}

func getFileForType(tType string) ([]byte, error) {
	template, err := os.ReadFile(fmt.Sprintf("../../templates/%s", tType))
	if err != nil {
		return nil, err
	}
	return template, nil
}
