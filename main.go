package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/lib/pq"
	"github.com/petrostrak/xm-companies/internal/adapters/handlers"
	"github.com/petrostrak/xm-companies/internal/adapters/repository"
	"github.com/petrostrak/xm-companies/internal/core/services"
)

var (
	companyService *services.CompanyService
	companyHandler *handlers.CompanyHandler
	tokenAuth      *jwtauth.JWTAuth
)

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	store := repository.NewPostgresRepository()
	companyService = services.NewCompanyService(store.CompanyRepository)
	companyHandler = handlers.NewCompanyHandler(*companyService)
	tokenAuth = jwtauth.New("HS256", []byte("xm-companies"), nil)

	srv := &http.Server{
		Addr:        fmt.Sprintf(":%d", 8080),
		Handler:     Routes(),
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
	}

	logger.Printf("starting development server on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}

func Routes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/companies", func(r chi.Router) {
		r.Get("/{id}", companyHandler.GetCompany)

		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(tokenAuth))
			r.Use(jwtauth.Authenticator)
			r.Post("/", companyHandler.CreateCompany)
			r.Patch("/{id}", companyHandler.UpdateCompany)
			r.Delete("/{id}", companyHandler.DeleteCompany)
		})
	})

	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})

	return r
}
