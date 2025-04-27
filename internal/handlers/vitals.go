package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	"go.uber.org/zap"
)

func TestHandleCaptureVitals(t *testing.T) {
	tests := []struct {
		name               string
		patientID          string
		body               dto.CreateVitalReq
		mockSetup          func(*mocks.VitalsStorer)
		expectedStatusCode int
	}{
		{
			name:      "Invalid Patient ID",
			patientID: "invalid-uuid",
			body:      dto.CreateVitalReq{},
			mockSetup: func(cs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "Valid Request",
			patientID: "550e8400-e29b-41d4-a716-446655440000",
			body: dto.CreateVitalReq{
				PatientID: "550e8400-e29b-41d4-a716-446655440000",
				Temperature: 98.6,
				BloodPressure: "120/80",
			},
			mockSetup: func(cs *mocks.VitalsStorer) {
				cs.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:      "DB Error",
			patientID: "550e8400-e29b-41d4-a716-446655440000",
			body: dto.CreateVitalReq{
				PatientID:   "550e8400-e29b-41d4-a716-446655440000",
				Temperature: 98.6,
				BloodPressure: "120/80",
			},
			mockSetup: func(cs *mocks.VitalsStorer) {
				cs.On("Create", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := mocks.NewVitalsStorer(t)
			tt.mockSetup(cs)

			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Vitals: cs},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.patientID)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/patients/"+tt.patientID+"/vitals", bytes.NewReader(body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleCaptureVitals(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleUpdatingVitals(t *testing.T) {
	tests := []struct {
		name               string
		patientID          string
		body               dto.UpdateVitalReq
		mockSetup          func(*mocks.VitalsStorer)
		expectedStatusCode int
	}{
		{
			name:      "Invalid Patient ID",
			patientID: "invalid-uuid",
			body:      dto.UpdateVitalReq{},
			mockSetup: func(cs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "Valid Request",
			patientID: "550e8400-e29b-41d4-a716-446655440000",
			body: dto.UpdateVitalReq{
				PatientID: "550e8400-e29b-41d4-a716-446655440000",
				Temperature: 99.5,
				BloodPressure: "130/85",
			},
			mockSetup: func(cs *mocks.VitalsStorer) {
				cs.On("Update", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:      "DB Error",
			patientID: "550e8400-e29b-41d4-a716-446655440000",
			body: dto.UpdateVitalReq{
				PatientID:   "550e8400-e29b-41d4-a716-446655440000",
				Temperature: 99.5,
				BloodPressure: "130/85",
			},
			mockSetup: func(cs *mocks.VitalsStorer) {
				cs.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := mocks.NewVitalsStorer(t)
			tt.mockSetup(cs)

			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Vitals: cs},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.patientID)

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/patients/"+tt.patientID+"/vitals", bytes.NewReader(body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleUpdatingVitals(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleDeleteVitals(t *testing.T) {
	tests := []struct {
		name               string
		patientID          string
		mockSetup          func(*mocks.VitalsStorer)
		expectedStatusCode int
	}{
		{
			name:      "Invalid Patient ID",
			patientID: "invalid-uuid",
			mockSetup: func(cs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:      "Valid Request",
			patientID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(cs *mocks.VitalsStorer) {
				cs.On("Delete", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:      "DB Error",
			patientID: "550e8400-e29b-41d4-a716-446655440000",
			mockSetup: func(cs *mocks.VitalsStorer) {
				cs.On("Delete", mock.Anything, "550e8400-e29b-41d4-a716-446655440000").Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := mocks.NewVitalsStorer(t)
			tt.mockSetup(cs)

			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Vitals: cs},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.patientID)

			req := httptest.NewRequest(http.MethodDelete, "/patients/"+tt.patientID+"/vitals", nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleDeleteVitals(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
