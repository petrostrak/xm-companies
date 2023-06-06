package services

import (
	"github.com/google/uuid"
	"github.com/petrostrak/xm-companies/internal/core/domain"
	"github.com/petrostrak/xm-companies/internal/core/ports"
)

type CompanyService struct {
	repo ports.CompanyRepository
}

func NewCompanyService(repo ports.CompanyRepository) *CompanyService {
	return &CompanyService{repo}
}

func (c *CompanyService) Create(company *domain.Company) error {
	return c.repo.Create(company)
}

func (c *CompanyService) Update(company *domain.Company) error {
	return c.repo.Update(company)
}

func (c *CompanyService) Delete(id uuid.UUID) error {
	return c.repo.Delete(id)
}

func (c *CompanyService) Get(id uuid.UUID) (*domain.Company, error) {
	return c.repo.Get(id)
}
