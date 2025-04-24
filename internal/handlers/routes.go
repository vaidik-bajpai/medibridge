package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func (h *handler) Router() http.Handler {
	r := chi.NewRouter()

	// Middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	// Routes
	r.Route("/v1", func(r chi.Router) {
		r.Post("/signup", h.HandleUserSignup)
		r.Post("/signin", h.HandleUserLogin)

		r.Route("/patient", func(r chi.Router) {
			r.Use(h.RequireAuth)
			r.Post("/", h.HandleRegisterPatient)
			r.Get("/list", h.HandleListPatients)

			r.Route("/{patientID}", func(r chi.Router) {
				r.Put("/", h.HandleUpdatePatientDetails)
				r.Delete("/", h.HandleDeletePatientDetails)
				r.Post("/diagnoses", h.HandleAddDiagnoses)
				r.Post("/vitals", h.HandleCaptureVitals)
				r.Put("/vitals", h.HandleUpdatingVitals)
				r.Delete("/vitals", h.HandleDeleteVitals)
				r.Post("/condition", h.HandleAddCondition)
				r.Post("/allergy", h.HandleRecordAllergy)
			})
		})

		r.Route("/diagnoses/{diagnosesID}", func(r chi.Router) {
			r.Put("/", h.HandleUpdateDiagnoses)
			r.Delete("/", h.HandleDeleteDiagnoses)
		})

		r.Route("/conditions/{conditionID}", func(r chi.Router) {
			r.Delete("/", h.HandleInactiveCondition)
		})

		r.Route("/allergy/{allergyID}", func(r chi.Router) {
			r.Put("/", h.HandleUpdateAllergy)
			r.Delete("/", h.HandleDeleteAllergy)
		})
	})

	return r
}
