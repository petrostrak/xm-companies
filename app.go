package main

import (
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/petrostrak/xm-companies/internal/adapters/handlers"
	"github.com/petrostrak/xm-companies/internal/adapters/kafka/producer"
	"github.com/petrostrak/xm-companies/internal/adapters/repository"
	"github.com/petrostrak/xm-companies/internal/core/services"
	"github.com/petrostrak/xm-companies/utils"
)

type Application struct {
	Logger              *log.Logger
	Config              *utils.Config
	CompanyService      *services.CompanyService
	CompanyHandler      *handlers.CompanyHandler
	Routes              *chi.Mux
	AuthenticationToken *jwtauth.JWTAuth
	KafkaProducer       *kafka.Producer
}

func (app *Application) New() error {
	config, err := utils.LoadConfig(".")
	if err != nil {
		return err
	}

	store := repository.NewPostgresRepository()
	companyService := services.NewCompanyService(store.CompanyRepository)
	companyHandler := handlers.NewCompanyHandler(*companyService)
	tokenAuth := jwtauth.New("HS256", []byte("xm-companies"), nil)

	prod, err := producer.GetNewProducer()
	if err != nil {
		return err
	}

	app.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	app.Config = config
	app.CompanyService = companyService
	app.CompanyHandler = companyHandler
	app.Routes = app.routes().(*chi.Mux)
	app.AuthenticationToken = tokenAuth
	app.KafkaProducer = prod

	return nil

}
