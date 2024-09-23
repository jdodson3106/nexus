package controllers

import (
	"github.com/jdodson3106/nexus"
	"github.com/jdodson3106/nexus/app/controllers"
	"github.com/jdodson3106/nexus/app/models"
)

func HomeRouter(r *nexus.Router) {
	/* Will create all new routes:
	View all - [GET: /blog]
	Create - [POST: /blog]
	Update - [POST: /blog/:id]
	Get by ID - [GET: /blog/:id]
	Delete - [GET: /blog/delete/:id]
	*/
	r.NewModelCrudRoutes(&models.Blog{})

}

func ModelRouter(r *nexus.Router) {
	model := r.NewRouteGroup("/model")
	model.GET("/create", controllers.View)
}
