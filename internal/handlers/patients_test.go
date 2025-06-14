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
	"github.com/vaidik-bajpai/medibridge/internal/helpers"
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	"github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

func TestHandleUpdatePatientDetails(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := models.UpdatePatientReq{
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
				ps.On("Update", mock.Anything, mock.MatchedBy(func(r *models.UpdatePatientReq) bool {
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
				ps.On("Update", mock.Anything, mock.MatchedBy(func(r *models.UpdatePatientReq) bool {
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

			req := helpers.InjectURLParam(http.MethodPut, tt.body, "/v1/patient/"+tt.urlID, "patientID", tt.urlID)

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
