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
	"github.com/vaidik-bajpai/medibridge/internal/helpers"
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
// @Router /v1/user/signup [post]
func (h *handler) HandleUserSignup(w http.ResponseWriter, r *http.Request) {
	var req models.SignupReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error decoding request body: ", err)
		unprocessableEntityResponse(w, r)
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	req.Role = strings.TrimSpace(req.Role)
	req.Fullname = strings.TrimSpace(req.Fullname)

	if err := h.validate.Struct(&req); err != nil {
		badRequestResponse(w, r)
		return
	}

	hash, err := helpers.MakeHashFromToken(req.Password)
	if err != nil {
		badRequestResponse(w, r)
		return
	}
	req.Password = string(hash)
	req.Activated = false

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = h.store.User.Create(ctx, &req)
	if err != nil {
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("user registered successfully", zap.String("username", req.Fullname))

	status := http.StatusCreated
	render.Status(r, status)
	render.JSON(w, r, models.SuccessResponse{
		Status:  status,
		Message: "user registered successfully",
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
// @Router /v1/user/signin [post]
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
		serverErrorResponse(w, r)
		return
	}

	ok, err := helpers.MatchPassword(user.Password, req.Password)
	if err != nil {
		serverErrorResponse(w, r)
		return
	}
	if !ok {
		unauthorisedErrorResponse(w, r, "invalid email or password")
		return
	}

	var cs models.CreateSessReq
	cs.Token, err = helpers.GenerateSessionToken()
	if err != nil {
		serverErrorResponse(w, r)
		return
	}
	cs.UserID = user.ID
	cs.Expiry = time.Now().Add(7 * 24 * time.Hour)

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
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Expires:  cs.Expiry,
	})

	h.logger.Info("user login successful", zap.String("user id", cs.UserID))

	status := http.StatusOK
	render.Status(r, status)
	render.JSON(w, r, models.SuccessResponse{
		Status:  status,
		Message: "user login successful",
	})
}

// HandleUserLogout godoc
// @Summary Logs out a user
// @Description Clears the user's session cookie
// @Tags Users
// @Produce  json
// @Success 200 {object} models.SuccessResponse
// @Router /v1/user/logout [post]
func (h *handler) HandleUserLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "medibridge-token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	h.logger.Info("user logout successful")

	render.Status(r, http.StatusOK)
	render.JSON(w, r, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "user logged out successfully",
	})
}
