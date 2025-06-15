package models

import "time"

type Condition struct {
	ID        string    `json:"id"`
	PatientID string    `json:"patientID"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// AddConditionReq represents the request body for adding a condition.
// swagger:parameters addConditionReq
type AddConditionReq struct {
	// The condition to be added.
	// required: true
	// min length: 2
	// max length: 30
	Condition string `json:"condition" validate:"required,min=2,max=30"`

	// PatientID is excluded from the API payload.
	// It's not included in the JSON body.
	PatientID string `json:"-"`
}
