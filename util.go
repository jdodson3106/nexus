package nexus

import (
	"fmt"
	"os/exec"
	"plugin"
	"reflect"

	"github.com/a-h/templ"

	"github.com/jdodson3106/nexus/log"
)

func printAppString() {
	name := ` 
 ________    _______       ___    ___  ___  ___   ________      
|\   ___  \ |\  ___ \     |\  \  /  /||\  \|\  \ |\   ____\     
\ \  \\ \  \\ \   __/|    \ \  \/  / /\ \  \\\  \\ \  \___|_    
 \ \  \\ \  \\ \  \_|/__   \ \    / /  \ \  \\\  \\ \_____  \   
  \ \  \\ \  \\ \  \_|\ \   /     \/    \ \  \\\  \\|____|\  \  
   \ \__\\ \__\\ \_______\ /  /\   \     \ \_______\ ____\_\  \ 
    \|__| \|__| \|_______|/__/ /\ __\     \|_______||\_________\
                          |__|/ \|__|               \|_________|
	`
	fmt.Printf("\033[1;32m%s\033[0m", name)
	fmt.Println()
}

func tidy() error {
	err := exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		log.Error(fmt.Sprintf("Error runing mod tidy :: %s", err))
		return err
	}

	return nil
}

func installTempl() error {
	const templLib = "github.com/a-h/templ"
	out, err := exec.Command("go", "version").Output()
	if err != nil {
		panic(err)
	}

	log.Trace(fmt.Sprintf("Running with Go version %s", out))
	out, err = exec.Command("go", "list", "-deps").Output()
	depBuffer := []byte{}
	deps := []string{}
	for _, v := range out {
		if v == '\n' {
			deps = append(deps, string(depBuffer))
			depBuffer = nil
			continue
		}
		depBuffer = append(depBuffer, v)
	}

	hasTempl := false
	for _, dep := range deps {
		if dep == templLib {
			hasTempl = true
			log.Trace(fmt.Sprintf("%s already installed.", dep))
			break
		}
	}

	if !hasTempl {
		log.Trace(fmt.Sprintf("installing %s", templLib))
		err = exec.Command("go", "install", "github.com/a-h/templ/cmd/templ@latest").Run()
		if err != nil {
			log.Error(fmt.Sprintf("Error installing %s :: %s", templLib, err))
			return err
		}
	}
	tidy()
	return nil
}

func compileTemplates() error {
	return exec.Command("templ", "generate").Run()
}


// reflectiveRender dynamically searches and executes the compiled
// <arg: name>_templ.go file (e.g. todoShow_templ.go). 
//
// If the RenderArgs.RenderFunc is provided then the provide view function is 
// called, otherwise the method is inferred from the template name using the 
// TemplateNameView() templ.Component {...} defintion convention (e.g. TodoShowView()...)
func reflectiveRender(name, path string, args *RenderArgs) (templ.Component, error) {
    if args.RenderFunc == "" {
        args.RenderFunc = ""
    }

    file, err := plugin.Open(path)
    if err != nil {
        log.Error(fmt.Sprintf("Error loading template %s :: %s", name, err))
        return nil, err
    }

    sym, err := file.Lookup(args.RenderFunc) 
    if err != nil {
        log.Error(fmt.Sprintf("Error executing function %s :: %s", args.RenderFunc, err))
        return nil, err
    } 

    if reflect.TypeOf(sym).Kind() == reflect.Func {
        switch f := sym.(type) {
			case *func() templ.Component:
				// Call functions with no parameters
                return (*f)(), nil 
			case *func(map[string]interface{}) templ.Component :
				// Call functions with parameters
				return (*f)(args.Args), nil
            default:
                msg := fmt.Sprintf("Unsupported function definition :: %v", f)
                log.Error(msg)
                return nil, fmt.Errorf(msg)
            }
    }
    t := reflect.TypeOf(sym).Kind().String()
    return nil, fmt.Errorf("Invalid type definition for %s. Expected: func() Received: %s", args.RenderFunc, t) 
}

func renderWithoutFunctions(path string) error {
     
    return nil
}
