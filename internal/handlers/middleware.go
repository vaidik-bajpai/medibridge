package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
	"go.uber.org/zap"
)

func (h *handler) RequirePaginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		lastID := query.Get("lastID")
		pageSize := query.Get("pageSize")

		var err error
		var paginate dto.Paginate
		paginate.LastID = lastID
		paginate.PageSize, err = strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			badRequestResponse(w, r)
			return
		}

		if err := h.validate.Struct(paginate); err != nil {
			badRequestResponse(w, r)
			return
		}

		pCtx := context.WithValue(r.Context(), paginateCtx, &paginate)
		next.ServeHTTP(w, r.WithContext(pCtx))
	})
}

func (h *handler) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("medibridge-token")
		if err != nil {
			h.logger.Error("unauthorised", zap.Error(err))
			unauthorisedErrorResponse(w, r, "you are unauthorized")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		user, err := h.store.Session.FindUserByToken(ctx, cookie.Value)
		if err != nil {
			h.logger.Error("unauthorised", zap.Error(err))
			unauthorisedErrorResponse(w, r, "you are unauthorised")
			return
		}
		fmt.Println(user)
		fmt.Println(cookie.Value)

		uCtx := context.WithValue(r.Context(), userCtx, user)
		next.ServeHTTP(w, r.WithContext(uCtx))
	})
}

func (h *handler) RequireRole(roles ...db.Role) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user := getUserFromCtx(r)

			for _, role := range roles {
				if user.Role == string(role) {
					next.ServeHTTP(w, r)
					return
				}
			}

			h.logger.Warn("user not permitted", zap.String("role", user.Role))
			forbiddenErrorResponse(w, r)
		})
	}
}
