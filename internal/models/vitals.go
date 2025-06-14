package models

// CreateVitalReq represents the request body for creating a patient's vitals.
// @Description Request payload to capture new vital signs of a patient.
type CreateVitalReq struct {
	PatientID              string   `json:"-" swaggerignore:"true"`
	HeightCm               *float64 `json:"heightCm" validate:"omitempty,gt=0" example:"170"`
	WeightKg               *float64 `json:"weightKg" validate:"omitempty,gt=0" example:"65"`
	BMI                    *float64 `json:"bmi" validate:"omitempty,gt=0" example:"22.5"`
	TemperatureC           *float64 `json:"temperatureC" validate:"omitempty,gt=29,lte=45" example:"36.5"`
	Pulse                  *int     `json:"pulse" validate:"omitempty,gt=0" example:"72"`
	RespiratoryRate        *int     `json:"respiratoryRate" validate:"omitempty,gt=0" example:"18"`
	BloodPressureSystolic  *int     `json:"bloodPressureSystolic" validate:"omitempty,gt=0" example:"120"`
	BloodPressureDiastolic *int     `json:"bloodPressureDiastolic" validate:"omitempty,gt=0" example:"80"`
	OxygenSaturation       *float64 `json:"oxygenSaturation" validate:"omitempty,gt=0,lte=100" example:"98.0"`
}

// UpdateVitalReq represents the request body for updating a patient's vitals.
// @Description Request payload to update existing vital signs of a patient. All fields are optional.
type UpdateVitalReq struct {
	PatientID              string   `json:"-" swaggerignore:"true"`
	HeightCm               *float64 `json:"heightCm" validate:"omitempty,gte=0" example:"172"`
	WeightKg               *float64 `json:"weightKg" validate:"omitempty,gte=0" example:"68"`
	BMI                    *float64 `json:"bmi" validate:"omitempty,gte=0" example:"23.0"`
	TemperatureC           *float64 `json:"temperatureC" validate:"omitempty,gte=30,lte=45" example:"37.0"`
	Pulse                  *int     `json:"pulse" validate:"omitempty,gte=0" example:"75"`
	RespiratoryRate        *int     `json:"respiratoryRate" validate:"omitempty,gte=0" example:"20"`
	BloodPressureSystolic  *int     `json:"bloodPressureSystolic" validate:"omitempty,gte=0" example:"122"`
	BloodPressureDiastolic *int     `json:"bloodPressureDiastolic" validate:"omitempty,gte=0" example:"82"`
	OxygenSaturation       *float64 `json:"oxygenSaturation" validate:"omitempty,gte=0,lte=100" example:"97.0"`
}
