package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/vaidik-bajpai/medibridge/internal/models"
	"go.uber.org/zap"
)

var (
	ErrPatientNotFound = errors.New("patient record not found")
)

// HandleRegisterPatient godoc
// @Summary Register a new patient
// @Description Registers a new patient with the provided details.
// @Tags Patients
// @Accept json
// @Produce json
// @Param body body models.RegPatientReq true "Patient Registration Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 422 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/patient [post]
func (h *handler) HandleRegisterPatient(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	var req models.RegPatientReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err.Error())
		unprocessableEntityResponse(w, r)
		return
	}

	req.Sanitize()
	req.RegByID = user.ID

	if err := h.validate.Struct(req); err != nil {
		log.Println("error: ", err)
		badRequestResponse(w, r)
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

// HandleUpdatePatientDetails godoc
// @Summary Update an existing patient
// @Description Updates details of an existing patient identified by patient ID.
// @Tags Patients
// @Accept json
// @Produce json
// @Param patientID path string true "Patient ID"
// @Param body body models.UpdatePatientReq true "Updated Patient Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 422 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/patient/{patientID} [put]
func (h *handler) HandleUpdatePatientDetails(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	h.logger.Info("identifier", zap.String("patient", patientID))
	if err := h.validate.Var(patientID, "required,uuid"); err != nil {
		unprocessableEntityResponse(w, r)
		return
	}

	var req models.UpdatePatientReq
	if err := DecodeJSON(r, &req); err != nil {
		badRequestResponse(w, r)
		return
	}

	req.Sanitize()
	req.ID = patientID

	if err := h.validate.Struct(req); err != nil {
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.store.Patient.Update(ctx, &req); err != nil {
		serverErrorResponse(w, r)
		return
	}

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "patient data updated successfully",
	})
}

// HandleDeletePatientDetails godoc
// @Summary Delete a patient
// @Description Deletes a patient by their patient ID.
// @Tags Patients
// @Accept json
// @Produce json
// @Param patientID path string true "Patient ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 404 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/patient/{patientID} [delete]
func (h *handler) HandleDeletePatientDetails(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(patientID, "required,uuid"); err != nil {
		log.Println(err)
		unprocessableEntityResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := h.store.Patient.Delete(ctx, patientID)
	if err != nil {
		log.Println(err)
		if ok := errors.Is(err, ErrPatientNotFound); ok {
			notFoundError(w, r)
			return
		}
		serverErrorResponse(w, r)
		return
	}

	WriteJSONResponse(w, r, http.StatusOK, map[string]string{
		"message": "patient deleted successfully",
	})
}

// HandleListPatients godoc
// @Summary List patients
// @Description Lists all patients with pagination (optional).
// @Tags Patients
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/patient [get]
func (h *handler) HandleListPatients(w http.ResponseWriter, r *http.Request) {
	paginate := getPaginateFromContext(r)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	list, err := h.store.Patient.List(ctx, paginate)
	if err != nil {
		log.Println(err)
		if ok := errors.Is(err, ErrPatientNotFound); ok {
			notFoundError(w, r)
			return
		}
		serverErrorResponse(w, r)
		return
	}

	WriteJSONResponse(w, r, http.StatusOK, map[string]interface{}{
		"list": list,
	})
}

// HandleGetPatient godoc
// @Summary Get patient by ID
// @Description Retrieves a patient's details using their patient ID.
// @Tags Patients
// @Accept json
// @Produce json
// @Param patientID path string true "Patient ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} models.FailureResponse
// @Failure 404 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/patient/{patientID} [get]
func (h *handler) HandleGetPatient(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	if err := h.validate.Var(patientID, "required,uuid"); err != nil {
		badRequestResponse(w, r)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	record, err := h.store.Patient.Get(ctx, patientID)
	if err != nil {
		log.Println(err)
		if ok := errors.Is(err, ErrPatientNotFound); ok {
			notFoundError(w, r)
			return
		}
		serverErrorResponse(w, r)
		return
	}

	WriteJSONResponse(w, r, http.StatusOK, map[string]interface{}{
		"record": record,
	})
}
