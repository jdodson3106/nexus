package nexus

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jdodson3106/nexus/internal/db"
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

// HandlerFunc the stub of a handler method that should be passed in any Nexus route
type HandlerFunc func(ctx *Context) error

type RenderArgs struct {
	Args       map[string]interface{}
	RenderFunc string
}

func (n *Nexus) CloseDatabase() error {
	return n.db.Close()
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

	db, err := db.NewDefaultDbConnection()
	if err != nil {
		return nil, err
	}

	return &Nexus{
		router:  NewRouter(""),
		port:    ":8080",
		appName: DEF_APP_NAME,
		db:      db,
	}, nil
}

// the custom setup given cli flags
func New(c NexusConfig) (*Nexus, error) {
	// WIP
	viewsPath = c.ViewPath
	return &Nexus{
		router: NewRouter(""),
		port:   c.Port,
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

	printAppString()
	log.Info(fmt.Sprintf("Nexus server started at http://localhost%s", n.port))
	return http.ListenAndServe(n.port, n.router.httpRouter)
}

func (n *Nexus) GET(path string, handler HandlerFunc) {
	n.createHttpHandle(GET, path, handler)
}

func (n *Nexus) POST(path string, handler HandlerFunc) {
	n.createHttpHandle(POST, path, handler)
}

func (n *Nexus) PUT(path string, handler HandlerFunc) {
	n.createHttpHandle(PUT, path, handler)
}

func (n *Nexus) PATCH(path string, handler HandlerFunc) {
	n.createHttpHandle(PATCH, path, handler)
}

func (n *Nexus) DELETE(path string, handler HandlerFunc) {
	n.createHttpHandle(DELETE, path, handler)
}

// TODO: Git rid of all this and use updates in the new Router
func (n *Nexus) createHttpHandle(method string, path string, handler HandlerFunc) {
	handle := func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		var newParams Params
		for _, oldParam := range p {
			newParams = append(newParams, Param{Key: oldParam.Key, Value: oldParam.Value})
		}
		ctx := NewContext(w, r, newParams)
		if err := handler(ctx); err != nil {
			// todo: figure out handling errors in the handler calls;
			//       need a default mechanism to handle these gracefully
			panic(err)
		}
	}
	n.router.Handle(method, path, handle)
}
