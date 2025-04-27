package models

import "time"

// CreateSessReq represents the request body for creating a new session.
// swagger:parameters createSessReq
type CreateSessReq struct {
	// UserID is the unique identifier for the user creating the session.
	// required: true
	// example: "user123"
	UserID string `json:"userID"`

	// Token is the authentication token for the session.
	// required: true
	// example: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJ1c2VyMTIzIn0.V9L8YrJshFkJcP-VEC9eFw"
	Token string `json:"token"`

	// Expiry is the expiration time of the session token.
	// required: true
	// example: "2025-12-31T23:59:59Z"
	Expiry time.Time `json:"expiry"`
}
