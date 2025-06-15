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
// @Summary      Add a new medical condition
// @Description  Adds a new medical condition associated with a patient ID.
// @Tags         Conditions
// @Accept       json
// @Produce      json
// @Param        patientID  path      string                  true  "Patient ID (UUID)"
// @Param        body       body      models.AddConditionReq  true  "Condition details"
// @Success      201        {object}  models.SuccessResponse
// @Failure      400        {object}  models.FailureResponse
// @Failure      422        {object}  models.FailureResponse
// @Failure      500        {object}  models.FailureResponse
// @Router       /v1/patient/{patientID}/condition [post]
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

	condition, err := h.store.Conditions.Add(ctx, &req)
	if err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("condition added successfully")
	helpers.WriteJSONResponse(w, r, http.StatusCreated, &models.SuccessResponse{
		Status:  http.StatusCreated,
		Message: "condition added successfully",
		Data:    condition,
	})
}

// HandleInactiveCondition godoc
// @Summary      Inactivate a medical condition
// @Description  Marks an existing condition as inactive by its ID.
// @Tags         Conditions
// @Accept       json
// @Produce      json
// @Param        conditionID  path      string  true  "Condition ID (UUID)"
// @Success      200          {object}  models.SuccessResponse
// @Failure      400          {object}  models.FailureResponse
// @Failure      500          {object}  models.FailureResponse
// @Router       /v1/condition/{conditionID} [delete]
func (h *handler) HandleInactiveCondition(w http.ResponseWriter, r *http.Request) {
	cID := chi.URLParam(r, "conditionID")
	if err := h.validate.Var(cID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Conditions.Delete(ctx, cID); err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("condition made inactive successfully")
	helpers.WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "condition made inactive successfully",
	})
}
