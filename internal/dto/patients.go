package dto

import (
	"strings"
	"time"
)

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
	Age int `json:"age" validate:"required,numeric,lte=100"`

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

// PatientListItem represents a patient in the list response.
// swagger:response patientListItem
type PatientListItem struct {
	// ID is the unique identifier of the patient.
	ID string `json:"id"`

	// Username is the username of the patient.
	Username string `json:"username"`

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
	Patient    PatientModel     `json:"patient"`
	Allergies  []AllergyModel   `json:"allergies"`
	Conditions []ConditionModel `json:"conditions"`
	Diagnoses  []DiagnosesModel `json:"diagnoses"`
	Vitals     VitalModel       `json:"vitals"`
}

type PatientModel struct {
	ID                string    `json:"id"`
	FullName          string    `json:"full_name"`
	Age               int       `json:"age"`
	Gender            string    `json:"gender"`
	DateOfBirth       time.Time `json:"date_of_birth"`
	ContactNumber     string    `json:"contact_number"`
	Address           string    `json:"address"`
	EmergencyName     string    `json:"emergency_name"`
	EmergencyRelation string    `json:"emergency_relation"`
	EmergencyPhone    string    `json:"emergency_phone"`
	RegisteredByID    string    `json:"registered_by_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type AllergyModel struct {
	ID         string     `json:"id"`
	PatientID  string     `json:"patient_id"`
	Name       string     `json:"name"`
	Reaction   string     `json:"reaction"`
	Severity   string     `json:"severity"`
	RecordedAt time.Time  `json:"recorded_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

type ConditionModel struct {
	ID        string     `json:"id"`
	PatientID string     `json:"patient_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type DiagnosesModel struct {
	ID        string     `json:"id"`
	PatientID string     `json:"patient_id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

type VitalModel struct {
	ID                     string    `json:"id"`
	PatientID              string    `json:"patient_id"`
	HeightCm               *float64  `json:"height_cm,omitempty"`
	WeightKg               *float64  `json:"weight_kg,omitempty"`
	BMI                    *float64  `json:"bmi,omitempty"`
	TemperatureC           *float64  `json:"temperature_c,omitempty"`
	Pulse                  *int      `json:"pulse,omitempty"`
	RespiratoryRate        *int      `json:"respiratory_rate,omitempty"`
	BloodPressureSystolic  *int      `json:"blood_pressure_systolic,omitempty"`
	BloodPressureDiastolic *int      `json:"blood_pressure_diastolic,omitempty"`
	OxygenSaturation       *float64  `json:"oxygen_saturation,omitempty"`
	CreatedAt              time.Time `json:"created_at"`
	UpdatedAt              time.Time `json:"updated_at"`
}
