package server

import (
	"net/http"

	"github.com/jdodson3106/nexus"
)

type RequestBody struct {
	FormData    map[string]interface{}
	targetModel *interface{}
}

type Request struct {
	body     RequestBody
	params   map[string]interface{}
	ctx      *nexus.Context
	response *Response
	w        http.ResponseWriter
	r        *http.Request
}

func (r *Request) Handle() *Response {
	return r.response
}
