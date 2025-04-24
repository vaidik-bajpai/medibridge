package dto

type DiagnosesReq struct {
	PID  string `json:"-"`
	Name string `json:"name" required:"required,min=2,max=30"`
}

type UpdateDiagnosesReq struct {
	DID  string
	Name string `json:"name" required:"required,min=2,max=30"`
}
