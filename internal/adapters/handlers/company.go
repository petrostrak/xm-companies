package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/petrostrak/xm-companies/internal/adapters/kafka/producer"
	"github.com/petrostrak/xm-companies/internal/adapters/repository"
	"github.com/petrostrak/xm-companies/internal/core/domain"
	"github.com/petrostrak/xm-companies/internal/core/services"
	"github.com/petrostrak/xm-companies/utils"
)

const (
	URL = "http://localhost:8082/topics/%s/partitions/%d"
)

type CompanyHandler struct {
	service services.CompanyService
}

func NewCompanyHandler(companyService services.CompanyService) *CompanyHandler {
	return &CompanyHandler{companyService}
}

func (a *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name              string `json:"name"`
		Description       string `json:"description"`
		NumberOfEmployees int    `json:"number_of_employees"`
		Registered        bool   `json:"registered"`
		Type              int64  `json:"type"`
	}

	err := utils.ReadJSON(w, r, &input)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
	}

	company := &domain.Company{
		Name:              input.Name,
		Description:       input.Description,
		NumberOfEmployees: input.NumberOfEmployees,
		Registered:        input.Registered,
		Type:              domain.CompanyType(input.Type),
	}

	err = a.service.Create(company)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/companies/%d", company.ID))

	err = producer.ProduceCompany(company, http.MethodPost)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}

	err = utils.WriteJSON(w, http.StatusCreated, utils.Envelope{"Company": company}, headers)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (a *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request) {
	id := utils.ReadIDParam(r)

	company, err := a.service.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			utils.NotFoundResponse(w, r)
		default:
			utils.ServerErrorResponse(w, r, err)
		}
		return
	}

	var comp struct {
		ID                uuid.UUID `json:"id"`
		Name              string    `json:"name"`
		Description       string    `json:"description"`
		NumberOfEmployees int       `json:"number_of_employees"`
		Registered        bool      `json:"registered"`
		Type              string    `json:"type"`
	}
	comp.ID = company.ID
	comp.Name = company.Name
	comp.Description = company.Description
	comp.NumberOfEmployees = company.NumberOfEmployees
	comp.Registered = company.Registered
	comp.Type = company.Type.ToString()

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"Company": comp}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (a *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request) {
	id := utils.ReadIDParam(r)

	company, err := a.service.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			utils.NotFoundResponse(w, r)
		default:
			utils.ServerErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name              *string `json:"name"`
		Description       *string `json:"description"`
		NumberOfEmployees *int    `json:"number_of_employees"`
		Registered        *bool   `json:"registered"`
		Type              *int64  `json:"type"`
	}

	err = utils.ReadJSON(w, r, &input)
	if err != nil {
		utils.BadRequestResponse(w, r, err)
	}
	if input.Name != nil {
		company.Name = *input.Name
	}
	if input.Description != nil {
		company.Description = *input.Description
	}
	if input.NumberOfEmployees != nil {
		company.NumberOfEmployees = *input.NumberOfEmployees
	}
	if input.Registered != nil {
		company.Registered = *input.Registered
	}
	if input.Type != nil {
		company.Type = domain.CompanyType(*input.Type)
	}

	err = a.service.Update(company)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
		return
	}

	err = producer.ProduceCompany(company, http.MethodPatch)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"Company": company}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}

func (a *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	id := utils.ReadIDParam(r)

	err := a.service.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRecordNotFound):
			utils.NotFoundResponse(w, r)
		default:
			utils.ServerErrorResponse(w, r, err)
		}
		return
	}

	var company domain.Company
	company.ID = id

	err = producer.ProduceCompany(&company, http.MethodDelete)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}

	err = utils.WriteJSON(w, http.StatusOK, utils.Envelope{"message": "Company successfully deleted"}, nil)
	if err != nil {
		utils.ServerErrorResponse(w, r, err)
	}
}
