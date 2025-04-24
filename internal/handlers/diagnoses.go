package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/store"
)

// HandleAddDiagnoses godoc
// @Summary Add a new diagnosis for a patient
// @Description Adds a new diagnosis for a patient identified by their patient ID.
// @Tags Diagnoses
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID" // Path parameter for patient ID
// @Param body body dto.DiagnosesReq true "Diagnosis Details" // Body parameter for the diagnosis data
// @Success 200 {object} map[string]string {"message": "diagnoses added successfully"}
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /patients/{patientID}/diagnoses [post]
func (h *handler) HandleAddDiagnoses(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(pID, "required,uuid4"); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	var req dto.DiagnosesReq
	if err := DecodeJSON(r, &req); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.PID = pID

	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Diagnoses.Add(ctx, &req); err != nil {
		log.Println(err)
		if ok := errors.Is(err, store.ErrPatientNotFound); !ok {
			notFoundError(w, r)
			return
		}
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("diagnoses added successfully")

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "diagnoses added successfully",
	})
}

// HandleUpdateDiagnoses godoc
// @Summary Update an existing diagnosis for a patient
// @Description Updates a diagnosis based on the provided diagnosis ID and details.
// @Tags Diagnoses
// @Accept  json
// @Produce  json
// @Param diagnosesID path string true "Diagnosis ID" // Path parameter for diagnosis ID
// @Param body body dto.UpdateDiagnosesReq true "Updated Diagnosis Details" // Body parameter for the updated diagnosis data
// @Success 200 {object} map[string]string {"message": "diagnoses updated successfully"}
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /diagnoses/{diagnosesID} [put]
func (h *handler) HandleUpdateDiagnoses(w http.ResponseWriter, r *http.Request) {
	dID := chi.URLParam(r, "diagnosesID")
	if err := h.validate.Var(dID, "required,uuid"); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	var req dto.UpdateDiagnosesReq
	if err := DecodeJSON(r, &req); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.DID = dID

	fmt.Println(req.Name)

	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Diagnoses.Update(ctx, &req); err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("diagnoses updated successfully")

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "diagnoses updated successfully",
	})
}

// HandleDeleteDiagnoses godoc
// @Summary Delete a diagnosis for a patient
// @Description Deletes a diagnosis for a patient based on the provided diagnosis ID.
// @Tags Diagnoses
// @Accept  json
// @Produce  json
// @Param diagnosesID path string true "Diagnosis ID" // Path parameter for diagnosis ID
// @Success 200 {object} map[string]string {"message": "diagnoses deleted successfully"}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /diagnoses/{diagnosesID} [delete]
func (h *handler) HandleDeleteDiagnoses(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "diagnosesID")
	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Diagnoses.Delete(ctx, pID); err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("diagnoses deleted successfully")

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "diagnoses deleted successfully",
	})
}
