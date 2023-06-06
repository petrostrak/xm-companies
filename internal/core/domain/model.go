package domain

import (
	"github.com/google/uuid"
)

type Company struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	NumberOfEmployees int       `json:"number_of_employees"`
	Registered        bool      `json:"registered"`
	Type              string    `json:"type"`
}
