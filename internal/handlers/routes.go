package handlers

import (
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/vaidik-bajpai/medibridge/docs"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

func (h *handler) Router() http.Handler {
	r := chi.NewRouter()

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

	r.Route("/v1", func(r chi.Router) {
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/v1/swagger/doc.json"),
		))

		r.Route("/user", func(r chi.Router) {
			r.Post("/signup", h.HandleUserSignup)
			r.Post("/signin", h.HandleUserLogin)
		})

		r.Route("/patient", func(r chi.Router) {
			r.Use(h.RequireAuth)
			r.With(h.RequirePaginate).Get("/", h.HandleListPatients)
			r.Post("/", h.HandleRegisterPatient)

			r.Route("/{patientID}", func(r chi.Router) {
				r.Get("/", h.HandleGetPatient)
				r.Put("/", h.HandleUpdatePatientDetails)
				r.Delete("/", h.HandleDeletePatientDetails)

				r.Group(func(r chi.Router) {
					r.Use(h.RequireRole(db.RoleDoctor))
					r.Post("/condition", h.HandleAddCondition)
					r.Post("/allergy", h.HandleRecordAllergy)
					r.Post("/diagnoses", h.HandleAddDiagnoses)

					r.Post("/vitals", h.HandleCaptureVitals)
					r.Put("/vitals", h.HandleUpdatingVitals)
					r.Delete("/vitals", h.HandleDeleteVitals)
				})
			})
		})

		r.Route("/condition/{conditionID}", func(r chi.Router) {
			r.Use(h.RequireAuth)
			r.Use(h.RequireRole(db.RoleDoctor))
			r.Delete("/", h.HandleInactiveCondition)
		})

		r.Route("/allergy/{allergyID}", func(r chi.Router) {
			r.Use(h.RequireAuth)
			r.Use(h.RequireRole(db.RoleDoctor))
			r.Put("/", h.HandleUpdateAllergy)
			r.Delete("/", h.HandleDeleteAllergy)
		})

		r.Route("/diagnoses/{diagnosesID}", func(r chi.Router) {
			r.Use(h.RequireAuth)
			r.Use(h.RequireRole(db.RoleDoctor))
			r.Put("/", h.HandleUpdateDiagnoses)
			r.Delete("/", h.HandleDeleteDiagnoses)
		})
	})

	return r
}
