package handlers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
)

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
