package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

func (app *Application) routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/companies", func(r chi.Router) {
		r.Get("/{id}", app.CompanyHandler.GetCompany)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(app.AuthenticationToken))
			r.Use(jwtauth.Authenticator)
			r.Post("/", app.CompanyHandler.CreateCompany)
			r.Patch("/{id}", app.CompanyHandler.UpdateCompany)
			r.Delete("/{id}", app.CompanyHandler.DeleteCompany)
		})
	})

	err := chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
	if err != nil {
		app.Logger.Println(err)
	}

	return r
}
