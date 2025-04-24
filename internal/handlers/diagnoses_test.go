package handlers

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	dto "github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

func TestHandleAddDiagnoses(t *testing.T) {
	validPatientID := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.DiagnosesStorer)
		expectedStatusCode int
	}{
		{
			name:  "VALID PATIENT ID",
			urlID: validPatientID,
			body:  []byte(`{"name":"Arthiritis"}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Add", mock.Anything, mock.MatchedBy(func(d *dto.DiagnosesReq) bool {
					return d.PID == validPatientID && d.Name == "Arthiritis"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		}, {
			name:  "IN-VALID PATIENT ID",
			urlID: "invalid",
			body: []byte(`{
				"name":"Arthiritis",
			}`),
			mockSetup:          func(ds *mocks.DiagnosesStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		}, {
			name:  "NOT FOUND",
			urlID: validPatientID,
			body:  []byte(`{"name":"Arthiritis"}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Add", mock.Anything, mock.MatchedBy(func(r *dto.DiagnosesReq) bool {
					return r.PID == validPatientID
				})).Return(ErrPatientNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		}, {
			name:  "INTERNAL SERVER ERROR",
			urlID: validPatientID,
			body:  []byte(`{"name":"Arthiritis"}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Add", mock.Anything, mock.MatchedBy(func(r *dto.DiagnosesReq) bool {
					return r.PID == validPatientID
				})).Return(errors.New("some error"))
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

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.urlID)

			req := httptest.NewRequest(http.MethodPost, "/v1/patients/"+tt.urlID+"/diagnoses", bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleAddDiagnoses(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleUpdateDiagnoses(t *testing.T) {
	validDiagnosesID := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.DiagnosesStorer)
		expectedStatusCode int
	}{
		{
			name:  "VALID UPDATE",
			urlID: validDiagnosesID,
			body:  []byte(`{"name":"Arthritis"}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Update", mock.Anything, mock.MatchedBy(func(req *dto.UpdateDiagnosesReq) bool {
					return req.DID == validDiagnosesID && req.Name == "Arthritis"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "VALID UPDATE WITH TRAILING SPACES",
			urlID: validDiagnosesID,
			body:  []byte(`{"name":"   Arthritis   "}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Update", mock.Anything, mock.MatchedBy(func(req *dto.UpdateDiagnosesReq) bool {
					return req.DID == validDiagnosesID && req.Name == "Arthritis"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "INVALID UUID IN URL",
			urlID:              "not-a-valid-uuid",
			body:               []byte(`{"name":"Arthritis"}`),
			mockSetup:          func(ds *mocks.DiagnosesStorer) {}, // Won’t be called
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "NAME FIELD TOO SHORT",
			urlID:              validDiagnosesID,
			body:               []byte(`{"name":"A"}`),
			mockSetup:          func(ds *mocks.DiagnosesStorer) {}, // Validation fails
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "NAME FIELD EMPTY",
			urlID:              validDiagnosesID,
			body:               []byte(`{"name":""}`),
			mockSetup:          func(ds *mocks.DiagnosesStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "MISSING NAME FIELD",
			urlID:              validDiagnosesID,
			body:               []byte(`{}`),
			mockSetup:          func(ds *mocks.DiagnosesStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "INTERNAL SERVER ERROR",
			urlID: validDiagnosesID,
			body:  []byte(`{"name":"Arthritis"}`),
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Update", mock.Anything, mock.MatchedBy(func(req *dto.UpdateDiagnosesReq) bool {
					return req.DID == validDiagnosesID
				})).Return(errors.New("some DB error"))
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

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("diagnosesID", tt.urlID)

			req := httptest.NewRequest(http.MethodPut, "/v1/patient/"+tt.urlID+"/diagnoses", bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleUpdateDiagnoses(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleDeleteDiagnoses(t *testing.T) {
	validPatientID := "f47ac10b-58cc-4372-a567-0e02b2c3d479"

	tests := []struct {
		name               string
		urlID              string
		mockSetup          func(ds *mocks.DiagnosesStorer)
		expectedStatusCode int
	}{
		{
			name:  "VALID DELETE",
			urlID: validPatientID,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Delete", mock.Anything, validPatientID).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "INVALID UUID FORMAT",
			urlID:              "not-a-uuid",
			mockSetup:          func(ds *mocks.DiagnosesStorer) {}, // won't be called
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "INTERNAL SERVER ERROR FROM STORE",
			urlID: validPatientID,
			mockSetup: func(ds *mocks.DiagnosesStorer) {
				ds.On("Delete", mock.Anything, validPatientID).Return(errors.New("db error"))
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

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("diagnosesID", tt.urlID)

			req := httptest.NewRequest(http.MethodDelete, "/v1/patient/"+tt.urlID+"/diagnoses", bytes.NewReader([]byte("")))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleDeleteDiagnoses(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
