package dto

type AddConditionReq struct {
	Condition string `json:"condition" validate:"required,min=2,max=30"`
	PatientID string `json:"-"`
}
