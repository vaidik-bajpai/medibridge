package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vaidik-bajpai/medibridge/internal/helpers"
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	dto "github.com/vaidik-bajpai/medibridge/internal/models"
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
			expectedStatusCode: http.StatusBadRequest,
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
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "TEMPERATURE BELOW MINIMUM",
			urlID:              validPatientID,
			body:               []byte(`{"temperatureC": 29.9}`),
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "TEMPERATURE ABOVE MAXIMUM",
			urlID:              validPatientID,
			body:               []byte(`{"temperatureC": 45.1}`),
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "OXYGEN SATURATION ABOVE MAXIMUM",
			urlID:              validPatientID,
			body:               []byte(`{"oxygenSaturation": 101}`),
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "NEGATIVE HEIGHT",
			urlID:              validPatientID,
			body:               []byte(`{"heightCm": -150}`),
			mockSetup:          func(vs *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusBadRequest,
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

			req := helpers.InjectURLParam(http.MethodPost, tt.body, "/v1/patient/"+tt.urlID+"/vitals", "patientID", tt.urlID)
			rr := httptest.NewRecorder()

			h.HandleUpdatePatientDetails(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleUpdatingVitals(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	// **Valid** JSON object
	reqBody := []byte(`{
		"heightCm":               176.5,
		"weightKg":               71.0,
		"bmi":                    23.0,
		"temperatureC":           37.0,
		"pulse":                  85,
		"respiratoryRate":        18,
		"bloodPressureSystolic":  125,
		"bloodPressureDiastolic": 85,
		"oxygenSaturation":       97.0
	}`)

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
			body:               reqBody,
			mockSetup:          func(ps *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Malformed JSON",
			urlID:              validUUID,
			body:               []byte(`{invalid-json}`),
			mockSetup:          func(ps *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:               "Validation fails on DTO",
			urlID:              validUUID,
			body:               []byte(`{"heightCm":-1}`),
			mockSetup:          func(ps *mocks.VitalsStorer) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Vitals update DB error",
			urlID: validUUID,
			body:  reqBody,
			mockSetup: func(ps *mocks.VitalsStorer) {
				ps.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:  "Vitals update success",
			urlID: validUUID,
			body:  reqBody,
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
			vs := mocks.NewVitalsStorer(t)
			tt.mockSetup(vs)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Vitals: vs},
				validate: validator.New(),
			}

			req := helpers.InjectURLParam(http.MethodPut, tt.body, "/v1/patient/"+tt.urlID+"/vitals", "patientID", tt.urlID)

			rr := httptest.NewRecorder()
			h.HandleUpdatingVitals(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
