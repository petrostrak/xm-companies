// go:build integration
package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/petrostrak/xm-companies/internal/core/domain"
)

var (
	resource *dockertest.Resource
	pool     *dockertest.Pool
	testDB   *sql.DB
	testRepo PostgresRepository
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
				{HostIP: "0.0.0.0", HostPort: "5435"},
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
		testDB, err = sql.Open("postgres", "host=localhost port=5435 user=postgres password=password dbname=xm_companies_test sslmode=disable timezone=UTC connect_timeout=5")
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

	testRepo = PostgresRepository{
		&CompanyRepository{testDB},
	}

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

var testCompanyID = uuid.MustParse("0e6c0248-a659-41d0-b860-795df3a53f44")

func Test_PostgresDBRepoGetCompany(t *testing.T) {
	company, err := testRepo.CompanyRepository.Get(testCompanyID)
	if err != nil {
		t.Errorf("error getting company by id: %s", err)
	}

	if company.Name != "Petros Trak Inc" {
		t.Errorf("wrong company name returned. expected 'Petros Trak Inc' but got %s", company.Name)
	}

	if company.NumberOfEmployees != 4 {
		t.Errorf("wrong company number of employees returned. expected 4 but got %v", company.NumberOfEmployees)
	}

	if company.Type.ToString() != "Sole Proprietorship" {
		t.Errorf("wrong company type returned. expected 'Sole Proprietorship' but got %s", company.Name)
	}
}

func Test_PostgresDBRepoCreateCompany(t *testing.T) {
	testCompany := domain.Company{
		Name:              "Golang inc",
		Description:       "A small family firm",
		NumberOfEmployees: 4,
		Registered:        true,
		Type:              domain.SoleProprietorship,
	}

	err := testRepo.CompanyRepository.Create(&testCompany)
	if err != nil {
		t.Errorf("insert company returned an error: %s", err)
	}
}

func Test_PostgresDBRepoUpdateCompany(t *testing.T) {
	company, _ := testRepo.CompanyRepository.Get(testCompanyID)
	company.NumberOfEmployees = 6
	company.Type = domain.Cooperative

	err := testRepo.CompanyRepository.Update(company)
	if err != nil {
		t.Errorf("error updating company: %s", err)
	}

	company, _ = testRepo.CompanyRepository.Get(testCompanyID)
	if company.NumberOfEmployees != 6 || company.Type != domain.Cooperative {
		t.Errorf("expected updated record to have 6 number of employees and Cooperative type, but got %v and %d", company.NumberOfEmployees, company.Type)
	}
}

func Test_PostgresDBRepoDeleteCompany(t *testing.T) {
	err := testRepo.CompanyRepository.Delete(testCompanyID)
	if err != nil {
		t.Errorf("error deleting company: %s", err)
	}

	_, err = testRepo.CompanyRepository.Get(testCompanyID)
	if err == nil {
		t.Errorf("got company %v, which should have been deleted", testCompanyID)
	}
}
