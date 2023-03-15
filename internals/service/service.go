package service

import (
	"github.com/ZakirAvrora/TechHome/internals/controller"
	customMiddleware "github.com/ZakirAvrora/TechHome/internals/middleware"
	"github.com/ZakirAvrora/TechHome/internals/repository"
	"github.com/ZakirAvrora/TechHome/pkg/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func NewService(conn *sqlx.DB, memCache cache.Cache) http.Handler {
	c := controller.New(repository.NewPostgresRepo(conn), memCache)
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/admin", func(r chi.Router) {
		r.With(customMiddleware.Pagination).Get("/redirects", c.GetAllLinks)
		r.Get("/redirects/{id}", c.GetLink)
		r.Post("/redirects", c.CreateLink)
		r.Patch("/redirects/{id}", c.UpdateLink)
		r.Delete("/redirects/{id}", c.DeleteLink)
	})

	router.Get("/redirects", c.Redirect)

	return router
}
