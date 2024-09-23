package nexus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/jdodson3106/nexus/internal/db"
	"github.com/jdodson3106/nexus/internal/router"
	"github.com/jdodson3106/nexus/log"
	"github.com/julienschmidt/httprouter"
)

type Nexus struct {
	// router  *httprouter.Router TODO: *httprouter.Router moved to internal nexus Router
	router  router.Router
	port    string
	appName string
	config  NexusConfig

	// TODO: implement db abstraction
	db *db.DB
}

type Context struct {
	Request  *http.Request
	Response http.ResponseWriter
	Params   httprouter.Params
	ctx      context.Context
}

func (ctx *Context) Render(name string, args *RenderArgs) error {
	path := fmt.Sprintf("%s/%s_templ.go", viewsPath, name)
	comp, err := reflectiveRender(name, path, args)
	if err != nil {
		return err
	}
	return comp.Render(ctx.ctx, ctx.Response)
}

func (ctx *Context) RenderComponent(component templ.Component) error {
	return component.Render(ctx.ctx, ctx.Response)
}

func NewContext(w http.ResponseWriter, r *http.Request, params httprouter.Params) *Context {
	return &Context{
		Request:  r,
		Response: w,
		Params:   params,
		ctx:      context.Background(),
	}
}

// NexusConfig: Defines all custom config in here
// this will be generated during the scaffolding proces
// if the user provides flags to the cli tool
type NexusConfig struct {
	Port          string
	ViewPath      string
	controllerDir string
}

type controller struct {
	path   string
	name   string
	routes []router.Route
}

type controllerRegister struct {
	controllers []controller
}

func InitNexus() (*Nexus, error) {
	n := &Nexus{}

	// register the existing routes and controller methods (handlers)
	n.router = router.Router{BasePath: "/"}

	// TODO: Call all found router methods

	// parse the controller files
	err := n.parseControllers()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return n, nil
}

func (n *Nexus) parseControllers() error {
	// TODO: search through the application's path and find the controller dir
	// and hash the files and store the hashes in a application data file
	// diff the hashes of each file to determine if the file needs to be reparsed

	return nil
}
