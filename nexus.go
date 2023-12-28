package nexus

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/a-h/templ"
	log "github.com/jdodson3106/nexus/log"
	"github.com/julienschmidt/httprouter"
)

var viewsPath string

const (
	GET     = "GET"
	POST    = "POST"
	DELETE  = "DELETE"
	PUT     = "PUT"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"

	DEF_APP_NAME = "app"
)

// Handler the stub of a handler method that should be passed in any Nexus route
type Handler func(ctx *Context) error

type RenderArgs struct {
	Args       map[string]interface{}
	RenderFunc string
}

type TemplateEngine struct {
	// stub for the eventual multiple templates support

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
	Port     string
	Engine   TemplateEngine
	ViewPath string
}

type Nexus struct {
	// Needed to build the engine: router, db, templating,
	router  *httprouter.Router
	port    string
	appName string

	// todo: implement db abstraction
}

// the default setup
func NewDefault() (*Nexus, error) {
	dir, err := os.Getwd()
	if err != nil {
		dir = fmt.Sprintf("./%s", DEF_APP_NAME)
	}
	p := getPathVar()
	dir += p
	dir += "/views"
	viewsPath = dir

	return &Nexus{
		router:  httprouter.New(),
		port:    ":3000",
		appName: DEF_APP_NAME,
	}, nil
}

// the custom setup given cli flags
func New(c NexusConfig) (*Nexus, error) {
	// WIP
	viewsPath = c.ViewPath
	return &Nexus{
		router: httprouter.New(),
		port:   c.Port,
	}, nil
}
func (n *Nexus) Run() error {
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

	printAppString()
	log.Info(fmt.Sprintf("Nexus server started at http://localhost%s", n.port))
	return http.ListenAndServe(n.port, n.router)
}

func (n *Nexus) GET(path string, handler Handler) {
	n.createHttpHandle(GET, path, handler)
}

func (n *Nexus) POST(path string, handler Handler) {
	n.createHttpHandle(POST, path, handler)
}

func (n *Nexus) PUT(path string, handler Handler) {
	n.createHttpHandle(PUT, path, handler)
}

func (n *Nexus) PATCH(path string, handler Handler) {
	n.createHttpHandle(PATCH, path, handler)
}

func (n *Nexus) DELETE(path string, handler Handler) {
	n.createHttpHandle(DELETE, path, handler)
}

func (n *Nexus) createHttpHandle(method string, path string, handler Handler) {
	handle := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := NewContext(w, r, p)
		if err := handler(ctx); err != nil {
			// todo: figure out handling errors in the handler calls;
			//       need a default mechanism to handle these gracefully
			panic(err)
		}
	}
	n.router.Handle(method, path, handle)
}
