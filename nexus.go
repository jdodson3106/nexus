package nexus

import (
	"context"
	"fmt"
	"github.com/a-h/templ"
	"github.com/jdodson3106/nexus/internal/db"
	"github.com/jdodson3106/nexus/log"
	"net/http"
)

const (
	GET     = "GET"
	POST    = "POST"
	DELETE  = "DELETE"
	PUT     = "PUT"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"

	DEF_APP_NAME = "app"
)

type RenderArgs struct {
	Args       map[string]interface{}
	RenderFunc string
}

type Nexus struct {
	port    string
	router  *Router
	appName string
	config  NexusConfig

	// TODO: implement db abstraction
	db *db.DB
}

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Context  context.Context
	Params   Params
	context  context.Context
}

func (ctx *Context) Render(name string, args *RenderArgs) error {
	path := fmt.Sprintf("%s/%s_templ.go", viewsPath, name)
	comp, err := reflectiveRender(name, path, args)
	if err != nil {
		return err
	}
	return comp.Render(ctx.context, ctx.Response)
}

func (ctx *Context) RenderComponent(component templ.Component) error {
	return component.Render(ctx.context, ctx.Response)
}

func NewContext(w http.ResponseWriter, r *http.Request, params Params) *Context {
	return &Context{
		Request:  r,
		Response: w,
		Params:   params,
		context:  context.Background(),
	}
}

// NexusConfig: Defines all custom config in here
// this will be generated during the scaffolding proces
// if the user provides flags to the cli tool
type NexusConfig struct {
	Port          string
	ViewPath      string
	ControllerDir string
	ModelsDir     string
}

func DefaultNexusConfig() NexusConfig {
	return NexusConfig{
		Port:          ":8080",
		ViewPath:      getRelativeViewsPath(),
		ControllerDir: getRelativeControllerPath(),
	}
}

type controller struct {
	path   string
	name   string
	routes []Route
}

type controllerRegister struct {
	controllers []controller
}

func InitNexus() (*Nexus, error) {
	// create the new nexus instance
	conf := DefaultNexusConfig()
	n := &Nexus{port: ":8080", appName: DEF_APP_NAME, config: conf}

	// register the existing routes and controller methods (handlers)
	n.router = NewRouter("/")

	// parse the controller files
	err := n.parseControllers()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return n, nil
}

func InitNexusWithConfig(config NexusConfig) (*Nexus, error) {
	viewsPath = config.ViewPath
	return &Nexus{
		appName: DEF_APP_NAME,
		router:  NewRouter(""),
		port:    config.Port,
	}, nil
}

func (n *Nexus) Run() error {
	defer n.db.Close()
	err := tidy()
	if err != nil {
		panic(err)
	}

	err = installTempl()
	if err != nil {
		panic(err)
	}

	err = compileTemplates()
	if err != nil {
		panic(err)
	}

	// TODO: Call all found router methods
	n.generateRoutes(n.router)

	for _, r := range n.router.routes {
		log.Info(fmt.Sprintf("Registering route: [%s] - %s", r.Method, r.Path))
	}

	printAppString()
	log.Info(fmt.Sprintf("Nexus server started at http://localhost%s", n.port))
	return http.ListenAndServe(n.port, n.router.httpRouter)
}

func (n *Nexus) CloseDatabase() error {
	return n.db.Close()
}

func (n *Nexus) parseControllers() error {
	// TODO: search through the application's path and find the controller dir
	// and hash the files and store the hashes in a application data file
	// diff the hashes of each file to determine if the file needs to be reparsed

	return nil
}

func (n *Nexus) generateRoutes(router *Router) {
	path := getRelativeControllerPath()
	if err := reflectiveRouteLoader(path, router); err != nil {
		log.Error(err.Error())
	}
}
