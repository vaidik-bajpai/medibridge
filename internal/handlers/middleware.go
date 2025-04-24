package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

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
