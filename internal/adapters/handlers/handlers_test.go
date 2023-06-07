// go:build integration
package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/petrostrak/xm-companies/internal/adapters/repository"
	"github.com/petrostrak/xm-companies/internal/core/services"
)

var (
	resource       *dockertest.Resource
	pool           *dockertest.Pool
	testDB         *sql.DB
	testRepo       repository.PostgresRepository
	companyService *services.CompanyService
	companyHandler *CompanyHandler
)

func TestMain(m *testing.M) {
	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	options := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=postgres",
			"POSTGRES_PASSWORD=password",
			"POSTGRES_DB=xm_companies_test",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: "5436"},
			},
		},
	}

	resource, err = pool.RunWithOptions(&options)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not start resource: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = sql.Open("postgres", "host=localhost port=5436 user=postgres password=password dbname=xm_companies_test sslmode=disable timezone=UTC connect_timeout=5")
		if err != nil {
			log.Println("error:", err)
			return err
		}
		return testDB.Ping()
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to DB: %s", err)
	}

	err = createTables()
	if err != nil {
		log.Fatalf("error creating tables: %s", err)
	}

	testRepo = repository.PostgresRepository{
		CompanyRepository: &repository.CompanyRepository{DB: testDB},
	}

	companyService = services.NewCompanyService(testRepo.CompanyRepository)
	companyHandler = NewCompanyHandler(*companyService)

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tableSQL, err := os.ReadFile("./testdata/init_schema.sql")
	if err != nil {
		fmt.Println(err)
		return err
	}

	_, err = testDB.Exec(string(tableSQL))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func Test_PingDB(t *testing.T) {
	err := testDB.Ping()
	if err != nil {
		t.Error("cannot ping DB")
	}
}

func Test_Handlers(t *testing.T) {
	testCases := []struct {
		name               string
		method             string
		json               string
		paramID            string
		handler            http.HandlerFunc
		expectedStatusCode int
	}{
		{
			"createCompany",
			"POST",
			`{
				"name": "Petros Inc.",
				"number_of_employees": 50,
				"registered": false,
				"type": "Sole Proprietorship"
			  }`,
			"",
			companyHandler.CreateCompany,
			http.StatusCreated,
		},
		{"getCompany", "GET", "", "0e6c0248-a659-41d0-b860-795df3a53f44", companyHandler.GetCompany, http.StatusOK},
		{"getCompany-Invalid", "", "", "121f03cd-ce8c-447d-8747-fb8cb7aa3a52", companyHandler.GetCompany, http.StatusMethodNotAllowed},
		{
			"updateCompany",
			"PATCH",
			`{
				"name": "Google Inc.",
				"description": "A short desc of my company",
				"number_of_employees": 5,
				"registered": true,
				"type": "NonProfit"
			  }`,
			"0e6c0248-a659-41d0-b860-795df3a53f44",
			companyHandler.UpdateCompany,
			http.StatusOK,
		},
		{
			"updateCompany-duplicate name",
			"PATCH",
			`{
				"name": "Petros Inc.",
				"description": "A short desc of my company",
				"number_of_employees": 5,
				"registered": true,
				"type": "NonProfit"
			  }`,
			"0e6c0248-a659-41d0-b860-795df3a53f44",
			companyHandler.UpdateCompany,
			http.StatusInternalServerError,
		},
		{"deleteCompany", "DELETE", "", "0e6c0248-a659-41d0-b860-795df3a53f44", companyHandler.DeleteCompany, http.StatusOK},
	}

	for _, tt := range testCases {
		var req *http.Request
		if tt.json == "" {
			req, _ = http.NewRequest(tt.method, "/", nil)
		} else {
			req, _ = http.NewRequest(tt.method, "/", strings.NewReader(tt.json))
		}

		if tt.paramID != "" {
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", tt.paramID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(tt.handler)
		handler.ServeHTTP(rr, req)

		if rr.Code != tt.expectedStatusCode {
			t.Errorf("%s: wrong status returned; expected %d but got %d", tt.name, tt.expectedStatusCode, rr.Code)
		}
	}
}
