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
	"github.com/vaidik-bajpai/medibridge/internal/models"

	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

func TestHandleDeleteVitals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name               string
		patientID          string
		mockSetup          func(*mocks.VitalsStorer)
		expectedStatusCode int
	}{
		{
			name:      "Invalid UUID",
			patientID: "not-a-uuid",
			mockSetup: func(ps *mocks.VitalsStorer) {
				// Should not be called
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Vitals deleted successfully",
			patientID: validUUID,
			mockSetup: func(ps *mocks.VitalsStorer) {
				ps.On("Delete", mock.Anything, validUUID).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:      "Vitals delete DB error",
			patientID: validUUID,
			mockSetup: func(ps *mocks.VitalsStorer) {
				ps.On("Delete", mock.Anything, validUUID).Return(errors.New("db error")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVitals := mocks.NewVitalsStorer(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockVitals)
			}

			h := NewHandler(validator.New(), zap.NewNop(), &store.Store{
				Vitals: mockVitals,
			})

			req := httptest.NewRequest(http.MethodDelete, "/v1/patients/"+tt.patientID+"/vitals", nil)

			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("patientID", tt.patientID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rec := httptest.NewRecorder()
			h.HandleDeleteVitals(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			require.Equal(t, tt.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestHandleUpdatingVitals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	validBody := []byte(`{
		"heightCm":               170,
		"weightKg":               65,
		"bmi":                    22.5,
		"temperatureC":           36.6,
		"pulse":                  80,
		"respiratoryRate":        18,
		"bloodPressureSystolic":  120,
		"bloodPressureDiastolic": 80,
		"oxygenSaturation":       98
	}`)

	tests := []struct {
		name               string
		patientID          string
		body               []byte
		mockSetup          func(*mocks.VitalsStorer)
		expectedStatusCode int
	}{
		{
			name:               "Invalid UUID",
			patientID:          "invalid-uuid",
			body:               validBody,
			mockSetup:          func(ps *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Malformed JSON",
			patientID:          validUUID,
			body:               []byte(`{"heightCm":}`),
			mockSetup:          func(ps *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Validation fails on models",
			patientID:          validUUID,
			body:               []byte(`{"heightCm": -1}`),
			mockSetup:          func(ps *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Vitals update DB error",
			patientID: validUUID,
			body:      validBody,
			mockSetup: func(ps *mocks.VitalsStorer) {
				ps.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:      "Vitals update success",
			patientID: validUUID,
			body:      validBody,
			mockSetup: func(ps *mocks.VitalsStorer) {
				ps.On("Update", mock.Anything, mock.MatchedBy(func(r *models.UpdateVitalReq) bool {
					return r.PatientID == validUUID && *r.BMI == 22.5
				})).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVitals := mocks.NewVitalsStorer(t)
			if tt.mockSetup != nil {
				tt.mockSetup(mockVitals)
			}

			h := NewHandler(validator.New(), zap.NewNop(), &store.Store{
				Vitals: mockVitals,
			})

			req := httptest.NewRequest(http.MethodPut, "/v1/patient/"+tt.patientID+"/vitals", bytes.NewReader(tt.body))

			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("patientID", tt.patientID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rec := httptest.NewRecorder()
			h.HandleUpdatingVitals(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			require.Equal(t, tt.expectedStatusCode, res.StatusCode)
		})
	}
}

func TestHandleCaptureVitals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	validBody := []byte(`{
		"heightCm":               170,
		"weightKg":               60,
		"bmi":                    20.8,
		"temperatureC":           36.8,
		"pulse":                  78,
		"respiratoryRate":        16,
		"bloodPressureSystolic":  120,
		"bloodPressureDiastolic": 80,
		"oxygenSaturation":       97
	}`)

	tests := []struct {
		name               string
		patientID          string
		body               []byte
		setupMock          func(*mocks.VitalsStorer)
		expectedStatusCode int
	}{
		{
			name:               "Invalid UUID",
			patientID:          "not-a-uuid",
			body:               validBody,
			setupMock:          func(m *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Malformed JSON body",
			patientID:          validUUID,
			body:               []byte(`{"heightCm":}`),
			setupMock:          func(m *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Validation failure",
			patientID:          validUUID,
			body:               []byte(`{"heightCm": -1}`), // invalid field value
			setupMock:          func(m *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:      "Unique constraint violation",
			patientID: validUUID,
			body:      validBody,
			setupMock: func(m *mocks.VitalsStorer) {
				m.On("Create", mock.Anything, mock.Anything).Return(store.ErrUniqueConstraintViolated).Once()
			},
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:      "Internal server error",
			patientID: validUUID,
			body:      validBody,
			setupMock: func(m *mocks.VitalsStorer) {
				m.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:      "Vitals captured successfully",
			patientID: validUUID,
			body:      validBody,
			setupMock: func(m *mocks.VitalsStorer) {
				m.On("Create", mock.Anything, mock.MatchedBy(func(r *models.CreateVitalReq) bool {
					return r.PatientID == validUUID && *r.BMI == 20.8
				})).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockVitals := mocks.NewVitalsStorer(t)
			tt.setupMock(mockVitals)

			h := NewHandler(validator.New(), zap.NewNop(), &store.Store{
				Vitals: mockVitals,
			})

			req := httptest.NewRequest(http.MethodPost, "/v1/patient/"+tt.patientID+"/vitals", bytes.NewReader(tt.body))
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("patientID", tt.patientID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			rec := httptest.NewRecorder()
			h.HandleCaptureVitals(rec, req)

			res := rec.Result()
			defer res.Body.Close()
			require.Equal(t, tt.expectedStatusCode, res.StatusCode)
		})
	}
}
