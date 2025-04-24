package dto

import (
	"strings"
	"time"
)

type RegPatientReq struct {
	FullName          string    `json:"fullname" validate:"required,min=2,max=100"`
	Gender            string    `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`
	DOB               time.Time `json:"dob" validate:"required"`
	Age               int       `json:"age" validate:"required,numeric,lte=100"`
	ContactNumber     string    `json:"contactNo" validate:"required,numeric,len=10"`
	Address           string    `json:"address" validate:"required,min=5,max=255"`
	EmergencyName     string    `json:"emergencyName" validate:"required"`
	EmergencyRelation string    `json:"emergencyRelation" validate:"required"`
	EmergencyPhone    string    `json:"emergencyPhone" validate:"required,numeric,len=10"`
	RegByID           string    `json:"-"`
}

func (r *RegPatientReq) Sanitize() {
	r.FullName = strings.TrimSpace(r.FullName)
	r.Gender = strings.ToUpper(strings.TrimSpace(r.Gender))
	r.ContactNumber = strings.TrimSpace(r.ContactNumber)
	r.Address = strings.TrimSpace(r.Address)
	r.EmergencyName = strings.TrimSpace(r.EmergencyName)
	r.EmergencyRelation = strings.TrimSpace(r.EmergencyRelation)
	r.EmergencyPhone = strings.TrimSpace(r.EmergencyPhone)
}

type PatientListItem struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Gender   string    `json:"gender"`
	Age      int       `json:"age"`
	DOB      time.Time `json:"dob"`
}

type UpdatePatientReq struct {
	ID                string     `json:"-"`
	FullName          *string    `json:"fullname" validate:"omitempty,min=2,max=100"`
	Gender            *string    `json:"gender" validate:"omitempty,oneof=MALE FEMALE OTHER"`
	DOB               *time.Time `json:"dob" validate:"omitempty"`
	Age               *int       `json:"age" validate:"omitempty,numeric,lte=100"`
	ContactNumber     *string    `json:"contactNo" validate:"omitempty,numeric,len=10"`
	Address           *string    `json:"address" validate:"omitempty,min=5,max=255"`
	EmergencyName     *string    `json:"emergencyName" validate:"omitempty"`
	EmergencyRelation *string    `json:"emergencyRelation" validate:"omitempty"`
	EmergencyPhone    *string    `json:"emergencyPhone" validate:"omitempty,numeric,len=10"`
}

func (p *UpdatePatientReq) Sanitize() {
	p.ID = strings.TrimSpace(p.ID)

	if p.FullName != nil {
		trimmed := strings.TrimSpace(*p.FullName)
		p.FullName = &trimmed
	}

	if p.Gender != nil {
		trimmed := strings.TrimSpace(*p.Gender)
		upper := strings.ToUpper(trimmed)
		p.Gender = &upper
	}

	if p.ContactNumber != nil {
		trimmed := strings.TrimSpace(*p.ContactNumber)
		p.ContactNumber = &trimmed
	}

	if p.Address != nil {
		trimmed := strings.TrimSpace(*p.Address)
		p.Address = &trimmed
	}

	if p.EmergencyName != nil {
		trimmed := strings.TrimSpace(*p.EmergencyName)
		p.EmergencyName = &trimmed
	}

	if p.EmergencyRelation != nil {
		trimmed := strings.TrimSpace(*p.EmergencyRelation)
		p.EmergencyRelation = &trimmed
	}

	if p.EmergencyPhone != nil {
		trimmed := strings.TrimSpace(*p.EmergencyPhone)
		p.EmergencyPhone = &trimmed
	}
}
