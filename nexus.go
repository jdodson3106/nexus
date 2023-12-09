package nexus

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Handler func() error

type Context struct {
	request  *http.Request
	response http.ResponseWriter
	ctx      context.Context
}

type Nexus struct {
	// Needed to build the engine: router, context, db, templating,
	router  *httprouter.Router
	context Context

	// todo: implement gorm db abstraction
	// todo: decide on default templateing (tmpl with htmx most likely)
}

type NexusConfig struct {
	// todo define all custom config in here
	// this will be generated during the scaffolding proces
	// if the user provides flags to the cli tool
}

// the default setup
func NewDefault() (*Nexus, error) {
	return &Nexus{}, nil
}

// the custom setup given cli flags
func New(c NexusConfig) (*Nexus, error) {
	return &Nexus{}, nil
}
