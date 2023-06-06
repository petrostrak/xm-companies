package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/petrostrak/xm-companies/internal/core/domain"
)

var (
	ErrRecordNotFound = errors.New("record not Found")
)

type PostgresRepository struct {
	*CompanyRepository
}

var dsn = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"

func NewPostgresRepository() *PostgresRepository {

	db, err := sql.Open("postgres", fmt.Sprintf(dsn,
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB_NAME"),
	))
	if err != nil {
		return nil
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	duration, err := time.ParseDuration("15m")
	if err != nil {
		return nil
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil
	}

	return &PostgresRepository{
		&CompanyRepository{db},
	}
}

type CompanyRepository struct {
	DB *sql.DB
}

func (a *CompanyRepository) Create(company *domain.Company) error {
	query := `
		INSERT INTO companies (name, description, number_of_employees, registered, type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, number_of_employees, registered, type`

	args := []any{company.Name, company.Description, company.NumberOfEmployees, company.Registered, company.Type}

	return a.DB.QueryRow(query, args...).Scan(
		&company.ID,
		&company.Name,
		&company.Description,
		&company.NumberOfEmployees,
		&company.Registered,
		&company.Type,
	)
}

func (a *CompanyRepository) Get(id uuid.UUID) (*domain.Company, error) {
	query := `
		SELECT id, name, description, number_of_employees, registered, type
		FROM companies
		WHERE id = $1`

	var company domain.Company

	err := a.DB.QueryRow(query, id).Scan(
		&company.ID,
		&company.Name,
		&company.Description,
		&company.NumberOfEmployees,
		&company.Registered,
		&company.Type,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &company, nil
}

func (a *CompanyRepository) Update(company *domain.Company) error {
	query := `
		UPDATE companies
		SET name = $1, description = $2, number_of_employees = $3, registered = $4, type = $5
		WHERE id = $6
		RETURNING id, name, description, number_of_employees, registered, type`

	args := []any{
		company.Name,
		company.Description,
		company.NumberOfEmployees,
		company.Registered,
		company.Type,
		company.ID,
	}

	return a.DB.QueryRow(query, args...).Scan(
		&company.ID,
		&company.Name,
		&company.Description,
		&company.NumberOfEmployees,
		&company.Registered,
		&company.Type,
	)
}

func (a *CompanyRepository) Delete(id uuid.UUID) error {
	query := `
		DELETE FROM companies
		WHERE id = $1`

	result, err := a.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
