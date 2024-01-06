package nexus

import (
	"fmt"
	"os"
	"reflect"
)

// templates
const (
	DB_MODEL = "dbModel.nexus"
	HANDLER  = "handler.nexus"
	MODEL    = "model.nexus"
	VIEW     = "view.nexus"
)

// tokens
const (
	leader = '#'
	open   = '{'
	close  = '}'
)

var keywords = []string{"app_name", "name", "obj_name", "fields"}

func GenerateFromTemplate(templatType string, app *App) error {

	switch templatType {
	case DB_MODEL:
		if err := genDbModel(app); err != nil {
			return err
		}
		return nil
	case HANDLER:
		if err := genHandler(app); err != nil {
			return err
		}
		return nil
	case MODEL:
		if err := genModel(app); err != nil {
			return err
		}
		return nil
	case VIEW:
		if err := genView(app); err != nil {
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
		if err := parseTemplate(&template, app, model); err != nil {
			return err
		}
	}
	return nil
}

func genModel(app *App) error {
	template, err := getFileForType(MODEL)
	if err != nil {
		return err
	}

	if err := parseTemplate(&template, app); err != nil {
		return err
	}
	return nil
}

func genHandler(app *App) error {
	template, err := getFileForType(HANDLER)
	if err != nil {
		return err
	}

	if err := parseTemplate(&template, app); err != nil {
		return err
	}
	return nil
}

func genView(app *App) error {
	template, err := getFileForType(VIEW)
	if err != nil {
		return err
	}

	if err := parseTemplate(&template, app); err != nil {
		return err
	}
	return nil
}

func getFileForType(tType string) ([]byte, error) {
	template, err := os.ReadFile(fmt.Sprintf("./templates/%s", tType))
	if err != nil {
		return nil, err
	}
	return template, nil
}

func parseTemplate(templContents *[]byte, app *App, args ...interface{}) error {
	// ACTIONS:
	// 1. check for interfaces to determine the types required and how to inject into the templates
	for _, arg := range args {
		typeName := reflect.TypeOf(arg).Name()
		switch typeName {
		case "Model":
			field := reflect.ValueOf(arg).FieldByName("Name").String()
			PrintInfo(fmt.Sprintf("creating model template for %s...\n", field))
			// handle the db data injection
			file, err := generateDbModelFile(templContents, app, field)
			if err != nil {
				PrintWarningInfo("Error generating Model file at location\n")
				PrintWarningInfo(fmt.Sprintf("%s\n", file))
				return err
			}
			continue
		default:
			PrintWarningInfo(fmt.Sprintf("Unknown type %s provided. Ignoring\n", typeName))
			continue
		}
	}
	// 2. get the new app dir
	// 3. parse the template tokens and inject the appropriate information into a new file
	// 4. save the new file into the relevant directory location
	return nil
}

func generateDbModelFile(template *[]byte, app *App, modelName string) (string, error) {
	var data []byte
	line := make([]byte, 0)

	for i := 0; i < len(*template); i++ {
		b := (*template)[i]
		if b == '/' && (*template)[i+1] == '/' {
			// skip ahead to next line
			tmp := (*template)[i:]
			for j, c := range tmp {
				if c == '\n' {
					i = j + 1
					b = (*template)[i]
					break
				}
			}
			continue
		}

		if b == '\n' {
			processed := processLine(&line, app, modelName)
			appendLineToFileData(&data, processed)
			line = make([]byte, 0)
		} else {
			line = append(line, b)
		}

	}

	fileName := fmt.Sprintf("%s/%s.go", app.Directory.Models, modelName)
	printFile(data)
	if err := os.WriteFile(fileName, data, os.ModePerm); err != nil {
		return "", err
	}
	return fileName, nil
}

func printFile(data []byte) {
	var line []byte
	for _, d := range data {
		if d == '\n' {
			fmt.Printf("%s\n", string(line))
			line = nil
		}
		line = append(line, d)
	}
}

func processLine(line *[]byte, app *App, args ...string) []byte {
	// first make sure it's not just a blank line
	if len(*line) == 1 && (*line)[0] == '\n' {
		return *line
	}

	processed := make([]byte, len((*line)))
	for i := 0; i < len(*line); i++ {
		b := (*line)[i]
		if b == leader {
			// verfiy there isn't a wild leader token out there
			if (*line)[i+1] != open {
				msg := fmt.Errorf("Invalid token found after leader %v\n", (*line)[i+1])
				panic(msg)
			}

			// parse the keyword from start to finish
			var kw []byte
			start := i + 2 // jump over the leader and open tokens
			curr := (*line)[start]
			for curr != close {
				kw = append(kw, curr)
				start++
				curr = (*line)[start]
			}

			// look up the keyword and replace it with the data
			skw := string(kw)
			val, ok := getValueForKeyword(skw, app)
			if !ok {
				msg := fmt.Errorf("Unknown keyword #{%s} found in template\n", skw)
				panic(msg)

			}
			// append the keyword value to the line
			for _, v := range val {
				processed = append(processed, v)
			}

			// move the cursor up to after the close token
			i = start
			if i <= len(*line)-1 {
				b = (*line)[i]
			} else {
				b = (*line)[len(*line)-1]
			}
			continue
		}
		processed = append(processed, b)
	}
	processed = append(processed, '\n')
	return processed
}

func getValueForKeyword(keyword string, app *App, args ...string) ([]byte, bool) {
	switch keyword {
	case "app_name":
		return []byte(app.Name), true
	case "name":
		return []byte(args[0]), true
	case "obj_name":
		return []byte(args[0]), true
	case "fields":
		return []byte("fiels"), true
	default:
		return nil, false
	}
}

func appendLineToFileData(data *[]byte, processed []byte) {
	for _, d := range processed {
		*data = append(*data, d)
	}
}
