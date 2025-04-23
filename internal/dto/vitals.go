package dto

type CreateVitalReq struct {
	PatientID              string  `json:"patientId" validate:"required,uuid4"`
	HeightCm               float64 `json:"heightCm" validate:"gte=0"`
	WeightKg               float64 `json:"weightKg" validate:"gte=0"`
	BMI                    float64 `json:"bmi" validate:"gte=0"`
	TemperatureC           float64 `json:"temperatureC" validate:"gte=30,lte=45"`
	Pulse                  int     `json:"pulse" validate:"gte=0"`
	RespiratoryRate        int     `json:"respiratoryRate" validate:"gte=0"`
	BloodPressureSystolic  int     `json:"bloodPressureSystolic" validate:"gte=0"`
	BloodPressureDiastolic int     `json:"bloodPressureDiastolic" validate:"gte=0"`
	OxygenSaturation       float64 `json:"oxygenSaturation" validate:"gte=0,lte=100"`
}

type UpdateVitalReq struct {
	PatientID              string   `json:"-"` // hidden, used internally
	HeightCm               *float64 `json:"heightCm,omitempty" validate:"omitempty,gte=0"`
	WeightKg               *float64 `json:"weightKg,omitempty" validate:"omitempty,gte=0"`
	BMI                    *float64 `json:"bmi,omitempty" validate:"omitempty,gte=0"`
	TemperatureC           *float64 `json:"temperatureC,omitempty" validate:"omitempty,gte=30,lte=45"`
	Pulse                  *int     `json:"pulse,omitempty" validate:"omitempty,gte=0"`
	RespiratoryRate        *int     `json:"respiratoryRate,omitempty" validate:"omitempty,gte=0"`
	BloodPressureSystolic  *int     `json:"bloodPressureSystolic,omitempty" validate:"omitempty,gte=0"`
	BloodPressureDiastolic *int     `json:"bloodPressureDiastolic,omitempty" validate:"omitempty,gte=0"`
	OxygenSaturation       *float64 `json:"oxygenSaturation,omitempty" validate:"omitempty,gte=0,lte=100"`
}
