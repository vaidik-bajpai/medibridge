package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	dto "github.com/vaidik-bajpai/medibridge/internal/models"
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
// @Failure 400 {object} ErrorResponse // Invalid input or bad request
// @Failure 422 {object} ErrorResponse // Unprocessable entity (validation errors)
// @Failure 500 {object} ErrorResponse // Internal server error
// @Router /patients/{patientID}/allergies [post]
func (h *handler) HandleRecordAllergy(w http.ResponseWriter, r *http.Request) {
	// Get the patient ID from the URL
	pID := chi.URLParam(r, "patientID")
	// Validate patient ID (it should be a UUID)
	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r) // Invalid request if patient ID is not valid
		return
	}

	var req dto.RegAllergyReq
	// Decode the request body to get allergy details
	if err := DecodeJSON(r, &req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r) // Invalid or malformed JSON body
		return
	}

	// Attach patient ID to the allergy details
	req.PatientID = pID

	// Validate the allergy request data
	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r) // Validation error
		return
	}

	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to record the allergy in the database
	if err := h.store.Allergy.Record(ctx, &req); err != nil {
		log.Println(err)
		serverErrorResponse(w, r) // Database or server error
		return
	}

	// Log successful allergy recording
	h.logger.Info("allergy recorded successfully")

	// Send a successful response to the client
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
// @Failure 400 {object} ErrorResponse // Invalid input or bad request
// @Failure 422 {object} ErrorResponse // Unprocessable entity (validation errors)
// @Failure 500 {object} ErrorResponse // Internal server error
// @Router /allergies/{allergyID} [put]
func (h *handler) HandleUpdateAllergy(w http.ResponseWriter, r *http.Request) {
	// Get the allergy ID from the URL
	aID := chi.URLParam(r, "allergyID")
	// Validate allergy ID (it should be a UUID)
	if err := h.validate.Var(aID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r) // Invalid request if allergy ID is not valid
		return
	}

	var req dto.UpdateAllergyReq
	// Decode the request body to get updated allergy details
	if err := DecodeJSON(r, &req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r) // Invalid or malformed JSON body
		return
	}

	// Attach allergy ID to the updated allergy details
	req.AllergyID = aID

	// Validate the updated allergy request data
	if err := h.validate.Struct(req); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r) // Validation error
		return
	}

	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to update the allergy in the database
	if err := h.store.Allergy.Update(ctx, &req); err != nil {
		log.Println(err)
		serverErrorResponse(w, r) // Database or server error
		return
	}

	// Log successful allergy update
	h.logger.Info("allergy updated successfully")

	// Send a successful response to the client
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
// @Failure 400 {object} ErrorResponse // Invalid input or bad request
// @Failure 500 {object} ErrorResponse // Internal server error
// @Router /allergies/{allergyID} [delete]
func (h *handler) HandleDeleteAllergy(w http.ResponseWriter, r *http.Request) {
	// Get the allergy ID from the URL
	aID := chi.URLParam(r, "allergyID")
	// Validate allergy ID (it should be a UUID)
	if err := h.validate.Var(aID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r) // Invalid request if allergy ID is not valid
		return
	}

	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Attempt to delete the allergy from the database
	if err := h.store.Allergy.Delete(ctx, aID); err != nil {
		log.Println(err)
		serverErrorResponse(w, r) // Database or server error
		return
	}

	// Log successful allergy deletion
	h.logger.Info("allergy deleted successfully")

	// Send a successful response to the client
	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "allergy deleted successfully",
	})
}
