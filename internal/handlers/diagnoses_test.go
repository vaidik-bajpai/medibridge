package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	"github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

func TestHandleAddDiagnoses(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := models.DiagnosesReq{
		Name: "Asthma",
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.DiagnosesStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Patient UUID",
			urlID: "invalid-uuid",
			body:  body,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Malformed JSON",
			urlID: validUUID,
			body:  []byte(`{invalid-json}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Validation fails on DTO",
			urlID: validUUID,
			body:  []byte(`{"name":""}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Success",
			urlID: validUUID,
			body:  body,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Add", mock.Anything, mock.MatchedBy(func(r *models.DiagnosesReq) bool {
					return r.PID == validUUID && r.Name == "Asthma"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			body:  body,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Add", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := mocks.NewDiagnosesStorer(t)
			tt.mockSetup(ds)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Diagnoses: ds},
				validate: validator.New(),
			}

			req := InjectURLParam(http.MethodPost, tt.body, "/v1/patient/"+tt.urlID+"/diagnoses", "patientID", tt.urlID)

			rr := httptest.NewRecorder()
			h.HandleAddDiagnoses(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleUpdateDiagnoses(t *testing.T) {
	validDiagnosisID := "c0f1e1de-e2d6-4ecf-9320-0dbb3b8c02aa"
	reqBody := models.UpdateDiagnosesReq{
		Name: "Asthma",
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.DiagnosesStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Diagnosis UUID",
			urlID: "invalid-uuid",
			body:  body,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Malformed JSON",
			urlID: validDiagnosisID,
			body:  []byte(`{invalid-json}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Validation fails on DTO",
			urlID: validDiagnosisID,
			body:  []byte(`{"name":""}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Success",
			urlID: validDiagnosisID,
			body:  body,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Update", mock.Anything, mock.MatchedBy(func(r *models.UpdateDiagnosesReq) bool {
					return r.DID == validDiagnosisID && r.Name == "Asthma"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validDiagnosisID,
			body:  body,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := mocks.NewDiagnosesStorer(t)
			tt.mockSetup(ds)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Diagnoses: ds},
				validate: validator.New(),
			}

			req := InjectURLParam(http.MethodPut, tt.body, "/v1/diagnoses/"+tt.urlID, "diagnosesID", tt.urlID)

			rr := httptest.NewRecorder()
			h.HandleUpdateDiagnoses(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleDeleteDiagnoses(t *testing.T) {
	validDiagnosisID := "c0f1e1de-e2d6-4ecf-9320-0dbb3b8c02aa"

	tests := []struct {
		name               string
		urlID              string
		mockSetup          func(*mocks.DiagnosesStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Diagnosis UUID",
			urlID: "invalid-uuid",
			mockSetup: func(ds *mocks.DiagnosesStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Success",
			urlID: validDiagnosisID,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Delete", mock.Anything, validDiagnosisID).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validDiagnosisID,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Delete", mock.Anything, validDiagnosisID).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := mocks.NewDiagnosesStorer(t)
			tt.mockSetup(ds)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Diagnoses: ds},
				validate: validator.New(),
			}

			req := InjectURLParam(http.MethodDelete, []byte(""), "/v1/diagnoses/"+tt.urlID, "diagnosesID", tt.urlID)

			rr := httptest.NewRecorder()
			h.HandleDeleteDiagnoses(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
