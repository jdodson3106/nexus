package templates

import (
	"fmt"
	"io/ioutil"
)

const (
	DB_MODEL = "dbModel.nexus"
	HANDLER  = "handler.nexus"
	MODEL    = "model.nexus"
	VIEW     = "view.nexus"
)

func GenerateFromTemplate(templatType string) error {

	switch templatType {
	case DB_MODEL:
		genDbModel()
		return nil
	case HANDLER:
		genHandler()
		return nil
	case MODEL:
		genModel()
		return nil
	case VIEW:
		genView()
		return nil
	default:
		return fmt.Errorf("Invalid Template Type Provided: %s", templatType)
	}
}

func genDbModel() error {
	template, err := ioutil.ReadFile("./dbModel.nexus")
	if err != nil {
		return err
	}

	if err := parseTemplate(template); err != nil {
		return err
	}
	return nil
}

func genModel() error {
	template, err := ioutil.ReadFile("./model.nexus")
	if err != nil {
		return err
	}

	if err := parseTemplate(template); err != nil {
		return err
	}
	return nil
}

func genHandler() error {
	template, err := ioutil.ReadFile("./hanler.nexus")
	if err != nil {
		return err
	}

	if err := parseTemplate(template); err != nil {
		return err
	}
	return nil
}

func genView() error {
	template, err := ioutil.ReadFile("./view.nexus")
	if err != nil {
		return err
	}

	if err := parseTemplate(template); err != nil {
		return err
	}
	return nil
}

func parseTemplate(templContents []byte) error {
	return nil
}
