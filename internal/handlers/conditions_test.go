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
	"github.com/vaidik-bajpai/medibridge/internal/helpers"
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	"github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

func TestHandleAddCondition(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"
	reqBody := models.AddConditionReq{
		Condition: "Asthma",
	}
	body, _ := json.Marshal(reqBody)

	tests := []struct {
		name               string
		urlID              string
		body               []byte
		mockSetup          func(*mocks.ConditionStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Patient UUID",
			urlID: "invalid-uuid",
			body:  body,
			mockSetup: func(cs *mocks.ConditionStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Malformed JSON",
			urlID: validUUID,
			body:  []byte(`{invalid-json}`),
			mockSetup: func(cs *mocks.ConditionStorer) {
			},
			expectedStatusCode: http.StatusUnprocessableEntity,
		},
		{
			name:  "Validation fails on DTO",
			urlID: validUUID,
			body:  []byte(`{"conditionName":""}`),
			mockSetup: func(cs *mocks.ConditionStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Success",
			urlID: validUUID,
			body:  body,
			mockSetup: func(cs *mocks.ConditionStorer) {
				cs.On("Add", mock.Anything, mock.MatchedBy(func(r *models.AddConditionReq) bool {
					return r.PatientID == validUUID && r.Condition == "Asthma"
				})).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			body:  body,
			mockSetup: func(cs *mocks.ConditionStorer) {
				cs.On("Add", mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := mocks.NewConditionStorer(t)
			tt.mockSetup(cs)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Conditions: cs},
				validate: validator.New(),
			}

			req := helpers.InjectURLParam(http.MethodPost, tt.body, "/v1/patient/"+tt.urlID+"/condition", "patientID", tt.urlID)

			rr := httptest.NewRecorder()
			h.HandleAddCondition(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}

func TestHandleInactiveCondition(t *testing.T) {
	validUUID := "550e8400-e29b-41d4-a716-446655440000"

	tests := []struct {
		name               string
		urlID              string
		mockSetup          func(*mocks.ConditionStorer)
		expectedStatusCode int
	}{
		{
			name:  "Invalid Condition UUID",
			urlID: "invalid-uuid",
			mockSetup: func(cs *mocks.ConditionStorer) {
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:  "Success",
			urlID: validUUID,
			mockSetup: func(cs *mocks.ConditionStorer) {
				cs.On("Delete", mock.Anything, validUUID).Return(nil)
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "DB Error",
			urlID: validUUID,
			mockSetup: func(cs *mocks.ConditionStorer) {
				cs.On("Delete", mock.Anything, validUUID).Return(errors.New("db error"))
			},
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := mocks.NewConditionStorer(t)
			tt.mockSetup(cs)
			l, _ := zap.NewDevelopment()

			h := &handler{
				logger:   l,
				store:    &store.Store{Conditions: cs},
				validate: validator.New(),
			}

			req := helpers.InjectURLParam(http.MethodDelete, []byte(""), "/v1/condition/"+tt.urlID, "conditionID", tt.urlID)

			rr := httptest.NewRecorder()
			h.HandleInactiveCondition(rr, req)

			require.Equal(t, tt.expectedStatusCode, rr.Code)
		})
	}
}
