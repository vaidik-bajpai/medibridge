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
	"github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

// HandleUserSignup godoc
// @Summary User signup
// @Description Registers a new user with the provided information including email, password, and role.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param body body models.SignupReq true "User Signup Information"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 422 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /users/signup [post]
func (h *handler) HandleUserSignup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding request body: ", err)
		badRequestResponse(w, r)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Role = strings.TrimSpace(req.Role)
	req.Fullname = strings.TrimSpace(req.Fullname)

	if err := h.validate.Struct(&req); err != nil {
		log.Println("error validating request: ", err)
		unprocessableEntityResponse(w, r)
		return
	}

	hash, err := MakeHashFromToken(req.Password)
	if err != nil {
		log.Println("error hashing password: ", err)
		unprocessableEntityResponse(w, r)
		return
	}
	req.Password = string(hash)
	req.Activated = false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.store.User.Create(ctx, &req)
	if err != nil {
		log.Println("error creating user: ", err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("user registered successfully", zap.String("username", req.Fullname))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{
		"message": "user registered successfully",
	})
}

// HandleUserLogin godoc
// @Summary User login
// @Description Logs in a user by verifying their email and password, and returns a session token if valid.
// @Tags Users
// @Accept  json
// @Produce  json
// @Param body body models.SigninReq true "User Login Information"
// @Success 200 {string} string "user login successful"
// @Failure 400 {object} models.FailureResponse
// @Failure 422 {object} models.FailureResponse
// @Failure 401 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /users/signin [post]
func (h *handler) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var req *models.SigninReq

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding request body: ", err)
		badRequestResponse(w, r)
		return
	}

	req.Email = strings.TrimSpace(req.Email)

	if err := h.validate.Struct(req); err != nil {
		log.Println("error validating login request: ", err)
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
		log.Println("error finding user: ", err)
		serverErrorResponse(w, r)
		return
	}

	ok, err := MatchPassword(user.Password, req.Password)
	if err != nil {
		log.Println("error matching password: ", err)
		serverErrorResponse(w, r)
		return
	}
	if !ok {
		log.Println("error: invalid email or password")
		unauthorisedErrorResponse(w, r, "invalid email or password")
		return
	}

	var cs models.CreateSessReq
	cs.Token, err = GenerateSessionToken()
	if err != nil {
		log.Println("error generating session token: ", err)
		serverErrorResponse(w, r)
		return
	}
	cs.UserID = user.ID
	cs.Expiry = time.Now().Add(7 * 24 * time.Hour) // Token valid for 7 days

	err = h.store.Session.Create(ctx, &cs)
	if err != nil {
		log.Println("error creating session: ", err)
		serverErrorResponse(w, r)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "medibridge-token",
		Value:    cs.Token,
		Path:     "/",
		HttpOnly: true,
		Secure:   false, // IMPORTANT: Should be true in production
		SameSite: http.SameSiteLaxMode,
		Expires:  cs.Expiry,
	})

	h.logger.Info("user login successful", zap.String("user id", cs.UserID))

	render.Status(r, http.StatusOK)
	render.JSON(w, r, "user login successful")
}
