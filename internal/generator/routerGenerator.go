package generator

import (
	"bytes"
	"go/parser"
	"go/token"
	"os"
)

type RouterGenerator struct {
	AppBasePath string
	AppName     string
	RoutesPath  string
	buildPath   string
	routePaths  []string
}

func NewRouterGenerator(basePath, appName, routesPath string) *RouterGenerator {
	return &RouterGenerator{
		AppBasePath: basePath,
		AppName:     appName,
		RoutesPath:  routesPath,
		buildPath:   PLUGINS_DIR_PATH,
	}
}

func (rg *RouterGenerator) GeneratePlugin() (string, error) {
	// find all routes files
	routeFiles, err := walkControllersDir()
	if err != nil {
		return "", err
	}

	fileBuf := make([]byte, 0)
	importsBuf := make([]byte, 0)
	bodyBuf := make([]byte, 0)
	fileBuf = append(fileBuf, []byte("package main\n\nimports (\n")...)

	for _, rf := range routeFiles {
		body, imports, err := parseGoFileWithSeparateImports(rf)
		if err != nil {
			return "", err
		}
		bodyBuf = append(bodyBuf, body...)
		for _, i := range imports {
			impString := "\t" + string(i) + "\n"
			importsBuf = append(importsBuf, impString...)
		}
	}

	fileBuf = append(fileBuf, importsBuf...)
	fileBuf = append(fileBuf, ")\n"...)
	fileBuf = append(fileBuf, bodyBuf...)

	return "", nil
}

func createMain() {
	// TODO: grab all route funcs and generate a large main file for plugin creation
	//os.ReadDir()
}

// parseGoFileWithSeparateImports parses the entire go file and returns the
// file contents and import values as separate variables.
//
// Note: This does NOT return the package declaration as part of the contents. Use parseGoPackageWithPackage
// to get package name.
// Note: The imports return is a list of byte slices that only include the import statements, NOT the import keyword
func parseGoFileWithSeparateImports(fileName string) (contents []byte, imports [][]byte, err error) {
	contentsBuf := bytes.NewBuffer(contents)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	var bodyStart int
	for i, imp := range f.Imports {
		imports = append(imports, []byte(imp.Name.Name))
		if i == len(f.Imports)-1 {
			bodyStart = fset.Position(imp.End()).Offset
		}
	}

	src, err := os.ReadFile(fileName)
	srcBody := src[bodyStart:f.FileEnd]
	contentsBuf.Write(srcBody)
	return contentsBuf.Bytes(), imports, nil
}

func parseGoPackageWithPackage(fileName string) ([]byte, error) {
	return nil, nil
}

func parseImports(fileName string) ([][]byte, error) {
	imports := make([][]byte, 0)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	for _, imp := range f.Imports {
		imports = append(imports, []byte(imp.Name.Name))
	}
	return imports, nil
}

func parseGoFileBody(fileName string) ([]byte, error) {
	contentBuffer := make([]byte, 0)
	writeBuf := bytes.NewBuffer(contentBuffer)
	// parser.Par
	return writeBuf.Bytes(), nil
}

func walkControllersDir() ([]string, error) {
	return nil, nil
}
