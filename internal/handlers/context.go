package handlers

import (
	"net/http"

	"github.com/vaidik-bajpai/medibridge/internal/models"
)

type userKey string

const userCtx userKey = "user"

func getUserFromCtx(r *http.Request) *models.UserModel {
	user, _ := r.Context().Value(userCtx).(*models.UserModel)
	return user
}

type paginateKey string

const paginateCtx paginateKey = "paginate"

func getPaginateFromContext(r *http.Request) *models.Paginate {
	paginate, _ := r.Context().Value(userCtx).(*models.Paginate)
	return paginate
}
