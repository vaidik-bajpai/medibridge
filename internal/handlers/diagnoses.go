package handlers

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
)

func (h *handler) HandleAddDiagnoses(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(pID, "required,uuid"); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	var req dto.DiagnosesReq
	if err := DecodeJSON(r, req); err != nil {
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
		serverErrorResponse(w, r)
		return
	}

	h.logger.Info("diagnoses added successfully")

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "diagnoses added successfully",
	})
}

func (h *handler) HandleUpdateDiagnoses(w http.ResponseWriter, r *http.Request) {
	dID := chi.URLParam(r, "diagnosesID")
	if err := h.validate.Var(dID, "required,uuid"); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	var req dto.UpdateDiagnosesReq
	if err := DecodeJSON(r, req); err != nil {
		log.Println(err)
		badRequestResponse(w, r)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.DID = dID

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

func (h *handler) HandleDeleteDiagnoses(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "patientID")
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
