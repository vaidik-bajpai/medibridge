package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"github.com/vaidik-bajpai/medibridge/internal/store/mocks"
	"go.uber.org/zap"
)

func TestHandleUpdatePatientDetails(t *testing.T) {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"
	invalidUUID := "not-a-uuid"
	validReqJSON := `{"fullname":"John Doe"}`

	tests := []struct {
		name           string
		patientID      string
		body           string
		mockSetup      func(*mocks.PatientStorer)
		expectedStatus int
	}{
		{
			name:           "Invalid UUID",
			patientID:      invalidUUID,
			body:           validReqJSON,
			mockSetup:      func(p *mocks.PatientStorer) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Invalid JSON",
			patientID:      validUUID,
			body:           "{invalid",
			mockSetup:      func(p *mocks.PatientStorer) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Validation fails",
			patientID:      validUUID,
			body:           `{}`,
			mockSetup:      func(p *mocks.PatientStorer) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:      "DB error",
			patientID: validUUID,
			body:      validReqJSON,
			mockSetup: func(p *mocks.PatientStorer) {
				p.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error")).Once()
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:      "Success",
			patientID: validUUID,
			body:      validReqJSON,
			mockSetup: func(p *mocks.PatientStorer) {
				p.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup the mock store
			mockPatient := mocks.NewPatientStorer(t)
			mockSession := mocks.NewSessionStorer(t)
			mockUser := mocks.NewUserStorer(t)
			mockStore := store.NewMockStore(mockPatient, mockSession, mockUser)
			tt.mockSetup(mockPatient)

			// Create handler
			hdl := &handler{
				store:    mockStore,
				validate: validator.New(),
				logger:   zap.NewNop(),
			}

			// Set up the router
			r := chi.NewRouter()
			r.Route("/v1/patients/{patientID}", func(r chi.Router) {
				r.Put("/", hdl.HandleUpdatePatientDetails)
			})

			req := httptest.NewRequest(http.MethodPut, "/v1/patients/"+tt.patientID, bytes.NewBufferString(tt.body))
			mux.Vars(req)
			mux.SetURLVars(req, map[string]string{"patientID": tt.patientID})
			fmt.Println("path:", req.URL.Path)
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Serve the request
			fmt.Println("server called")
			r.ServeHTTP(w, req)

			// Assert the response status
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
