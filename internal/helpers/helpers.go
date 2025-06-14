package helpers

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

func DecodeJSON(r *http.Request, into interface{}) error {
	return json.NewDecoder(r.Body).Decode(into)
}

func MakeHashFromToken(password string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func MatchPassword(hash string, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if ok := errors.Is(err, bcrypt.ErrMismatchedHashAndPassword); ok {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func GenerateSessionToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func WriteJSONResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	render.Status(r, status)
	render.JSON(w, r, data)
}

func InjectURLParam(method string, body []byte, path string, key, val string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, val)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
}
