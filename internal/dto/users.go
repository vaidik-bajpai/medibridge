package dto

type SignupReq struct {
	Fullname  string `json:"fullname" validate:"required,min=3,max=30"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=64"`
	Role      string `json:"role" validate:"required,oneof=doctor receptionist"`
	Activated bool
}

type SigninReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type UserModel struct {
	ID            string
	Username      string
	Email         string
	Password      string
	Activated     bool
	Role          string
	OAuthProvider string
	OAuthID       string
}
