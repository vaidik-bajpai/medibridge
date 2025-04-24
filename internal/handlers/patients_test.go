package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
	"gopkg.in/go-playground/assert.v1"
)

func TestHandleRegisterPatient(t *testing.T) {
	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.PatientStorer)
		expectedStatusCode int
		expectedBody       map[string]string
	}{
		{
			name:  "valid request",
			urlID: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			body: []byte(`{
				"fullName": "John Doe",
				"gender": "MALE",
				"dob": "1990-01-01T00:00:00Z",
				"age": 30,
				"contactNo": "1234567890",
				"address": "123 Main St",
				"emergencyName": "Jane Doe",
				"emergencyRelation": "Sister",
				"emergencyPhone": "0987654321",
				"regByID": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
			}`),
			mockSetup: func(m *mocks.PatientStorer) {
				// Mock the PatientStorer Create method to return nil (no error)
				m.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedBody: map[string]string{
				"message": "patient registered successfully",
			},
		},
		{
			name:  "invalid request - missing fullName",
			urlID: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			body: []byte(`{
				"fullName": "",
				"gender": "MALE",
				"dob": "1990-01-01T00:00:00Z",
				"age": 30,
				"contactNo": "1234567890",
				"address": "123 Main St",
				"emergencyName": "Jane Doe",
				"emergencyRelation": "Sister",
				"emergencyPhone": "0987654321",
				"regByID": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
			}`),
			mockSetup: func(m *mocks.PatientStorer) {
				// This test doesn't call the Create method, so no mock is needed here
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
			expectedBody: map[string]string{
				"error": "validation failed",
			},
		},
		{
			name:  "store error",
			urlID: "f47ac10b-58cc-4372-a567-0e02b2c3d479",
			body: []byte(`{
				"fullName": "John Doe",
				"gender": "MALE",
				"dob": "1990-01-01T00:00:00Z",
				"age": 30,
				"contactNo": "1234567890",
				"address": "123 Main St",
				"emergencyName": "Jane Doe",
				"emergencyRelation": "Sister",
				"emergencyPhone": "0987654321",
				"regByID": "f47ac10b-58cc-4372-a567-0e02b2c3d479"
			}`),
			mockSetup: func(m *mocks.PatientStorer) {
				// Mock the PatientStorer Create method to return an error
				m.On("Create", mock.Anything, mock.Anything).Return(errors.New("store error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedBody: map[string]string{
				"error": "something went wrong with our server",
			},
		},
	}

	// Set up the mocks
	mockPatientStorer := mocks.NewPatientStorer(t)
	mockStore := store.NewMockStore(t)

	v := validator.New()
	l := zap.NewNop()
	h := NewHandler(v, l, mockStore)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock behavior for the test case
			tt.mockSetup(mockPatientStorer)

			// Prepare the request
			req := httptest.NewRequest("POST", "/v1/patient", bytes.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			// Add user context if necessary
			ctx := context.WithValue(req.Context(), userCtx, &dto.UserModel{ID: tt.urlID})
			req = req.WithContext(ctx)

			// Record the response
			w := httptest.NewRecorder()

			// Call the handler
			h.HandleRegisterPatient(w, req)

			// Assert the status code
			assert.Equal(t, tt.expectedStatusCode, w.Code)

			// Assert the response body
			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			if err != nil {
				t.Fatalf("failed to unmarshal response: %v", err)
			}

			// Assert that the response matches the expected body
			assert.Equal(t, tt.expectedBody, response)

			// Assert that the mock store's Create method was called
			mockPatientStorer.AssertExpectations(t)
		})
	}
}

func TestHandleUpdatePatientDetails(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := dto.UpdatePatientReq{
		FullName:      ptrToString("John Doe"),
		Gender:        ptrToString("MALE"),
		Age:           ptrToInt(30),
		ContactNumber: ptrToString("1234567890"),
		Address:       ptrToString("123 Street"),
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.PatientStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid UUID",
			urlID: "invalid-uuid",
			body:  body,
			mockSetup: func(ps *mocks.PatientStorer) {

			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Malformed JSON",
			urlID: validUUID,
			body:  []byte(`{invalid-json}`),
			mockSetup: func(ps *mocks.PatientStorer) {

			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Validation fails on DTO",
			urlID: validUUID,
			body:  []byte(`{"fullname":""}`),
			mockSetup: func(ps *mocks.PatientStorer) {

			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Patient not found",
			urlID: validUUID,
			body:  body,
			mockSetup: func(ps *mocks.PatientStorer) {
				ps.On("Update", mock.Anything, mock.MatchedBy(func(r *dto.UpdatePatientReq) bool {
					return r.ID == validUUID
				})).Return(ErrPatientNotFound)
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			body:  body,
			mockSetup: func(ps *mocks.PatientStorer) {
				ps.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:  "Success",
			urlID: validUUID,
			body:  body,
			mockSetup: func(ps *mocks.PatientStorer) {
				ps.On("Update", mock.Anything, mock.MatchedBy(func(r *dto.UpdatePatientReq) bool {
					return r.ID == validUUID && *r.FullName == "John Doe"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := mocks.NewPatientStorer(t)
			tt.mockSetup(ps)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Patient: ps},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.urlID)

			req := httptest.NewRequest(http.MethodPut, "/patients/"+tt.urlID, bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleUpdatePatientDetails(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleDeletePatientDetails(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.PatientStorer)
		expectedStatusCode int
	}{
		{
			name:  "valid uuid",
			urlID: validUUID,
			body:  []byte(""),
			mockSetup: func(ps *mocks.PatientStorer) {
				ps.On("Delete", mock.Anything, mock.MatchedBy(func(pID string) bool {
					return pID == validUUID
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		}, {
			name:               "invalid-uuid",
			urlID:              "invalid",
			body:               []byte(""),
			mockSetup:          func(ps *mocks.PatientStorer) {},
			expectedStatusCode: http.StatusUnprocessableEntity,
		}, {
			name:  "patient not found",
			urlID: validUUID,
			body:  []byte(""),
			mockSetup: func(ps *mocks.PatientStorer) {
				ps.On("Delete", mock.Anything, mock.MatchedBy(func(pID string) bool {
					return pID == validUUID
				})).Return(ErrPatientNotFound)
			},
			expectedStatusCode: http.StatusNotFound,
		}, {
			name:  "server error",
			urlID: validUUID,
			body:  []byte(""),
			mockSetup: func(ps *mocks.PatientStorer) {
				ps.On("Delete", mock.Anything, mock.MatchedBy(func(pID string) bool {
					return pID == validUUID
				})).Return(errors.New("server error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := mocks.NewPatientStorer(t)
			tt.mockSetup(ps)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Patient: ps},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("patientID", tt.urlID)

			req := httptest.NewRequest(http.MethodDelete, "/patients/"+tt.urlID, bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleDeletePatientDetails(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func ptrToString(s string) *string {
	return &s
}

func ptrToInt(i int) *int {
	return &i
}
