package nexus

import (
	"errors"
	"fmt"
	ar "github.com/jdodson3106/nexus/internal/activeRecord"
	"github.com/jdodson3106/nexus/log"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"reflect"
	"text/template"
)

type RouteHandler func(*Request) *Response

type Route struct {
	Path    string
	Handler RouteHandler
	Method  string
	Name    string
	context Context
}

type Router struct {
	BasePath   string
	routes     []Route
	groups     []Router
	httpRouter *httprouter.Router
}

func NewRouter(path string, routes ...Route) *Router {
	return &Router{
		BasePath:   path,
		routes:     routes,
		httpRouter: httprouter.New(),
	}
}

func (r *Router) NewModelCrudRoutes(model ar.ActiveRecord) {
	t := reflect.ValueOf(model)
	if t.Kind() != reflect.Ptr {
		t = t.Elem()
	}
	name := t.Type().Name()
	if err := r.NewCrudRoutes(name, model); err != nil {
		log.Error(err.Error())
	}
}

func (r *Router) NewRouteGroup(prefix string) *Router {
	newRouter := *r
	newRouter.BasePath = r.BasePath + "/" + prefix
	_ = append(r.groups, newRouter)
	return &newRouter
}

// NOTE: May end up killing this and just moving into the cli for scaffolding
func (r *Router) NewCrudRoutes(modelName string, models ...ar.ActiveRecord) error {
	cruds := []Route{
		{
			Path:    fmt.Sprintf("%s/%s", r.BasePath, modelName),
			Handler: r.defaultViewHandler(models[0]),
			Method:  http.MethodGet,
			Name:    fmt.Sprintf("view_%s", modelName),
		},
		{
			Path:    fmt.Sprintf("%s/%s/:id", r.BasePath, modelName),
			Handler: func(req *Request) *Response { return &Response{} },
			Method:  http.MethodGet,
			Name:    fmt.Sprintf("get_%s", modelName),
		},
		{
			Path:    fmt.Sprintf("%s/%s/:id", r.BasePath, modelName),
			Handler: func(req *Request) *Response { return &Response{} },
			Method:  http.MethodPost,
			Name:    fmt.Sprintf("create_%s", modelName),
		},
		{
			Path:    fmt.Sprintf("%s/%s/:id", r.BasePath, modelName), // TODO: figure out query and path params
			Handler: func(req *Request) *Response { return &Response{} },
			Method:  http.MethodPost,
			Name:    fmt.Sprintf("update_%s", modelName),
		},
		{
			Path:    fmt.Sprintf("%s/%s", r.BasePath, modelName),
			Handler: func(req *Request) *Response { return &Response{} },
			Method:  http.MethodGet,
			Name:    fmt.Sprintf("delete_%s", modelName),
		},
	}

	for _, crud := range cruds {
		if err := r.registerNewCrudRoute(crud); err != nil {
			return err
		}
	}

	return nil
}

func (r *Router) newRoute(rt Route) error {
	ok, _ := r.checkForExistingRoutes(rt)
	if !ok {
		return errors.New(fmt.Sprintf("routes alread exist. cannot create :: %+v", rt))
	}
	r.routes = append(r.routes, rt)
	return nil
}

func (r *Router) registerNewCrudRoute(rt Route) error {
	if err := r.newRoute(rt); err != nil {
		return err
	}
	switch rt.Method {
	case http.MethodGet:
		r.GET(rt.Path, rt.Handler)
	case http.MethodPost:
		r.POST(rt.Path, rt.Handler)
	case http.MethodPut:
		r.PUT(rt.Path, rt.Handler)
	case http.MethodDelete:
		r.DELETE(rt.Path, rt.Handler)
	}

	return nil
}

func (r *Router) GET(path string, handler func(r *Request) *Response) {
	if err := r.newRoute(Route{Path: path, Handler: handler}); err != nil {
		log.Error(fmt.Sprintf("a route at %s already exists. will use that definition instead.", path))
		return
	}

	// create the internal handler func
	h := func(wr http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			wr.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		params := r.getParams(req)
		nexReq := &Request{
			body:   RequestBody{FormData: req.Form},
			ctx:    NewContext(wr, req, params),
			params: params,
			w:      wr,
			r:      req,
		}

		// call the user defined handler
		res := handler(nexReq)
		wr.WriteHeader(res.Status)

		// TODO: write the view generator to parse html files
		// 		 for now just use the built-ins
		temp, err := template.ParseFiles(res.View)
		if err != nil {
			wr.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = temp.Execute(wr, res.Model)
		if err != nil {
			wr.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	r.httpRouter.HandlerFunc(http.MethodGet, path, h)
}

func (r *Router) POST(path string, handler RouteHandler) {

}

func (r *Router) PUT(path string, handler RouteHandler) {

}

func (r *Router) DELETE(path string, handler RouteHandler) {

}

func (r *Router) defaultViewHandler(record ar.ActiveRecord) RouteHandler {
	return func(req *Request) *Response {
		rec, err := record.Get()
		if err != nil {
			// TODO: Build some default errors to handle here
			// TODO: Build default views to return
			return &Response{Status: http.StatusInternalServerError}
		}

		// TODO: parse the record's name and get the default views
		model := map[string]interface{}{"modelName": rec}
		return &Response{View: "view_model", Model: model, Status: http.StatusOK}
	}
}

func (r *Router) checkForExistingRoutes(newRoutes ...Route) (ok bool, badPaths []Route) {
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

func (r *Router) getParams(req *http.Request) Params {
	var params Params

	for k, v := range req.URL.Query() {
		params = append(params, Param{Key: k, Value: v[0]})
	}

	return params
}

func (r *Router) Handle(method string, path string, handler httprouter.Handle) {

}
