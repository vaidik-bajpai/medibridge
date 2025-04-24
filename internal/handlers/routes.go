package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/vaidik-bajpai/medibridge/docs"
)

// @title MediBridge API
// @version 1.0
// @description This is the API documentation for the MediBridge application.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

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

	r.Route("/v1", func(r chi.Router) {
		// Signup
		// @Summary User Signup
		// @Description Register a new user
		// @Tags Auth
		// @Accept json
		// @Produce json
		// @Param user body models.UserSignup true "User Signup Information"
		// @Success 200 {object} models.User
		// @Failure 400 {object} models.ErrorResponse
		// @Router /v1/signup [post]
		r.Post("/signup", h.HandleUserSignup)

		// Signin
		// @Summary User Login
		// @Description Login an existing user
		// @Tags Auth
		// @Accept json
		// @Produce json
		// @Param user body models.UserLogin true "User Login Information"
		// @Success 200 {object} models.User
		// @Failure 400 {object} models.ErrorResponse
		// @Router /v1/signin [post]
		r.Post("/signin", h.HandleUserLogin)

		r.Route("/patient", func(r chi.Router) {
			r.Use(h.RequireAuth)

			// Register Patient
			// @Summary Register Patient
			// @Description Register a new patient
			// @Tags Patient
			// @Accept json
			// @Produce json
			// @Param patient body models.Patient true "Patient Registration Information"
			// @Success 200 {object} models.Patient
			// @Failure 400 {object} models.ErrorResponse
			// @Router /v1/patient [post]
			r.Post("/", h.HandleRegisterPatient)

			// List Patients
			// @Summary List Patients
			// @Description Get a list of patients
			// @Tags Patient
			// @Produce json
			// @Success 200 {array} models.Patient
			// @Failure 400 {object} models.ErrorResponse
			// @Router /v1/patient/list [get]
			r.Get("/list", h.HandleListPatients)

			r.Route("/{patientID}", func(r chi.Router) {
				// Update Patient Details
				// @Summary Update Patient Details
				// @Description Update the details of an existing patient
				// @Tags Patient
				// @Accept json
				// @Produce json
				// @Param patientID path int true "Patient ID"
				// @Param patient body models.Patient true "Updated Patient Information"
				// @Success 200 {object} models.Patient
				// @Failure 400 {object} models.ErrorResponse
				// @Router /v1/patient/{patientID} [put]
				r.Put("/", h.HandleUpdatePatientDetails)

				// Delete Patient Details
				// @Summary Delete Patient Details
				// @Description Delete the details of an existing patient
				// @Tags Patient
				// @Param patientID path int true "Patient ID"
				// @Success 200 {object} models.SuccessResponse
				// @Failure 400 {object} models.ErrorResponse
				// @Router /v1/patient/{patientID} [delete]
				r.Delete("/", h.HandleDeletePatientDetails)

				// Add Diagnoses
				// @Summary Add Diagnoses
				// @Description Add diagnoses to a patient
				// @Tags Patient
				// @Accept json
				// @Produce json
				// @Param patientID path int true "Patient ID"
				// @Param diagnoses body models.Diagnoses true "Diagnoses Information"
				// @Success 200 {object} models.Diagnoses
				// @Failure 400 {object} models.ErrorResponse
				// @Router /v1/patient/{patientID}/diagnoses [post]
				r.Post("/diagnoses", h.HandleAddDiagnoses)

				// Capture Vitals
				// @Summary Capture Vitals
				// @Description Capture vitals of a patient
				// @Tags Patient
				// @Accept json
				// @Produce json
				// @Param patientID path int true "Patient ID"
				// @Param vitals body models.Vitals true "Vitals Information"
				// @Success 200 {object} models.Vitals
				// @Failure 400 {object} models.ErrorResponse
				// @Router /v1/patient/{patientID}/vitals [post]
				r.Post("/vitals", h.HandleCaptureVitals)

				// Update Vitals
				// @Summary Update Vitals
				// @Description Update vitals of a patient
				// @Tags Patient
				// @Accept json
				// @Produce json
				// @Param patientID path int true "Patient ID"
				// @Param vitals body models.Vitals true "Vitals Information"
				// @Success 200 {object} models.Vitals
				// @Failure 400 {object} models.ErrorResponse
				// @Router /v1/patient/{patientID}/vitals [put]
				r.Put("/vitals", h.HandleUpdatingVitals)

				// Delete Vitals
				// @Summary Delete Vitals
				// @Description Delete vitals of a patient
				// @Tags Patient
				// @Param patientID path int true "Patient ID"
				// @Success 200 {object} models.SuccessResponse
				// @Failure 400 {object} models.ErrorResponse
				// @Router /v1/patient/{patientID}/vitals [delete]
				r.Delete("/vitals", h.HandleDeleteVitals)

				// Add Condition
				// @Summary Add Condition
				// @Description Add a condition to a patient
				// @Tags Patient
				// @Param patientID path int true "Patient ID"
				// @Param condition body models.Condition true "Condition Information"
				// @Success 200 {object} models.Condition
				// @Failure 400 {object} models.ErrorResponse
				// @Router /v1/patient/{patientID}/condition [post]
				r.Post("/condition", h.HandleAddCondition)

				// Record Allergy
				// @Summary Record Allergy
				// @Description Record an allergy for a patient
				// @Tags Patient
				// @Param patientID path int true "Patient ID"
				// @Param allergy body models.Allergy true "Allergy Information"
				// @Success 200 {object} models.Allergy
				// @Failure 400 {object} models.ErrorResponse
				// @Router /v1/patient/{patientID}/allergy [post]
				r.Post("/allergy", h.HandleRecordAllergy)
			})
		})

		// Diagnoses Routes
		r.Route("/diagnoses/{diagnosesID}", func(r chi.Router) {
			// Update Diagnoses
			// @Summary Update Diagnoses
			// @Description Update the diagnoses for a patient
			// @Tags Diagnoses
			// @Param diagnosesID path int true "Diagnoses ID"
			// @Param diagnoses body models.Diagnoses true "Updated Diagnoses Information"
			// @Success 200 {object} models.Diagnoses
			// @Failure 400 {object} models.ErrorResponse
			// @Router /v1/diagnoses/{diagnosesID} [put]
			r.Put("/", h.HandleUpdateDiagnoses)

			// Delete Diagnoses
			// @Summary Delete Diagnoses
			// @Description Delete diagnoses from a patient
			// @Tags Diagnoses
			// @Param diagnosesID path int true "Diagnoses ID"
			// @Success 200 {object} models.SuccessResponse
			// @Failure 400 {object} models.ErrorResponse
			// @Router /v1/diagnoses/{diagnosesID} [delete]
			r.Delete("/", h.HandleDeleteDiagnoses)
		})

		// Conditions Routes
		r.Route("/conditions/{conditionID}", func(r chi.Router) {
			// Delete Condition
			// @Summary Delete Condition
			// @Description Delete condition from a patient
			// @Tags Conditions
			// @Param conditionID path int true "Condition ID"
			// @Success 200 {object} models.SuccessResponse
			// @Failure 400 {object} models.ErrorResponse
			// @Router /v1/conditions/{conditionID} [delete]
			r.Delete("/", h.HandleInactiveCondition)
		})

		// Allergy Routes
		r.Route("/allergy/{allergyID}", func(r chi.Router) {
			// Update Allergy
			// @Summary Update Allergy
			// @Description Update allergy details for a patient
			// @Tags Allergy
			// @Param allergyID path int true "Allergy ID"
			// @Param allergy body models.Allergy true "Updated Allergy Information"
			// @Success 200 {object} models.Allergy
			// @Failure 400 {object} models.ErrorResponse
			// @Router /v1/allergy/{allergyID} [put]
			r.Put("/", h.HandleUpdateAllergy)

			// Delete Allergy
			// @Summary Delete Allergy
			// @Description Delete allergy from a patient
			// @Tags Allergy
			// @Param allergyID path int true "Allergy ID"
			// @Success 200 {object} models.SuccessResponse
			// @Failure 400 {object} models.ErrorResponse
			// @Router /v1/allergy/{allergyID} [delete]
			r.Delete("/", h.HandleDeleteAllergy)
		})
	})

	return r
}
