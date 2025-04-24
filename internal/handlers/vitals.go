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

func (h *handler) HandleUpdatingVitals(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	h.logger.Info("identifier", zap.String("patient", patientID))
	if err := h.validate.Var(patientID, "required,uuid"); err != nil {
		unprocessableEntityResponse(w, r)
		return
	}

	var req dto.UpdateVitalReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("decode error:", err)
		badRequestResponse(w, r)
		return
	}

	req.PatientID = patientID

	if err := h.validate.Struct(req); err != nil {
		log.Println("validation error:", err)
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	if err := h.store.Vitals.Update(ctx, &req); err != nil {
		log.Println("update error:", err)
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("vitals updated successfully")
	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "vitals updated successfully",
	})
}

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
