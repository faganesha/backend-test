package main

import (
	"backend-test/controller"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/top-stories", controller.InsertTopStories)
	r.Post("/detail-stories", controller.InsertTopStoryDetail)
	r.Get("/get-top-stories", controller.GetTopStories)
	r.Get("/get-top-story/:id", controller.GetDetailStories)

	http.ListenAndServe(":3000", r)
}
