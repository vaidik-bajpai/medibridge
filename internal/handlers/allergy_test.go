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
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	"github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

func TestHandleRecordAllergy(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := models.RegAllergyReq{
		Name:     "Peanut",
		Severity: "moderate",
		Reaction: "Coughing",
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.AllergyStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Patient UUID",
			urlID: "invalid-uuid",
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Malformed JSON",
			urlID: validUUID,
			body:  []byte(`{invalid-json}`),
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Validation fails on models",
			urlID: validUUID,
			body:  []byte(`{"Name":""}`),
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Success",
			urlID: validUUID,
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Record", mock.Anything, mock.MatchedBy(func(r *models.RegAllergyReq) bool {
					return r.PatientID == validUUID
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Record", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := mocks.NewAllergyStorer(t)
			tt.mockSetup(as)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Allergy: as},
				validate: validator.New(),
			}

			req := InjectURLParam(http.MethodPost, tt.body, "/v1/patient"+tt.urlID+"/allergy", "patientID", tt.urlID)

			rr := httptest.NewRecorder()
			h.HandleRecordAllergy(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleUpdateAllergy(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	name := "Peanut"
	severity := "severe"
	reqBody := models.UpdateAllergyReq{
		Name:     &name,
		Severity: &severity,
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.AllergyStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Allergy UUID",
			urlID: "invalid-uuid",
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Malformed JSON",
			urlID: validUUID,
			body:  []byte(`{invalid-json}`),
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Validation fails on models",
			urlID: validUUID,
			body:  []byte(`{"Name":""}`),
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Success",
			urlID: validUUID,
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Update", mock.Anything, mock.MatchedBy(func(r *models.UpdateAllergyReq) bool {
					return r.AllergyID == validUUID && *r.Name == "Peanut"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			body:  body,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Update", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := mocks.NewAllergyStorer(t)
			tt.mockSetup(as)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Allergy: as},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("allergyID", tt.urlID)

			req := httptest.NewRequest(http.MethodPut, "/allergies/"+tt.urlID, bytes.NewReader(tt.body))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleUpdateAllergy(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleDeleteAllergy(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name               string
		urlID              string
		mockSetup          func(*mocks.AllergyStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Allergy UUID",
			urlID: "invalid-uuid",
			mockSetup: func(as *mocks.AllergyStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Success",
			urlID: validUUID,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Delete", mock.Anything, validUUID).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			mockSetup: func(as *mocks.AllergyStorer) {
				as.On("Delete", mock.Anything, validUUID).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			as := mocks.NewAllergyStorer(t)
			tt.mockSetup(as)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Allergy: as},
				validate: validator.New(),
			}

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("allergyID", tt.urlID)

			req := httptest.NewRequest(http.MethodDelete, "/allergies/"+tt.urlID, nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			rr := httptest.NewRecorder()
			h.HandleDeleteAllergy(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
