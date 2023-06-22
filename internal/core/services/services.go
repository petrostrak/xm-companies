package services

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/petrostrak/xm-companies/internal/adapters/kafka/producer"
	"github.com/petrostrak/xm-companies/internal/core/domain"
	"github.com/petrostrak/xm-companies/internal/core/ports"
	"github.com/petrostrak/xm-companies/utils"
)

type CompanyService struct {
	repo ports.CompanyRepository
}

func NewCompanyService(repo ports.CompanyRepository) *CompanyService {
	return &CompanyService{repo}
}

func (c *CompanyService) Create(company *domain.Company) error {
	err := producer.ProduceCompany(company, http.MethodPost)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
	return c.repo.Create(company)
}

func (c *CompanyService) Update(company *domain.Company) error {
	err := producer.ProduceCompany(company, http.MethodPatch)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
	return c.repo.Update(company)
}

func (c *CompanyService) Delete(id uuid.UUID) error {
	err := producer.ProduceCompany(&company, http.MethodDelete)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
	return c.repo.Delete(id)
}

func (c *CompanyService) Get(id uuid.UUID) (*domain.Company, error) {
	return c.repo.Get(id)
}
