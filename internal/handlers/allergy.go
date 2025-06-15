package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vaidik-bajpai/medibridge/internal/helpers"
	"github.com/vaidik-bajpai/medibridge/internal/models"
)

// HandleRecordAllergy godoc
// @Summary Record a new allergy for a patient
// @Description Records a new allergy for the patient by providing allergy details.
// @Tags Allergy
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID"
// @Param body body models.RegAllergyReq true "Allergy Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 422 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/patient/{patientID}/allergy [post]
func (h *handler) HandleRecordAllergy(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "patientID")

	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	var req models.RegAllergyReq

	if err := helpers.DecodeJSON(r, &req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	req.PatientID = pID

	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	allergy, err := h.store.Allergy.Record(ctx, &req)
	if err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("allergy recorded successfully")

	helpers.WriteJSONResponse(w, r, http.StatusOK, &models.SuccessResponse{
		Status:  http.StatusCreated,
		Message: "allergy recorded successfully",
		Data:    allergy,
	})
}

// HandleUpdateAllergy godoc
// @Summary Update an existing allergy for a patient
// @Description Update the allergy details for a specific allergy by its ID.
// @Tags Allergy
// @Accept  json
// @Produce  json
// @Param allergyID path string true "Allergy ID"
// @Param body body models.UpdateAllergyReq true "Updated Allergy Details"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 422 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/allergy/{allergyID} [put]
func (h *handler) HandleUpdateAllergy(w http.ResponseWriter, r *http.Request) {
	aID := chi.URLParam(r, "allergyID")
	if err := h.validate.Var(aID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	var req models.UpdateAllergyReq

	if err := helpers.DecodeJSON(r, &req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	req.AllergyID = aID

	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	allergy, err := h.store.Allergy.Update(ctx, &req)
	if err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("allergy updated successfully")

	helpers.WriteJSONResponse(w, r, http.StatusOK, &models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "allergy updated successfully",
		Data:    allergy,
	})
}

// HandleDeleteAllergy godoc
// @Summary Delete an allergy from the patient's record
// @Description Delete a specific allergy by its ID.
// @Tags Allergy
// @Accept  json
// @Produce  json
// @Param allergyID path string true "Allergy ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/allergy/{allergyID} [delete]
func (h *handler) HandleDeleteAllergy(w http.ResponseWriter, r *http.Request) {
	aID := chi.URLParam(r, "allergyID")

	if err := h.validate.Var(aID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Allergy.Delete(ctx, aID); err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("allergy deleted successfully")

	helpers.WriteJSONResponse(w, r, http.StatusOK, &models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "allergy deleted successfully",
	})
}
