package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vaidik-bajpai/medibridge/internal/helpers"
	"github.com/vaidik-bajpai/medibridge/internal/models"
)

// HandleAddCondition godoc
// @Summary Add a new medical condition for a patient
// @Description Adds a new medical condition for a patient. The condition is associated with the patient's ID.
// @Tags Conditions
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID"
// @Param body body models.AddConditionReq true "Condition Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 422 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/patient/{patientID}/condition [post]
func (h *handler) HandleAddCondition(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	var req models.AddConditionReq
	if err := helpers.DecodeJSON(r, &req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	req.Condition = strings.TrimSpace(req.Condition)
	req.PatientID = pID

	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Conditions.Add(ctx, &req); err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("condition added successfully")
	helpers.WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "condition added successfully",
	})
}

// HandleInactiveCondition godoc
// @Summary Inactivate a patient's medical condition
// @Description Marks a condition as inactive for a patient by removing it from their active conditions.
// @Tags Conditions
// @Accept  json
// @Produce  json
// @Param conditionID path string true "Condtion ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/condition/{conditionID} [delete]
func (h *handler) HandleInactiveCondition(w http.ResponseWriter, r *http.Request) {
	// Extract condition ID from the URL parameter and validate it
	cID := chi.URLParam(r, "conditionID")
	if err := h.validate.Var(cID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r) // Return a bad request response if the validation fails
		return
	}

	// Set a timeout context for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to mark the condition as inactive by deleting it from the active conditions list
	if err := h.store.Conditions.Delete(ctx, cID); err != nil {
		log.Println(err)
		serverErrorResponse(w, r) // Return a server error response if the deletion fails
		return
	}

	// Log the successful inactivation of the condition and send the success response
	h.logger.Info("condition made inactive successfully")
	helpers.WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "condition made inactive successfully",
	})
}
