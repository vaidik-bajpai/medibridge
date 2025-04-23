package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"go.uber.org/zap"
)

func (h *handler) HandleRegisterPatient(w http.ResponseWriter, r *http.Request) {
	var req dto.RegPatientReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("error: ", err)
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{
			"error": "invalid request payload",
		})
		return
	}

	req.Sanitize()

	if err := h.validate.Struct(req); err != nil {
		log.Println("error: ", err)
		render.Status(r, http.StatusUnprocessableEntity)
		render.JSON(w, r, map[string]string{
			"error": "validation failed",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.store.Patient.Create(ctx, &req)
	if err != nil {
		log.Println("error: ", err)
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, map[string]string{
			"error": "something went wrong with our server",
		})
		return
	}

	h.logger.Info(
		"patient registered successfully",
		zap.String("patient name", req.FullName),
	)

	render.Status(r, http.StatusOK)
	render.JSON(w, r, map[string]string{
		"message": "patient registered successfully",
	})
}

func (h *handler) HandleUpdateBasicDetails(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
func (h *handler) HandleAddDiagnoses(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
func (h *handler) HandleAddConditions(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
func (h *handler) HandleAddAllergies(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
func (h *handler) HandleAddVitals(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
