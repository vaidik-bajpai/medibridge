package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/mock"
	"github.com/vaidik-bajpai/medibridge/internal/helpers"
	"github.com/vaidik-bajpai/medibridge/internal/mocks"
	"github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/store"
	"go.uber.org/zap"
)

func TestHandleUserLogin_AllScenarios(t *testing.T) {
	validEmail := "test@example.com"
	validPassword := "password"
	hashedPassword, _ := helpers.MakeHashFromToken(validPassword)

	baseUser := &models.UserModel{
		ID:        "user123",
		Email:     validEmail,
		Password:  string(hashedPassword),
		Activated: true,
		Role:      "doctor",
	}

	tests := []struct {
		name           string
		body           string
		mockUser       func(*mocks.UserStorer)
		mockSession    func(*mocks.SessionStorer)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"email":"test@example.com", "password":"password"}`,
			mockUser: func(m *mocks.UserStorer) {
				m.On("FindViaEmail", mock.Anything, validEmail).Return(baseUser, nil)
			},
			mockSession: func(m *mocks.SessionStorer) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.CreateSessReq")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "user not found",
			body: `{"email":"notfound@example.com", "password":"password"}`,
			mockUser: func(m *mocks.UserStorer) {
				m.On("FindViaEmail", mock.Anything, "notfound@example.com").Return(nil, store.ErrNotFound)
			},
			mockSession:    func(m *mocks.SessionStorer) {},
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "wrong password",
			body: `{"email":"test@example.com", "password":"wrongpass"}`,
			mockUser: func(m *mocks.UserStorer) {
				m.On("FindViaEmail", mock.Anything, validEmail).Return(baseUser, nil)
			},
			mockSession:    func(m *mocks.SessionStorer) {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name: "invalid json",
			body: `invalid_json_payload`,
			mockUser: func(m *mocks.UserStorer) {
			},
			mockSession:    func(m *mocks.SessionStorer) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "validation error",
			body: `{"email":"", "password":"password"}`,
			mockUser: func(m *mocks.UserStorer) {
			},
			mockSession:    func(m *mocks.SessionStorer) {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "session creation failure",
			body: `{"email":"test@example.com", "password":"password"}`,
			mockUser: func(m *mocks.UserStorer) {
				m.On("FindViaEmail", mock.Anything, validEmail).Return(baseUser, nil)
			},
			mockSession: func(m *mocks.SessionStorer) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.CreateSessReq")).Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(mocks.UserStorer)
			mockSessionStore := new(mocks.SessionStorer)

			if tt.mockUser != nil {
				tt.mockUser(mockUserStore)
			}
			if tt.mockSession != nil {
				tt.mockSession(mockSessionStore)
			}

			h := NewHandler(validator.New(), zap.NewNop(), &store.Store{
				User:    mockUserStore,
				Session: mockSessionStore,
			})

			req := httptest.NewRequest(http.MethodPost, "/v1/user/signin", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.HandleUserLogin(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestHandleUserSignup_AllScenarios(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		mockUser       func(*mocks.UserStorer)
		expectedStatus int
	}{
		{
			name: "success",
			body: `{"fullname":"John Doe","email":"john@example.com","password":"secure123","role":"doctor"}`,
			mockUser: func(m *mocks.UserStorer) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.SignupReq")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "invalid json",
			body: `{"fullname": "John", "email": "invalid`, // malformed JSON
			mockUser: func(m *mocks.UserStorer) {
				// no store call
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "validation error",
			body: `{"fullname":"","email":" ","password":"","role":""}`,
			mockUser: func(m *mocks.UserStorer) {
				// no store call
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "store error",
			body: `{"fullname":"Jane Doe","email":"jane@example.com","password":"secure123","role":"receptionist"}`,
			mockUser: func(m *mocks.UserStorer) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*models.SignupReq")).Return(errors.New("db error"))
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserStore := new(mocks.UserStorer)
			if tt.mockUser != nil {
				tt.mockUser(mockUserStore)
			}

			h := NewHandler(validator.New(), zap.NewNop(), &store.Store{
				User: mockUserStore,
			})

			req := httptest.NewRequest(http.MethodPost, "/v1/user/signup", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			h.HandleUserSignup(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}

func TestHandleUserLogout(t *testing.T) {
	h := NewHandler(validator.New(), zap.NewNop(), &store.Store{})

	req := httptest.NewRequest(http.MethodPost, "/v1/user/logout", nil)
	w := httptest.NewRecorder()

	h.HandleUserLogout(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check for the deleted cookie
	cookieFound := false
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "medibridge-token" {
			cookieFound = true
			assert.Equal(t, "", cookie.Value)
			assert.Equal(t, "/", cookie.Path)
			assert.Equal(t, true, cookie.HttpOnly)
			assert.Equal(t, int(-1), cookie.MaxAge)
			assert.Equal(t, true, cookie.Expires.Before(time.Now()))
		}
	}
	if !cookieFound {
		t.Error("Expected logout cookie to be set")
	}

	// Check the JSON body
	var res models.SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		t.Fatal("Expected valid JSON response")
	}

	assert.Equal(t, http.StatusOK, res.Status)
	assert.Equal(t, "user logged out successfully", res.Message)
}
