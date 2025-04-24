package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
)

// HandleRecordAllergy godoc
// @Summary Record a new allergy for a patient
// @Description Records a new allergy for the patient by providing allergy details.
// @Tags Allergy
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID" // Path parameter for patient ID
// @Param body body dto.RegAllergyReq true "Allergy Details" // Body parameter for allergy data
// @Success 200 {object} map[string]string {"message": "allergy recorded successfully"}
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /patients/{patientID}/allergies [post]
func (h *handler) HandleRecordAllergy(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	var req dto.RegAllergyReq
	if err := DecodeJSON(r, &req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	req.PatientID = pID

	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Allergy.Record(ctx, &req); err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("allergy recorded successfully")

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "allergy recorded successfully",
	})
}

// HandleUpdateAllergy godoc
// @Summary Update an existing allergy for a patient
// @Description Update the allergy details for a specific allergy by its ID.
// @Tags Allergy
// @Accept  json
// @Produce  json
// @Param allergyID path string true "Allergy ID" // Path parameter for allergy ID
// @Param body body dto.UpdateAllergyReq true "Updated Allergy Details" // Body parameter for updated allergy data
// @Success 200 {object} map[string]string {"message": "allergy updated successfully"}
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /allergies/{allergyID} [put]
func (h *handler) HandleUpdateAllergy(w http.ResponseWriter, r *http.Request) {
	aID := chi.URLParam(r, "allergyID")
	if err := h.validate.Var(aID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	var req dto.UpdateAllergyReq
	if err := DecodeJSON(r, &req); err != nil {
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

	if err := h.store.Allergy.Update(ctx, &req); err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("allergy updated successfully")

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "allergy updated successfully",
	})
}

// HandleDeleteAllergy godoc
// @Summary Delete an allergy from the patient's record
// @Description Delete a specific allergy by its ID.
// @Tags Allergy
// @Accept  json
// @Produce  json
// @Param allergyID path string true "Allergy ID" // Path parameter for allergy ID
// @Success 200 {object} map[string]string {"message": "allergy deleted successfully"}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /allergies/{allergyID} [delete]
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

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "allergy deleted successfully",
	})
}
