package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"go.uber.org/zap"
)

var (
	ErrPatientNotFound = errors.New("patient record not found")
)

func (h *handler) HandleRegisterPatient(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)
	fmt.Println("user", *user)

	var req dto.RegPatientReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err.Error())
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "invalid request payload",
		})
		return
	}

	req.Sanitize()
	req.RegByID = user.ID

	if err := h.validate.Struct(req); err != nil {
		log.Println("error: ", err)
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, map[string]string{
			"error": "validation failed",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.store.Patient.Create(ctx, &req)
	if err != nil {
		log.Println("error: ", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "something went wrong with our server",
		})
		return
	}

	h.logger.Info(
		"patient registered successfully",
		zap.String("patient name", req.FullName),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{
		"message": "patient registered successfully",
	})
}

func (h *handler) HandleUpdatePatientDetails(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")

	if err := h.validate.Var(patientID, "required,uuid"); err != nil {
		unprocessableEntityResponse(w, r)
		return
	}

	var req dto.UpdatePatientReq
	if err := DecodeJSON(r, &req); err != nil {
		badRequestResponse(w, r)
		return
	}

	req.Sanitize()
	req.ID = patientID

	if err := h.validate.Struct(req); err != nil {
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Patient.Update(ctx, &req); err != nil {
		serverErrorResponse(w, r)
		return
	}

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "patient data updated successfully",
	})
}

func (h *handler) HandleDeletePatientDetails(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(patientID, "required,uuid"); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.store.Patient.Delete(ctx, patientID)
	if err != nil {
		log.Println(err)
		if ok := errors.Is(err, ErrPatientNotFound); ok {
			notFoundError(w, r)
			return
		}
		serverErrorResponse(w, r)
		return
	}

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "patient deleted successfully",
	})
}
