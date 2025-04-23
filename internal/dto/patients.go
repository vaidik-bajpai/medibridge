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
	Email             string    `json:"email" validate:"email"`
	Address           string    `json:"address" validate:"required,min=5,max=255"`
	EmergencyName     string    `json:"emergencyName" validate:"required"`
	EmergencyRelation string    `json:"emergencyRelation" validate:"required"`
	EmergencyPhone    string    `json:"emergencyPhone" validate:"required,numeric,len=10"`
	RegByID           string    `json:"regBy" validate:"required,uuid"`
}

func (r *RegPatientReq) Sanitize() {
	r.FullName = strings.TrimSpace(r.FullName)
	r.Gender = strings.ToUpper(strings.TrimSpace(r.Gender))
	r.Email = strings.TrimSpace(r.Email)
	r.ContactNumber = strings.TrimSpace(r.ContactNumber)
	r.Address = strings.TrimSpace(r.Address)
	r.EmergencyName = strings.TrimSpace(r.EmergencyName)
	r.EmergencyRelation = strings.TrimSpace(r.EmergencyRelation)
	r.EmergencyPhone = strings.TrimSpace(r.EmergencyPhone)
	r.RegByID = strings.TrimSpace(r.RegByID)
}

type PatientListItem struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Gender   string    `json:"gender"`
	Age      int       `json:"age"`
	DOB      time.Time `json:"dob"`
}
