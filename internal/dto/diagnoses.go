package dto

type DiagnosesReq struct {
	PID  string `json:"-"`
	Name string `json:"name" validate:"required,min=2,max=30"`
}

type UpdateDiagnosesReq struct {
	DID  string `json:"-"`
	Name string `json:"name" validate:"required,min=2,max=30"`
}
