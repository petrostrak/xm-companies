package domain

import (
	"strings"

	"github.com/google/uuid"
)

type Company struct {
	ID                uuid.UUID   `json:"id"`
	Name              string      `json:"name"`
	Description       string      `json:"description"`
	NumberOfEmployees int         `json:"number_of_employees"`
	Registered        bool        `json:"registered"`
	Type              CompanyType `json:"type"`
}

type CompanyType int64

const (
	Corporations CompanyType = iota
	NonProfit
	Cooperative
	SoleProprietorship
	Unknown
)

func (c CompanyType) ToString() string {
	switch c {
	case Corporations:
		return "Corporations"
	case NonProfit:
		return "Non Profit"
	case Cooperative:
		return "Cooperative"
	default:
		return "Sole Proprietorship"
	}
}

func FromString(s string) CompanyType {
	switch strings.ToLower(s) {
	case "corporations":
		return Corporations
	case "non profit":
		return NonProfit
	case "cooperative":
		return Cooperative
	case "sole proprietorship":
		return SoleProprietorship
	default:
		return Unknown
	}
}
