package dto

type CreateVitalReq struct {
	PatientID              string   `json:"-" validate:"required,uuid"` // PatientID as required and should be a valid UUID
	HeightCm               *float64 `json:"heightCm" validate:"omitempty,gt=0"`
	WeightKg               *float64 `json:"weightKg" validate:"omitempty,gt=0"`
	BMI                    *float64 `json:"bmi" validate:"omitempty,gt=0"`
	TemperatureC           *float64 `json:"temperatureC" validate:"omitempty,gt=29,lte=45"`
	Pulse                  *int     `json:"pulse" validate:"omitempty,gt=0"`
	RespiratoryRate        *int     `json:"respiratoryRate" validate:"omitempty,gt=0"`
	BloodPressureSystolic  *int     `json:"bloodPressureSystolic" validate:"omitempty,gt=0"`
	BloodPressureDiastolic *int     `json:"bloodPressureDiastolic" validate:"omitempty,gt=0"`
	OxygenSaturation       *float64 `json:"oxygenSaturation" validate:"omitempty,gt=0,lte=100"`
}

// UpdateVitalReq represents the request body for updating an existing set of vital signs.
// swagger:parameters updateVitalReq
type UpdateVitalReq struct {
	// PatientID is the unique identifier of the patient.
	// required: true
	// example: "patient123"
	PatientID string `json:"-"`

	// HeightCm is the patient's height in centimeters.
	// optional: true
	// example: 176.5
	HeightCm *float64 `json:"heightCm,omitempty" validate:"omitempty,gte=0"`

	// WeightKg is the patient's weight in kilograms.
	// optional: true
	// example: 71.0
	WeightKg *float64 `json:"weightKg,omitempty" validate:"omitempty,gte=0"`

	// BMI is the patient's Body Mass Index.
	// optional: true
	// example: 23.0
	BMI *float64 `json:"bmi,omitempty" validate:"omitempty,gte=0"`

	// TemperatureC is the patient's body temperature in Celsius.
	// optional: true
	// example: 37.0
	TemperatureC *float64 `json:"temperatureC,omitempty" validate:"omitempty,gte=30,lte=45"`

	// Pulse is the patient's heart rate in beats per minute.
	// optional: true
	// example: 85
	Pulse *int `json:"pulse,omitempty" validate:"omitempty,gte=0"`

	// RespiratoryRate is the number of breaths per minute.
	// optional: true
	// example: 18
	RespiratoryRate *int `json:"respiratoryRate,omitempty" validate:"omitempty,gte=0"`

	// BloodPressureSystolic is the systolic blood pressure value in mmHg.
	// optional: true
	// example: 125
	BloodPressureSystolic *int `json:"bloodPressureSystolic,omitempty" validate:"omitempty,gte=0"`

	// BloodPressureDiastolic is the diastolic blood pressure value in mmHg.
	// optional: true
	// example: 85
	BloodPressureDiastolic *int `json:"bloodPressureDiastolic,omitempty" validate:"omitempty,gte=0"`

	// OxygenSaturation is the oxygen saturation percentage in the blood.
	// optional: true
	// example: 97.0
	OxygenSaturation *float64 `json:"oxygenSaturation,omitempty" validate:"omitempty,gte=0,lte=100"`
}
