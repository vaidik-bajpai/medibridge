package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/render"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

func (h *handler) HandleUserSignup(w http.ResponseWriter, r *http.Request) {
	var req dto.SignupReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error: ", err)
		badRequestResponse(w, r)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Role = strings.TrimSpace(req.Role)
	req.Fullname = strings.TrimSpace(req.Fullname)

	if err := h.validate.Struct(&req); err != nil {
		log.Println("error: ", err)
		unprocessableEntityResponse(w, r)
		return
	}

	hash, err := MakeHashFromToken(req.Password)
	if err != nil {
		log.Println("error: ", err)
		unprocessableEntityResponse(w, r)
		return
	}
	req.Password = string(hash)
	req.Activated = false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.store.User.Create(ctx, &req)
	if err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info(
		"user registered successfully",
		zap.String("username", req.Fullname),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{
		"message": "user registered successfully",
	})
}

func (h *handler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var req *dto.SigninReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user, err := h.store.User.FindViaEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			notFoundError(w, r)
			return
		}
		serverErrorResponse(w, r)
		return
	}

	ok, err := MatchPassword(user.Password, req.Password)
	if err != nil {
		log.Println("error: ", err)
		serverErrorResponse(w, r)
		return
	}
	if !ok {
		log.Println("error: ", err)
		unauthorisedErrorResponse(w, r, "invalid email or password")
		return
	}

	var cs dto.CreateSessReq
	cs.Token, err = GenerateSessionToken()
	if err != nil {
		log.Println("error: ", err)
		serverErrorResponse(w, r)
		return
	}
	cs.UserID = user.ID
	cs.Expiry = time.Now().Add(7 * 24 * time.Hour)

	err = h.store.Session.Create(ctx, &cs)
	if err != nil {
		log.Println("error: ", err)
		serverErrorResponse(w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "medibridge-token",
		Value:    cs.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, //set to true in prod
		SameSite: http.SameSiteLaxMode,
		Expires:  cs.Expiry,
	})

	h.logger.Info(
		"user login successful",
		zap.String("user id", cs.UserID),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, "user login successful")
}
