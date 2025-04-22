package handlers

import (
	"net/http"

	"github.com/vaidik-bajpai/medibridge/internal/store"
)

type userCtxKey string

const userCtx userCtxKey = "user"

func getUserFromCtx(r *http.Request) *store.User {
	user, _ := r.Context().Value(userCtx).(*store.User)
	return user
}
