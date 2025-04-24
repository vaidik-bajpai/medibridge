package dto

type RegAllergyReq struct {
	PatientID string `json:"-" validate:"required,uuid4"`
	Name      string `json:"name" validate:"required,min=2,max=100"`
	Severity  string `json:"severity" validate:"required,oneof=mild moderate severe"`
	Reaction  string `json:"reaction" validate:"required,min=2,max=255"`
}

type UpdateAllergyReq struct {
	AllergyID string  `json:"-"`
	Name      *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Severity  *string `json:"severity,omitempty" validate:"omitempty,oneof=mild moderate severe"`
	Reaction  *string `json:"reaction,omitempty" validate:"omitempty,min=2,max=255"`
}
