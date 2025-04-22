package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"go.uber.org/zap"
)

func (h *handler) HandleUserSignup(w http.ResponseWriter, r *http.Request) {
	var req dto.SignupReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error: ", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "invalid request payload",
		})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Role = strings.TrimSpace(req.Role)
	req.Username = strings.TrimSpace(req.Username)

	if err := h.validate.Struct(&req); err != nil {
		log.Println("error: ", err)
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, map[string]string{
			"error": "validation failed",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.store.User.Create(ctx, &req)
	if err != nil {
		log.Println(err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{"error": "something went wrong in our servers."})
		return
	}

	h.logger.Info(
		"user registered successfully",
		zap.String("username", req.Username),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{
		"message": "user created successfully",
	})
}
