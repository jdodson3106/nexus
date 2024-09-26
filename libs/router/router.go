package router

import (
	"context"
	"net/http"
)

type NexusContext struct {
	ctx context.Context
	w   http.ResponseWriter
	r   *http.Request
}

type NexusHandler func(ctx *NexusContext) error

func NexusRoute(path string, handler NexusHandler, middleware ...NexusHandler)
