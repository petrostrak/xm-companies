package ports

import (
	"github.com/google/uuid"
	"github.com/petrostrak/xm-companies/core/domain"
)

type CompanyRepository interface {
	Create(*domain.Company) error
	Update(*domain.Company) error
	Delete(uuid.UUID) error
	Get(uuid.UUID) (*domain.Company, error)
}
