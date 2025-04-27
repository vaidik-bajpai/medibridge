package models

// RegAllergyReq represents the request body for registering a new allergy
// @Description A request to register a new allergy for a patient
// @Tags allergies
// @Accept json
// @Produce json
type RegAllergyReq struct {
	// @Param patient_id query string true "Patient ID" validate:"required,uuid4"
	// @example "a4b6d789-60e0-4c93-876d-7a8b62e3b883"
	PatientID string `json:"-" validate:"required,uuid4"`

	// @Param name query string true "Allergy Name" validate:"required,min=2,max=100"
	// @example "Peanut"
	Name string `json:"name" validate:"required,min=2,max=100"`

	// @Param severity query string true "Severity of the allergy" validate:"required,oneof=mild moderate severe"
	// @example "mild"
	Severity string `json:"severity" validate:"required,oneof=mild moderate severe"`

	// @Param reaction query string true "Reaction to the allergy" validate:"required,min=2,max=255"
	// @example "Swelling"
	Reaction string `json:"reaction" validate:"required,min=2,max=255"`
}

// UpdateAllergyReq represents the request body for updating an allergy record
// @Description A request to update an existing allergy record
// @Tags allergies
// @Accept json
// @Produce json
type UpdateAllergyReq struct {
	// @example "1b33b7a0-5e83-4b67-b3b1-b249d3d4acfe"
	// @Param allergy_id query string true "Allergy ID"
	AllergyID string `json:"-"`

	// @example "Peanut"
	// @Param name query string false "Updated name of the allergy" validate:"omitempty,min=2,max=100"
	Name *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`

	// @example "moderate"
	// @Param severity query string false "Updated severity of the allergy" validate:"omitempty,oneof=mild moderate severe"
	Severity *string `json:"severity,omitempty" validate:"omitempty,oneof=mild moderate severe"`

	// @example "swelling"
	// @Param reaction query string false "Updated reaction to the allergy" validate:"omitempty,min=2,max=255"
	Reaction *string `json:"reaction,omitempty" validate:"omitempty,min=2,max=255"`
}
