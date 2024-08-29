package nexus

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
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
	LEADER = '#'
	OPEN   = '{'
	CLOSE  = '}'
)

type GeneratorArgs struct {
	Module    string
	App       *App
	ArgsMap   map[string]any
	Generator GeneratorFunc
}

// var keywords = []string{"app_name", "name", "obj_name", "fields"}

type GeneratorFunc func(*GeneratorArgs)

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
		argMap := map[string]any{
			"model": model,
		}
		args := GeneratorArgs{Module: "models", App: app, ArgsMap: argMap}
		if err := parseTemplate(&template, &args); err != nil {
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

	args := GeneratorArgs{Module: "models", App: app}
	if err := parseTemplate(&template, &args); err != nil {
		return err
	}
	return nil
}

func genHandler(app *App) error {
	template, err := getFileForType(HANDLER)
	if err != nil {
		return err
	}

	args := GeneratorArgs{Module: "handlers", App: app}
	if err := parseTemplate(&template, &args); err != nil {
		return err
	}
	return nil
}

func genView(app *App) error {
	template, err := getFileForType(VIEW)
	if err != nil {
		return err
	}

	args := GeneratorArgs{Module: "views", App: app}
	if err := parseTemplate(&template, &args); err != nil {
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

func parseTemplate(templContents *[]byte, args *GeneratorArgs) error {
	for k, v := range args.ArgsMap {
		switch k {
		case "model":
			modelName := reflect.ValueOf(v).FieldByName("Name").String()
			PrintInfo(fmt.Sprintf("creating model template for %s...\n", modelName))

			// handle the db data injection
			args.ArgsMap["modelName"] = modelName
			file, err := generateDbModelFile(templContents, args)
			if err != nil {
				PrintWarningInfo("Error generating Model file at location\n")
				PrintWarningInfo(fmt.Sprintf("%s\n", file))
				return err
			}
			continue

		default:
			// PrintWarningInfo(fmt.Sprintf("Unknown type %s provided. Ignoring\n", k))
			continue
		}
	}
	return nil
}

func generateDbModelFile(template *[]byte, args *GeneratorArgs) (string, error) {
	data := make([]byte, 0)
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
			processed := processLine(&line, args)
			data = append(data, processed...)
			line = make([]byte, 0)
		} else {
			line = append(line, b)
		}
	}

	t, modelName, found := getValueForMapKey("modelName", args)
	if !found {
		return "", errors.New("no argument found for key 'modelName'")
	}

	if t.Name() != "string" {
		return "", fmt.Errorf("value for key 'modelName' expected a string, got %v", t.Name())
	}

	name := strings.ToLower(reflect.ValueOf(modelName).String())
	fileName := fmt.Sprintf("%s/%s.go", args.App.Directory.Models, name)
	if err := os.WriteFile(fileName, data, os.ModePerm); err != nil {
		return "", err
	}
	return fileName, nil
}

func processLine(line *[]byte, args *GeneratorArgs) []byte {
	// first make sure it's not just a blank line
	if len(*line) == 1 && (*line)[0] == '\n' {
		return *line
	}

	processed := make([]byte, 0)
	for i := 0; i < len(*line); i++ {
		b := (*line)[i]
		if b == LEADER {
			// verfiy there isn't a wild LEADER token out there
			if (*line)[i+1] != OPEN {
				msg := fmt.Errorf("Invalid token found after LEADER %v\n", (*line)[i+1])
				panic(msg)
			}

			// parse the keyword from start to finish
			var kw []byte
			cursor := i + 2 // jump over the LEADER and open tokens
			curr := (*line)[cursor]
			for curr != CLOSE {
				kw = append(kw, curr)
				cursor++
				curr = (*line)[cursor]
			}

			// look up the keyword and replace it with the data
			skw := string(kw)
			val, ok := getValueForKeyword(skw, args) // TODO: handle the passed args
			if !ok {
				msg := fmt.Errorf("Unknown keyword #{%s} found in template\n", skw)
				panic(msg)

			}
			// append the keyword value to the line
			processed = append(processed, val...)

			// move the cursor up to after the close token
			i = cursor
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

func getValueForKeyword(keyword string, args *GeneratorArgs) ([]byte, bool) {
	switch keyword {
	case "app_name":
		if args.Module == "" {
			return nil, false
		}
		return []byte(args.Module), true
	case "name":
		rtn, err := getNameValue(args)
		if err != nil {
			printKeywordError(err)
			return nil, false
		}
		return rtn, true
	case "obj_name":
		rtn, err := getObjNameValue(args)
		if err != nil {
			printKeywordError(err)
			return nil, false
		}
		return rtn, true
	case "fields":
		fields, err := getFieldsValue(args)
		if err != nil {
			printKeywordError(err)
			return nil, false
		}

		rtn := processFields(fields)
		return rtn, true
	default:
		return nil, false
	}
}

func getNameValue(args *GeneratorArgs) ([]byte, error) {
	if name := args.ArgsMap["modelName"]; name != nil {
		val := name.(string)
		return []byte(val), nil
	}
	return nil, errors.New("no model name found.")
}

func getObjNameValue(args *GeneratorArgs) ([]byte, error) {

	return nil, nil
}

func getFieldsValue(args *GeneratorArgs) ([][]string, error) {
	// get the model using reflection
	model := reflect.ValueOf(args.ArgsMap["model"])
	if model.Kind() != reflect.Struct {
		return nil, fmt.Errorf("invalid type %s. expected struct type", model.Kind().String())
	}

	var fields [][]string

	for i := 0; i < model.NumField(); i++ {
		field := model.Field(i)
		if field.Type() == reflect.TypeOf([]ModelCol{}) {
			cols := field.Interface().([]ModelCol)
			for _, col := range cols {
				fields = append(fields, []string{
					col.Title,
					col.DataType,
				})
			}
		}
	}
	return fields, nil
}

func processFields(fieldSlice [][]string) []byte {
	// TODO: Extend this to add gorm tags to the models
	data := make([]byte, 0)

	for _, fieldSet := range fieldSlice {
		if len(fieldSet) != 2 {
			PrintWarningInfo(fmt.Sprintf("invalid fields found :: %v", fieldSet))
			continue
		}
		field := fmt.Sprintf("\t%s %s\n", fieldSet[0], fieldSet[1])
		data = append(data, field...)
	}

	return data
}

func getValueForMapKey(key string, args *GeneratorArgs) (reflect.Type, any, bool) {
	var val any
	for k, v := range args.ArgsMap {
		if k == key {
			val = v
			break
		}
	}

	// key no exists
	if val == nil {
		return nil, nil, false
	}

	t := reflect.TypeOf(val)
	return t, val, true
}

func printKeywordError(err error) {
	PrintWarningInfo(fmt.Sprintf("error getting keyword :: %s", err))
}
