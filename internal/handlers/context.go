package handlers

import (
	"net/http"

	"github.com/vaidik-bajpai/medibridge/internal/dto"
)

type userKey string

const userCtx userKey = "user"

func getUserFromCtx(r *http.Request) *dto.UserModel {
	user, _ := r.Context().Value(userCtx).(*dto.UserModel)
	return user
}
