package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

func TestHandleCaptureVitals(t *testing.T) {
	validPatientID := "f47ac10b-58cc-4372-a567-0e02b2c3d479"
	validBody := []byte(`{
		"heightCm":170.5,
		"weightKg":65.2,
		"bmi":22.4,
		"temperatureC":36.7,
		"pulse":72,
		"respiratoryRate":16,
		"bloodPressureSystolic":120,
		"bloodPressureDiastolic":80,
		"oxygenSaturation":98
	}`)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.VitalsStorer)
		expectedStatusCode int
	}{
		{
			name:               "INVALID UUID FORMAT",
			urlID:              "invalid-uuid",
			body:               validBody,
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "MALFORMED JSON BODY",
			urlID:              validPatientID,
			body:               []byte(`{bad json}`),
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "MISSING REQUIRED FIELD (Temperature)",
			urlID:              validPatientID,
			body:               []byte(`{"heightCm":170,"weightKg":70}`), // other fields missing
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "TEMPERATURE BELOW MINIMUM",
			urlID:              validPatientID,
			body:               []byte(`{"temperatureC": 29.9}`),
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "TEMPERATURE ABOVE MAXIMUM",
			urlID:              validPatientID,
			body:               []byte(`{"temperatureC": 45.1}`),
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "OXYGEN SATURATION ABOVE MAXIMUM",
			urlID:              validPatientID,
			body:               []byte(`{"oxygenSaturation": 101}`),
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "NEGATIVE HEIGHT",
			urlID:              validPatientID,
			body:               []byte(`{"heightCm": -150}`),
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vs := mocks.NewVitalsStorer(t)
			tt.mockSetup(vs)
			l, _ := zap.NewDevelopment()
			defer l.Sync()
			logger := l.With(zap.String("test_name", tt.name))
			v := validator.New()
			h := NewHandler(v, logger, store.NewMockStore(t))

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.urlID)

			req := httptest.NewRequest(http.MethodPost, "/patients/"+tt.urlID+"/vitals", bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleUpdatePatientDetails(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleUpdatingVitals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := dto.UpdateVitalReq{
		HeightCm:               ptrToFloat64(176.5),
		WeightKg:               ptrToFloat64(71.0),
		BMI:                    ptrToFloat64(23.0),
		TemperatureC:           ptrToFloat64(37.0),
		Pulse:                  ptrToInt(85),
		RespiratoryRate:        ptrToInt(18),
		BloodPressureSystolic:  ptrToInt(125),
		BloodPressureDiastolic: ptrToInt(85),
		OxygenSaturation:       ptrToFloat64(97.0),
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.VitalsStorer)
		expectedStatusCode int
	}{
		{
			name:               "Invalid UUID",
			urlID:              "invalid-uuid",
			body:               body,
			mockSetup:          func(ps *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Malformed JSON",
			urlID:              validUUID,
			body:               []byte(`{invalid-json}`),
			mockSetup:          func(ps *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Validation fails on DTO",
			urlID:              validUUID,
			body:               []byte(`{"heightCm":-1}`),
			mockSetup:          func(ps *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Vitals update DB error",
			urlID: validUUID,
			body:  body,
			mockSetup: func(ps *mocks.VitalsStorer) {
				ps.On("Update", mock.Anything, mock.MatchedBy(func(r *dto.UpdateVitalReq) bool {
					return r.PatientID == validUUID
				})).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:  "Vitals update success",
			urlID: validUUID,
			body:  body,
			mockSetup: func(ps *mocks.VitalsStorer) {
				ps.On("Update", mock.Anything, mock.MatchedBy(func(r *dto.UpdateVitalReq) bool {
					return r.PatientID == validUUID && *r.HeightCm == 176.5
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := mocks.NewVitalsStorer(t)
			tt.mockSetup(ps)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Vitals: ps},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.urlID)

			req := httptest.NewRequest(http.MethodPut, "/vitals/"+tt.urlID, bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleUpdatingVitals(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
