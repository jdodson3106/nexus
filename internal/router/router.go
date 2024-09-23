package router

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/jdodson3106/nexus/internal/server"
)

type RouteHandler func(*server.Request)

type Route struct {
	Path    string
	Handler RouteHandler
	Method  string
	Name    string
}

type Router struct {
	BasePath string
	routes   []Route
}

func NewRouter(path string, routes ...Route) *Router {
	return &Router{BasePath: path, routes: routes}
}

// TODO: Make this accept the ORM interface
func (r *Router) NewModelCrudRoutes(model interface{}) {
	// name := reflect.Type(model)
}

// NOTE: May end up killing this and just moving into the cli for scaffolding
func (r *Router) NewCrudRoutes(modelName string) error {
	cruds := []Route{
		{
			Path:    fmt.Sprintf("%s/%s/:id", r.BasePath, modelName),
			Handler: func(req *server.Request) { req.Handle() },
			Method:  http.MethodGet,
			Name:    fmt.Sprintf("get_%s", modelName),
		},
		{
			Path:    fmt.Sprintf("%s/%s/:id", r.BasePath, modelName),
			Handler: func(req *server.Request) { req.Handle() },
			Method:  http.MethodPost,
			Name:    fmt.Sprintf("create_%s", modelName),
		},
		{
			Path:    fmt.Sprintf("%s/%s/:id", r.BasePath, modelName), // TODO: figure out query and path params
			Handler: func(req *server.Request) { req.Handle() },
			Method:  http.MethodPost,
			Name:    fmt.Sprintf("update_%s", modelName),
		},
		{
			Path:    fmt.Sprintf("%s/%s", r.BasePath, modelName),
			Handler: func(req *server.Request) { req.Handle() },
			Method:  http.MethodGet,
			Name:    fmt.Sprintf("delete_%s", modelName),
		},
	}

	ok, existingPaths := r.checkForExistingRoutes(cruds...)
	if !ok {
		return errors.New(fmt.Sprintf("routes alread exist. cannot create :: %+v", existingPaths))
	}

	r.routes = append(r.routes, cruds...)
	return nil
}

func (r *Router) NewRoute(rt Route) {
	r.routes = append(r.routes, rt)
}

func (r Router) checkForExistingRoutes(newRoutes ...Route) (ok bool, badPaths []Route) {
	for _, p := range newRoutes {
		isBad := false
		for _, route := range r.routes {
			if p.Path == route.Path && p.Method == route.Method {
				isBad = true
				badPaths = append(badPaths, p)
				break
			}
		}
		if isBad {
			continue
		}
	}
	return len(badPaths) > 0, badPaths
}
