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

// HandleDeletePatientDetails godoc
// @Summary Delete a patient's details
// @Description Deletes a patient based on the provided patient ID.
// @Tags Patients
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} models.FailureResponse
// @Failure 404 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /v1/patient/{patientID} [delete]
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
// @Summary Update patient details
// @Description Updates the details of an existing patient.
// @Tags Patients
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID"
// @Param body body models.UpdatePatientReq true "Updated Patient Information"
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
// @Summary Delete a patient's details
// @Description Deletes a patient based on the provided patient ID.
// @Tags Patients
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID" // Path parameter for patient ID
// @Success 200 {object} map[string]string {"message": "patient deleted successfully"}
// @Failure 400 {object} models.FailureResponse
// @Failure 404 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /patients/{patientID} [delete]
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
// @Summary List all patients
// @Description Lists all patients with pagination support.
// @Tags Patients
// @Accept  json
// @Produce  json
// @Param page query int false "Page number" // Query parameter for page number
// @Param pageSize query int false "Page size" // Query parameter for page size
// @Success 200 {object} map[string]interface{} {"list": []models.Patient} // List of patients
// @Failure 400 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /patients [get]
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
// @Summary Get a patient by ID
// @Description Retrieves the details of a patient identified by their patient ID.
// @Tags Patients
// @Accept  json
// @Produce  json
// @Param patientID path string true "Patient ID" // Path parameter for patient ID
// @Success 200 {object} map[string]interface{} {"record": models.Patient} // Patient record
// @Failure 400 {object} models.FailureResponse
// @Failure 404 {object} models.FailureResponse
// @Failure 500 {object} models.FailureResponse
// @Router /patients/{patientID} [get]
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
