package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vaidik-bajpai/medibridge/internal/helpers"
	"github.com/vaidik-bajpai/medibridge/internal/models"
)

// HandleAddDiagnoses godoc
// @Summary      Add a new diagnosis
// @Description  Adds a new diagnosis for a patient using their patient ID.
// @Tags         Diagnoses
// @Accept       json
// @Produce      json
// @Param        patientID  path      string                true  "Patient ID (UUID)"
// @Param        body       body      models.DiagnosesReq  true  "Diagnosis details"
// @Success      201        {object}  models.SuccessResponse
// @Failure      400        {object}  models.FailureResponse
// @Failure      422        {object}  models.FailureResponse
// @Failure      404        {object}  models.FailureResponse
// @Failure      500        {object}  models.FailureResponse
// @Router       /v1/patient/{patientID}/diagnoses [post]
func (h *handler) HandleAddDiagnoses(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	var req models.DiagnosesReq
	if err := helpers.DecodeJSON(r, &req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.PID = pID

	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	diag, err := h.store.Diagnoses.Add(ctx, &req)
	if err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("diagnoses added successfully")

	helpers.WriteJSONResponse(w, r, http.StatusCreated, &models.SuccessResponse{
		Status:  http.StatusCreated,
		Message: "diagnoses added successfully",
		Data:    diag,
	})
}

// HandleUpdateDiagnoses godoc
// @Summary      Update an existing diagnosis
// @Description  Updates an existing diagnosis using the diagnosis ID.
// @Tags         Diagnoses
// @Accept       json
// @Produce      json
// @Param        diagnosesID  path      string                     true  "Diagnosis ID (UUID)"
// @Param        body         body      models.UpdateDiagnosesReq true  "Updated diagnosis details"
// @Success      200          {object}  models.SuccessResponse
// @Failure      400          {object}  models.FailureResponse
// @Failure      422          {object}  models.FailureResponse
// @Failure      500          {object}  models.FailureResponse
// @Router       /v1/diagnoses/{diagnosesID} [put]
func (h *handler) HandleUpdateDiagnoses(w http.ResponseWriter, r *http.Request) {
	dID := chi.URLParam(r, "diagnosesID")
	if err := h.validate.Var(dID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	var req models.UpdateDiagnosesReq
	if err := helpers.DecodeJSON(r, &req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.DID = dID

	fmt.Println(req.Name)

	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	diag, err := h.store.Diagnoses.Update(ctx, &req)
	if err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("diagnoses updated successfully")

	helpers.WriteJSONResponse(w, r, http.StatusOK, &models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "diagnoses updated successfully",
		Data:    diag,
	})
}

// HandleDeleteDiagnoses godoc
// @Summary      Delete a diagnosis
// @Description  Deletes a diagnosis using the diagnosis ID.
// @Tags         Diagnoses
// @Accept       json
// @Produce      json
// @Param        diagnosesID  path      string  true  "Diagnosis ID (UUID)"
// @Success      200          {object}  models.SuccessResponse
// @Failure      400          {object}  models.FailureResponse
// @Failure      500          {object}  models.FailureResponse
// @Router       /v1/diagnoses/{diagnosesID} [delete]
func (h *handler) HandleDeleteDiagnoses(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "diagnosesID")
	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
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

	helpers.WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "diagnoses deleted successfully",
	})
}
