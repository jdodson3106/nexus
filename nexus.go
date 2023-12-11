package nexus

import (
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	GET     = "GET"
	POST    = "POST"
	DELETE  = "DELETE"
	PUT     = "PUT"
	PATCH   = "PATCH"
	OPTIONS = "OPTIONS"
)

type Handler func(ctx *Context) error

type TemplateEngine struct {
	// todo: determine what is necessary to build this
}

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	params   httprouter.Params
	ctx      context.Context
}

func NewContext(w http.ResponseWriter, r *http.Request, params httprouter.Params) *Context {
	return &Context{
		request:  r,
		response: w,
		params:   params,
		ctx:      context.TODO(), // not quite sure what to do here yet...thinking context.Background() but needs more investigation
	}
}

type NexusConfig struct {
	// todo define all custom config in here
	// this will be generated during the scaffolding proces
	// if the user provides flags to the cli tool
	Port   string
	Engine TemplateEngine
}

type Nexus struct {
	// Needed to build the engine: router, db, templating,
	router *httprouter.Router
	port   string

	// todo: implement db abstraction
	// todo: decide on default templateing (tmpl with htmx most likely)
}

// the default setup
func NewDefault() (*Nexus, error) {
	return &Nexus{
		router: httprouter.New(),
		port:   ":3000",
	}, nil
}

// the custom setup given cli flags
func New(c NexusConfig) (*Nexus, error) {
	// WIP
	return &Nexus{
		router: httprouter.New(),
		port:   c.Port,
	}, nil
}

func (n *Nexus) Run() error {
	printAppString()
	log.Printf("Nexus server started at http://localhost%s\n", n.port)
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
