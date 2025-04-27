package models

// DiagnosesReq represents the request body for adding a diagnosis.
// swagger:parameters addDiagnosesReq
type DiagnosesReq struct {
	// PID is excluded from the API payload and not included in the JSON body.
	// It's a server-side value.
	PID string `json:"-"`

	// Name represents the name of the diagnosis.
	// required: true
	// min length: 2
	// max length: 30
	Name string `json:"name" validate:"required,min=2,max=30"`
}

// UpdateDiagnosesReq represents the request body for updating a diagnosis.
// swagger:parameters updateDiagnosesReq
type UpdateDiagnosesReq struct {
	// DID is excluded from the API payload and not included in the JSON body.
	// It's a server-side value used to identify the diagnosis being updated.
	DID string `json:"-"`

	// Name represents the updated name of the diagnosis.
	// required: true
	// min length: 2
	// max length: 30
	Name string `json:"name" validate:"required,min=2,max=30"`
}
