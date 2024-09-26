package nexus

import (
	"github.com/jdodson3106/nexus/log"
	"net/http"
)

type Response struct {
	// The path to the view inside the view folder
	// example: index.html => "index.html" /users/get.html => "/users/get.html"
	View string

	// the vars that are injected into the body of the template
	Model map[string]interface{}

	// status code for the response
	Status int

	// default writer to writer the template to
	w http.ResponseWriter

	// default go http request
	r *http.Request
}

func (r *Response) handle() (*Response, error) {
	// load the template and inject the variables

	// update the writer with the hydrated html file

	// set the codes

	// return the response
	return nil, nil
}

func (r *Response) HandlerFunc() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		res, err := r.handle()
		if err != nil {
			// TODO: Create an error handler config to inject
			log.Error(err.Error())
		}
		w.WriteHeader(res.Status)
		_, _ = w.Write([]byte(res.View))
	}
}
