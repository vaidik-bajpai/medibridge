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
	"github.com/vaidik-bajpai/medibridge/internal/helpers"
	"github.com/vaidik-bajpai/medibridge/internal/models"
	"go.uber.org/zap"
)

var (
	ErrPatientNotFound = errors.New("patient record not found")
)

// HandleRegisterPatient godoc
// @Summary      Register a new patient
// @Description  Registers a new patient with the provided details.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        body  body      models.RegPatientReq  true  "Patient registration data"
// @Success      201   {object}  models.SuccessResponse
// @Failure      400   {object}  models.FailureResponse
// @Failure      422   {object}  models.FailureResponse
// @Failure      500   {object}  models.FailureResponse
// @Router       /v1/patient [post]
func (h *handler) HandleRegisterPatient(w http.ResponseWriter, r *http.Request) {
	user := getUserFromCtx(r)

	var req models.RegPatientReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error(err.Error())
		unprocessableEntityResponse(w, r)
		return
	}

	req.Age = helpers.CalculateAge(time.Time(req.DOB))

	req.Sanitize()
	req.RegByID = user.ID

	if err := h.validate.Struct(req); err != nil {
		badRequestResponse(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p, err := h.store.Patient.Create(ctx, &req)
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, models.FailureResponse{
			Status: http.StatusInternalServerError,
			Error:  "something went wrong with our server",
		})
		return
	}

	h.logger.Info(
		"patient registered successfully",
		zap.String("patient name", req.FullName),
	)

	status := http.StatusOK
	render.Status(r, status)
	render.JSON(w, r, models.SuccessResponse{
		Status:  status,
		Message: "patient registered successfully",
		Data:    p,
	})
}

// HandleUpdatePatientDetails godoc
// @Summary      Update a patient
// @Description  Updates details of an existing patient by patient ID.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        patientID  path      string                   true  "Patient ID (UUID)"
// @Param        body       body      models.UpdatePatientReq  true  "Updated patient data"
// @Success      200        {object}  models.SuccessResponse
// @Failure      400        {object}  models.FailureResponse
// @Failure      422        {object}  models.FailureResponse
// @Failure      500        {object}  models.FailureResponse
// @Router       /v1/patient/{patientID} [put]
func (h *handler) HandleUpdatePatientDetails(w http.ResponseWriter, r *http.Request) {
	patientID := chi.URLParam(r, "patientID")
	h.logger.Info("identifier", zap.String("patient", patientID))
	if err := h.validate.Var(patientID, "required,uuid"); err != nil {
		unprocessableEntityResponse(w, r)
		return
	}

	var req models.UpdatePatientReq
	if err := helpers.DecodeJSON(r, &req); err != nil {
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

	p, err := h.store.Patient.Update(ctx, &req)
	if err != nil {
		if ok := errors.Is(err, ErrPatientNotFound); ok {
			notFoundError(w, r)
			return
		}
		serverErrorResponse(w, r)
		return
	}

	status := http.StatusOK
	helpers.WriteJSONResponse(w, r, status, models.SuccessResponse{
		Status:  status,
		Message: "patient data updated successfully",
		Data:    p,
	})
}

// HandleDeletePatientDetails godoc
// @Summary      Delete a patient
// @Description  Deletes a patient by their patient ID.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        patientID  path      string  true  "Patient ID (UUID)"
// @Success      200        {object}  models.SuccessResponse
// @Failure      400        {object}  models.FailureResponse
// @Failure      404        {object}  models.FailureResponse
// @Failure      500        {object}  models.FailureResponse
// @Router       /v1/patient/{patientID} [delete]
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

	helpers.WriteJSONResponse(w, r, http.StatusOK, &models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "patient deleted successfully",
	})
}

// HandleListPatients godoc
// @Summary      List patients
// @Description  Lists all registered patients with optional pagination and search.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        page        query     int     false  "Page number"
// @Param        pageSize    query     int     false  "Page size"
// @Param        searchTerm  query     string  false  "Search term (e.g., name or email)"
// @Success      200         {object}  models.SuccessResponse
// @Failure      400         {object}  models.FailureResponse
// @Failure      500         {object}  models.FailureResponse
// @Router       /v1/patient [get]
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

	helpers.WriteJSONResponse(w, r, http.StatusOK, &models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "patients fetched successfully",
		Data:    list,
	})
}

// HandleGetPatient godoc
// @Summary      Get patient details
// @Description  Retrieves a patient's details using their patient ID.
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        patientID  path      string  true  "Patient ID (UUID)"
// @Success      200        {object}  models.SuccessResponse
// @Failure      400        {object}  models.FailureResponse
// @Failure      404        {object}  models.FailureResponse
// @Failure      500        {object}  models.FailureResponse
// @Router       /v1/patient/{patientID} [get]
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

	helpers.WriteJSONResponse(w, r, http.StatusOK, models.SuccessResponse{
		Status:  http.StatusOK,
		Message: "record of the patient fetched successfully",
		Data:    record,
	})
}
