package main

import (
	"github.com/jdodson3106/nexus"
	"github.com/jdodson3106/nexus/app/views"
)

func main() {
	nx, err := nexus.NewDefault()
	if err != nil {
		panic(err)
	}

	nx.GET("/index", func(ctx *nexus.Context) error {
		return ctx.RenderComponent(views.Hello())
	})

	if err := nx.Run(); err != nil {
		panic(err)
	}
}
