package nexus

import "net/http"

type RequestBody struct {
	FormData    map[string][]string
	targetModel *interface{}
}

type Param struct {
	Key   string
	Value string
}

type Params []Param

type Request struct {
	body   RequestBody
	params Params
	ctx    *Context
	//	response *Response
	w http.ResponseWriter
	r *http.Request
}

func (r *Request) HandleCrudRoute() *Response {
	res := &Response{}
	return res
}
