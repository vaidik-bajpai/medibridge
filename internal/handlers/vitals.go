package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

// HandleCaptureVitals godoc
// @Summary Capture patient vitals
// @Description Captures the vital information for a specific patient identified by their patient ID.
// @Tags Vitals
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID" // The unique identifier for the patient
// @Param body body dto.CreateVitalReq true "Vitals Information" // Body parameter for capturing patient vitals
// @Success 200 {object} map[string]string {"message": "vitals captured successfully"}
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /patients/{patientID}/vitals [post]
func (h *handler) HandleCaptureVitals(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		h.logger.Info("bad request", zap.Error(err))
		badRequestResponse(w, r)
		return
	}

	var req dto.CreateVitalReq
	if err := DecodeJSON(r, &req); err != nil {
		h.logger.Info("unprocessable entity", zap.Error(err))
		unprocessableEntityResponse(w, r)
		return
	}

	req.PatientID = pID

	if err := h.validate.Struct(req); err != nil {
		h.logger.Info("unprocessable entity", zap.Error(err))
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Vitals.Create(ctx, &req); err != nil {
		h.logger.Info("internal server error", zap.Error(err))
		if ok := errors.Is(err, store.ErrUniqueConstraintViolated); ok {
			conflictErrorResponse(w, r)
			return
		}
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("vitals captured successfully")

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "vitals captured successfully",
	})
}

// HandleUpdatingVitals godoc
// @Summary Update patient vitals
// @Description Updates the vital information for a specific patient identified by their patient ID.
// @Tags Vitals
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID" // The unique identifier for the patient
// @Param body body dto.UpdateVitalReq true "Vitals Information" // Body parameter for updating patient vitals
// @Success 200 {object} map[string]string {"message": "vitals updated successfully"}
// @Failure 400 {object} ErrorResponse
// @Failure 422 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /patients/{patientID}/vitals [put]
func (h *handler) HandleUpdatingVitals(w http.ResponseWriter, r *http.Request) {
	// 1) UUID v4 validation → 422 if malformed
	pID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(pID, "required,uuid4"); err != nil {
		log.Println("invalid UUID:", err)
		unprocessableEntityResponse(w, r)
		return
	}

	// 2) JSON decoding → 400 if malformed
	var req dto.UpdateVitalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("decode error:", err)
		badRequestResponse(w, r)
		return
	}

	// 3) Inject path param
	req.PatientID = pID

	// 4) DTO validation → 422 if any field fails
	if err := h.validate.Struct(req); err != nil {
		log.Println("validation error:", err)
		unprocessableEntityResponse(w, r)
		return
	}

	// 5) Persist → 500 on DB errors
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.store.Vitals.Update(ctx, &req); err != nil {
		log.Println("update error:", err)
		serverErrorResponse(w, r)
		return
	}

	// 6) Success → 200
	h.logger.Info("vitals updated successfully")
	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "vitals updated successfully",
	})
}

// HandleDeleteVitals godoc
// @Summary Delete patient vitals
// @Description Deletes the vital information for a specific patient identified by their patient ID.
// @Tags Vitals
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID" // The unique identifier for the patient
// @Success 200 {object} map[string]string {"message": "vitals deleted successfully"}
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /patients/{patientID}/vitals [delete]
func (h *handler) HandleDeleteVitals(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Vitals.Delete(ctx, pID); err != nil {
		log.Println(err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("vitals deleted successfully")

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "vitals deleted successfully",
	})
}
