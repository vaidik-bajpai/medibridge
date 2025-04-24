package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
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
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // Point to the generated Swagger JSON URL
	))

	// User Authentication Routes
	r.Route("/v1", func(r chi.Router) {
		r.Post("/signup", h.HandleUserSignup)
		r.Post("/signin", h.HandleUserLogin)

		// Patient CRUD - Receptionist and Doctor have access
		r.Route("/patient", func(r chi.Router) {
			r.Use(h.RequireAuth)
			// Receptionist can perform Patient CRUD
			r.Use(h.RequireRole(db.RoleReceptionist, db.RoleDoctor)) // Both Receptionist and Doctor can access patient routes
			r.Post("/", h.HandleRegisterPatient)
			r.Get("/list", h.HandleListPatients)

			r.Route("/{patientID}", func(r chi.Router) {
				r.Put("/", h.HandleUpdatePatientDetails)
				r.Delete("/", h.HandleDeletePatientDetails)
			})
		})

		// Doctor-level routes - Doctor can manage everything
		r.Route("/patient/{patientID}/diagnoses", func(r chi.Router) {
			r.Use(h.RequireAuth)
			r.Use(h.RequireRole(db.RoleDoctor)) // Only Doctor can manage diagnoses
			r.Post("/", h.HandleAddDiagnoses)
			r.Put("/", h.HandleUpdateDiagnoses)
			r.Delete("/", h.HandleDeleteDiagnoses)
		})

		r.Route("/patient/{patientID}/vitals", func(r chi.Router) {
			r.Use(h.RequireAuth)
			r.Use(h.RequireRole(db.RoleDoctor)) // Only Doctor can manage vitals
			r.Post("/", h.HandleCaptureVitals)
			r.Put("/", h.HandleUpdatingVitals)
			r.Delete("/", h.HandleDeleteVitals)
		})

		r.Route("/patient/{patientID}/condition", func(r chi.Router) {
			r.Use(h.RequireAuth)
			r.Use(h.RequireRole(db.RoleDoctor)) // Only Doctor can manage conditions
			r.Post("/", h.HandleAddCondition)
			r.Delete("/", h.HandleInactiveCondition)
		})

		r.Route("/patient/{patientID}/allergy", func(r chi.Router) {
			r.Use(h.RequireAuth)
			r.Use(h.RequireRole(db.RoleDoctor)) // Only Doctor can manage allergies
			r.Post("/", h.HandleRecordAllergy)
			r.Put("/", h.HandleUpdateAllergy)
			r.Delete("/", h.HandleDeleteAllergy)
		})
	})

	return r
}
