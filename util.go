package nexus

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"reflect"

	"github.com/a-h/templ"

	"github.com/jdodson3106/nexus/log"
)

func printAppString() {
	name :=
		`
+--------------------------------------------------------------------------+
|       ________    _______       ___    ___  ___  ___   ________          |
|      |\   ___  \ |\  ___ \     |\  \  /  /||\  \|\  \ |\   ____\         |
|      \ \  \\ \  \\ \   __/|    \ \  \/  / /\ \  \\\  \\ \  \___|_        |
|       \ \  \\ \  \\ \  \_|/__   \ \    / /  \ \  \\\  \\ \_____  \       |
|        \ \  \\ \  \\ \  \_|\ \   /     \/    \ \  \\\  \\|____|\  \      |
|         \ \__\\ \__\\ \_______\ /  /\   \     \ \_______\ ____\_\  \     |
|          \|__| \|__| \|_______|/__/ /\ __\     \|_______||\_________\    |
|                                |__|/ \|__|               \|_________|    |
|                                                                          |
|                     An opinionated Web Framework in Go                   |
+--------------------------------------------------------------------------+
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
	// get all *_templ.go files, copy to a local tmp/main.go, then compile to plugin
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
		// TODO: Define func name using name & path convention
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
		case *func(map[string]interface{}) templ.Component:
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

func generatePlugin(path string) (string, error) {
	if buildDirPath == "" {
		buildDirPath = getRelativeAppPath() + "out"
	}
	err := os.MkdirAll(buildDirPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create build directory: %v", err)
	}
	baseName := filepath.Base(path)
	pluginName := baseName[:len(baseName)-len(filepath.Ext(baseName))] + ".so"

	// Construct the output path
	outputPath := filepath.Join(buildDirPath, pluginName)

	// Build the command
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputPath, path)
	fmt.Println(cmd)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to build plugin: %v", err)
	}

	fmt.Printf("Plugin built successfully: %s\n", outputPath)
	return outputPath, nil
}

func reflectiveRouteLoader(path string, router *Router) error {
	routeFile := fmt.Sprintf("%s/routes.go", path)
	funcs, err := parseFunctionsFromRouter(routeFile)
	if err != nil {
		return err
	}

	pluginFile, err := generatePlugin(routeFile)
	if err != nil {
		return err
	}

	file, err := plugin.Open(pluginFile)
	if err != nil {
		log.Error(fmt.Sprintf("error loading router file %s :: %s", pluginFile, err))
		return err
	}
	for _, fn := range funcs {
		sym, err := file.Lookup(fn)
		if err != nil {
			log.Error(fmt.Sprintf("error executing function %s :: %s", fn, err))
			return err
		}

		if reflect.TypeOf(sym).Kind() == reflect.Func {
			switch f := sym.(type) {
			case *func():
				msg := fmt.Sprintf("error calling router func: %s. all routers functions expect a *nexus.Router parameter", fn)
				return fmt.Errorf(msg)
			case *func(router *Router):
				(*f)(router)
			default:
				msg := fmt.Sprintf("Unsupported function definition :: %v", f)
				log.Error(msg)
				return fmt.Errorf(msg)
			}
		}
	}
	return nil
}

func parseFunctionsFromRouter(path string) ([]string, error) {
	fset := token.NewFileSet()

	node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	var functions []string
	ast.Inspect(node, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			functions = append(functions, funcDecl.Name.Name)
		}
		return true
	})
	return functions, nil
}
