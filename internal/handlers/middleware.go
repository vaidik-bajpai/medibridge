package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	dto "github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
	"go.uber.org/zap"
)

func (h *handler) RequirePaginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		page := query.Get("page")
		pageSize := query.Get("pageSize")
		searchTerm := query.Get("searchTerm")

		var err error
		var paginate dto.Paginate
		paginate.Page, err = strconv.ParseInt(page, 10, 64)
		if err != nil {
			badRequestResponse(w, r)
			return
		}

		paginate.PageSize, err = strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			badRequestResponse(w, r)
			return
		}

		paginate.SearchTerm = searchTerm

		if err := h.validate.Struct(paginate); err != nil {
			badRequestResponse(w, r)
			return
		}

		h.logger.Debug("parsed pagination data",
			zap.Int64("page", paginate.Page),
			zap.Int64("pagesize", paginate.PageSize),
			zap.String("search term", paginate.SearchTerm))

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
