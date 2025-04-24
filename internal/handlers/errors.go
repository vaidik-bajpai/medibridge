package handlers

import (
	"net/http"

	"github.com/go-chi/render"
)

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	render.Status(r, status)
	render.JSON(w, r, map[string]string{
		"error": message,
	})
}

func serverErrorResponse(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusInternalServerError, "something went wrong with our servers")
}

func badRequestResponse(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusBadRequest, "invalid request payload")
}

func unprocessableEntityResponse(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusUnprocessableEntity, "validation failed")
}

func notFoundError(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusNotFound, "404 not found")
}

func unauthorisedErrorResponse(w http.ResponseWriter, r *http.Request, message string) {
	errorResponse(w, r, http.StatusUnauthorized, message)
}

func conflictErrorResponse(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusConflict, "two unique resource are in conflict")
}

func forbiddenErrorResponse(w http.ResponseWriter, r *http.Request) {
	errorResponse(w, r, http.StatusForbidden, "forbidden resource")
}
