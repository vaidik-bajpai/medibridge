package models

type CreateVitalReq struct {
	PatientID              string   `json:"-"`
	HeightCm               *float64 `json:"heightCm" validate:"gt=0"`
	WeightKg               *float64 `json:"weightKg" validate:"gt=0"`
	BMI                    *float64 `json:"bmi" validate:"gt=0"`
	TemperatureC           *float64 `json:"temperatureC" validate:"gt=29,lte=45"`
	Pulse                  *int     `json:"pulse" validate:"gt=0"`
	RespiratoryRate        *int     `json:"respiratoryRate" validate:"gt=0"`
	BloodPressureSystolic  *int     `json:"bloodPressureSystolic" validate:"gt=0"`
	BloodPressureDiastolic *int     `json:"bloodPressureDiastolic" validate:"gt=0"`
	OxygenSaturation       *float64 `json:"oxygenSaturation" validate:"gt=0,lte=100"`
}

type UpdateVitalReq struct {
	PatientID              string   `json:"-"`
	HeightCm               *float64 `json:"heightCm" validate:"omitempty,gte=0"`
	WeightKg               *float64 `json:"weightKg" validate:"omitempty,gte=0"`
	BMI                    *float64 `json:"bmi" validate:"omitempty,gte=0"`
	TemperatureC           *float64 `json:"temperatureC" validate:"omitempty,gte=30,lte=45"`
	Pulse                  *int     `json:"pulse" validate:"omitempty,gte=0"`
	RespiratoryRate        *int     `json:"respiratoryRate" validate:"omitempty,gte=0"`
	BloodPressureSystolic  *int     `json:"bloodPressureSystolic" validate:"omitempty,gte=0"`
	BloodPressureDiastolic *int     `json:"bloodPressureDiastolic" validate:"omitempty,gte=0"`
	OxygenSaturation       *float64 `json:"oxygenSaturation" validate:"omitempty,gte=0,lte=100"`
}
