package models

import (
	"strings"
	"time"
)

type Patient struct {
	ID                string     `json:"id"`
	FullName          string     `json:"fullname"`
	Gender            string     `json:"gender"`
	DOB               DateOnly   `json:"dob"`
	Age               int        `json:"age"`
	ContactNumber     string     `json:"contactNo"`
	Address           string     `json:"address"`
	EmergencyName     string     `json:"emergencyName"`
	EmergencyPhone    string     `json:"emergencyPhone"`
	EmergencyRelation string     `json:"emergencyRelation"`
	RegByID           string     `json:"regById"`
	Registrar         string     `json:"registrar"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         *time.Time `json:"updatedAt,omitempty"`
	Version           int        `json:"version,omitempty"`
}

// RegPatientReq represents the request body for registering a new patient.
// swagger:parameters regPatientReq
type RegPatientReq struct {
	// FullName is the full name of the patient.
	// required: true
	// min length: 2
	// max length: 100
	FullName string `json:"fullname" validate:"required,min=2,max=100"`

	// Gender represents the patient's gender.
	// required: true
	// allowed values: MALE, FEMALE, OTHER
	Gender string `json:"gender" validate:"required,oneof=MALE FEMALE OTHER"`

	// DOB is the date of birth of the patient.
	// required: true
	DOB time.Time `json:"dob" validate:"required"`

	// Age is the age of the patient.
	// required: true
	// numeric: true
	// max value: 100
	// format: YYYY-MM-DD
	Age int `json:"-" validate:"required,numeric,lte=100"`

	// ContactNumber is the patient's contact number.
	// required: true
	// numeric: true
	// length: 10 digits
	ContactNumber string `json:"contactNo" validate:"required,numeric,len=10"`

	// Address is the patient's address.
	// required: true
	// min length: 5
	// max length: 255
	Address string `json:"address" validate:"required,min=5,max=255"`

	// EmergencyName is the name of the emergency contact person.
	// required: true
	EmergencyName string `json:"emergencyName" validate:"required"`

	// EmergencyRelation is the relationship to the emergency contact person.
	// required: true
	EmergencyRelation string `json:"emergencyRelation" validate:"required"`

	// EmergencyPhone is the phone number of the emergency contact person.
	// required: true
	// numeric: true
	// length: 10 digits
	EmergencyPhone string `json:"emergencyPhone" validate:"required,numeric,len=10"`

	// RegByID is the ID of the user who registered the patient.
	// It's not included in the API payload.
	RegByID string `json:"-"`
}

type ListPatientRes struct {
	Patients []*ListPatientItem   `json:"patients"`
	Meta     *ListPatientMetadata `json:"meta"`
}

type ListPatientMetadata struct {
	CurrentPage int64 `json:"currentPage"` // e.g., 2
	PageSize    int64 `json:"pageSize"`    // e.g., 10
	TotalItems  int64 `json:"totalItems"`  // e.g., 43
	TotalPages  int64 `json:"totalPages"`  // e.g., 5
	From        int64 `json:"from"`        // e.g., 11
	To          int64 `json:"to"`          // e.g., 20
	HasNext     bool  `json:"hasNext"`     // true if next page exists
	HasPrevious bool  `json:"hasPrevious"` // true if previous page exists
}

// ListPatientItem represents a patient in the list response.
// swagger:response listPatientItem
type ListPatientItem struct {
	// ID is the unique identifier of the patient.
	ID string `json:"id"`

	// Username is the username of the patient.
	FullName string `json:"fullName"`

	// Gender is the patient's gender.
	Gender string `json:"gender"`

	// Age is the patient's age.
	Age int `json:"age"`

	// DOB is the patient's date of birth.
	DOB time.Time `json:"dob"`
}

// UpdatePatientReq represents the request body for updating patient details.
// swagger:parameters updatePatientReq
type UpdatePatientReq struct {
	// ID is the unique identifier of the patient to update.
	// It's excluded from the API payload.
	ID string `json:"-"`

	// FullName is the updated full name of the patient.
	// optional: true
	// min length: 2
	// max length: 100
	FullName *string `json:"fullname" validate:"omitempty,min=2,max=100"`

	// Gender is the updated gender of the patient.
	// optional: true
	// allowed values: MALE, FEMALE, OTHER
	Gender *string `json:"gender" validate:"omitempty,oneof=MALE FEMALE OTHER"`

	// DOB is the updated date of birth of the patient.
	// optional: true
	DOB *time.Time `json:"dob" validate:"omitempty"`

	// Age is the updated age of the patient.
	// optional: true
	// numeric: true
	// max value: 100
	Age *int `json:"age" validate:"omitempty,numeric,lte=100"`

	// ContactNumber is the updated contact number of the patient.
	// optional: true
	// numeric: true
	// length: 10 digits
	ContactNumber *string `json:"contactNo" validate:"omitempty,numeric,len=10"`

	// Address is the updated address of the patient.
	// optional: true
	// min length: 5
	// max length: 255
	Address *string `json:"address" validate:"omitempty,min=5,max=255"`

	// EmergencyName is the updated emergency contact name.
	// optional: true
	EmergencyName *string `json:"emergencyName" validate:"omitempty"`

	// EmergencyRelation is the updated relationship to the emergency contact.
	// optional: true
	EmergencyRelation *string `json:"emergencyRelation" validate:"omitempty"`

	// EmergencyPhone is the updated emergency contact phone number.
	// optional: true
	// numeric: true
	// length: 10 digits
	EmergencyPhone *string `json:"emergencyPhone" validate:"omitempty,numeric,len=10"`
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

type Record struct {
	Patient    Patient          `json:"patient"`
	Allergies  []AllergyModel   `json:"allergies"`
	Conditions []ConditionModel `json:"conditions"`
	Diagnoses  []DiagnosesModel `json:"diagnoses"`
	Vitals     VitalModel       `json:"vitals"`
}

type PatientModel struct {
	ID                string    `json:"id"`
	FullName          string    `json:"fullName"`
	Age               int       `json:"age"`
	Gender            string    `json:"gender"`
	DateOfBirth       time.Time `json:"dob"`
	ContactNumber     string    `json:"contactNo"`
	Address           string    `json:"address"`
	EmergencyName     string    `json:"emergencyName"`
	EmergencyRelation string    `json:"emergencyRelation"`
	EmergencyPhone    string    `json:"emergencyPhone"`
	RegisteredByID    string    `json:"regByID"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type AllergyModel struct {
	ID         string     `json:"id"`
	PatientID  string     `json:"patientID"`
	Name       string     `json:"name"`
	Reaction   string     `json:"reaction"`
	Severity   string     `json:"severity"`
	RecordedAt time.Time  `json:"recordedAt"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

type ConditionModel struct {
	ID        string     `json:"id"`
	PatientID string     `json:"patientID"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type DiagnosesModel struct {
	ID        string     `json:"id"`
	PatientID string     `json:"patientID"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type VitalModel struct {
	ID                     string     `json:"id,omitempty"`
	PatientID              string     `json:"patientID,omitempty"`
	HeightCm               *float64   `json:"height_cm,omitempty"`
	WeightKg               *float64   `json:"weight_kg,omitempty"`
	BMI                    *float64   `json:"bmi,omitempty"`
	TemperatureC           *float64   `json:"temperature_c,omitempty"`
	Pulse                  *int       `json:"pulse,omitempty"`
	RespiratoryRate        *int       `json:"respiratory_rate,omitempty"`
	BloodPressureSystolic  *int       `json:"blood_pressure_systolic,omitempty"`
	BloodPressureDiastolic *int       `json:"blood_pressure_diastolic,omitempty"`
	OxygenSaturation       *float64   `json:"oxygen_saturation,omitempty"`
	CreatedAt              *time.Time `json:"createdAt,omitempty"`
	UpdatedAt              *time.Time `json:"updatedAt,omitempty"`
}
