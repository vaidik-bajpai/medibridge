package handlers

import (
	"net/http"

	dto "github.com/vaidik-bajpai/medibridge/internal/models"
)

type userKey string

const userCtx userKey = "user"

func getUserFromCtx(r *http.Request) *dto.UserModel {
	user, _ := r.Context().Value(userCtx).(*dto.UserModel)
	return user
}

type paginateKey string

const paginateCtx paginateKey = "paginate"

func getPaginateFromContext(r *http.Request) *dto.Paginate {
	paginate, _ := r.Context().Value(userCtx).(*dto.Paginate)
	return paginate
}
