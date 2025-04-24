package dto

// SignupReq represents the request body for user signup.
// swagger:parameters signupReq
type SignupReq struct {
	// Fullname is the full name of the user.
	// required: true
	// min length: 3
	// max length: 30
	Fullname string `json:"fullname" validate:"required,min=3,max=30"`

	// Email is the email address of the user.
	// required: true
	// format: email
	// example: "user@example.com"
	Email string `json:"email" validate:"required,email"`

	// Password is the password chosen by the user.
	// required: true
	// min length: 8
	// max length: 64
	Password string `json:"password" validate:"required,min=8,max=64"`

	// Role defines the user's role in the system.
	// required: true
	// allowed values: doctor, receptionist
	Role string `json:"role" validate:"required,oneof=doctor receptionist"`

	// Activated is a flag that determines whether the user is active.
	// optional: true
	Activated bool `json:"activated"`
}

// SigninReq represents the request body for user sign-in.
// swagger:parameters signinReq
type SigninReq struct {
	// Email is the email address of the user.
	// required: true
	// format: email
	// example: "user@example.com"
	Email string `json:"email" validate:"required,email"`

	// Password is the password chosen by the user.
	// required: true
	// min length: 8
	// max length: 64
	Password string `json:"password" validate:"required,min=8,max=64"`
}

// UserModel represents a user in the system.
// swagger:response userModel
type UserModel struct {
	// ID is the unique identifier of the user.
	ID string `json:"id"`

	// Username is the username of the user.
	Username string `json:"username"`

	// Email is the email address of the user.
	Email string `json:"email"`

	// Password is the user's hashed password.
	Password string `json:"password"`

	// Activated is a flag that indicates whether the user is active.
	Activated bool `json:"activated"`

	// Role is the user's role within the system.
	Role string `json:"role"`

	// OAuthProvider is the external OAuth provider (e.g., Google, Facebook) for the user.
	OAuthProvider string `json:"oauthProvider"`

	// OAuthID is the ID provided by the external OAuth provider.
	OAuthID string `json:"oauthID"`
}
