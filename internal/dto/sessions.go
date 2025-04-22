package dto

import "time"

type CreateSessReq struct {
	UserID string
	Token  string
	Expiry time.Time
}
