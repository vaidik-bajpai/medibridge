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
	ID                string    `json:"-"`
	FullName          string    `json:"fullname" validate:"omitempty,min=2,max=100"`
	Gender            string    `json:"gender" validate:"omitempty,oneof=MALE FEMALE OTHER"`
	DOB               time.Time `json:"dob" validate:"omitempty"`
	Age               int       `json:"age" validate:"omitempty,numeric,lte=100"`
	ContactNumber     string    `json:"contactNo" validate:"omitempty,numeric,len=10"`
	Address           string    `json:"address" validate:"omitempty,min=5,max=255"`
	EmergencyName     string    `json:"emergencyName" validate:"omitempty"`
	EmergencyRelation string    `json:"emergencyRelation" validate:"omitempty"`
	EmergencyPhone    string    `json:"emergencyPhone" validate:"omitempty,numeric,len=10"`
}

func (p *UpdatePatientReq) Sanitize() {
	p.ID = strings.TrimSpace(p.ID)
	p.FullName = strings.TrimSpace(p.FullName)
	p.Gender = strings.ToUpper(strings.TrimSpace(p.Gender))
	p.ContactNumber = strings.TrimSpace(p.ContactNumber)
	p.Address = strings.TrimSpace(p.Address)
	p.EmergencyName = strings.TrimSpace(p.EmergencyName)
	p.EmergencyRelation = strings.TrimSpace(p.EmergencyRelation)
	p.EmergencyPhone = strings.TrimSpace(p.EmergencyPhone)
}
