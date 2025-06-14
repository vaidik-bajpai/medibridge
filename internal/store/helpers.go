package store

import (
	"fmt"
	"strings"

	dto "github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

// PreparePatientUpdateParams prepares the parameters for updating a patient's data.
// It returns a slice of db.PatientSetParam, which can be used for the actual update query.
func preparePatientUpdateParams(input *dto.UpdatePatientReq) []db.PatientSetParam {
	var params []db.PatientSetParam

	// Helper function to conditionally add params to the slice
	addParam := func(p db.PatientSetParam) {
		params = append(params, p)
	}

	// Only add parameters if they are non-nil and non-empty
	if input.FullName != nil && *input.FullName != "" {
		addParam(db.Patient.FullName.Set(*input.FullName))
	}

	// Check Gender separately to avoid dereferencing nil pointer
	if input.Gender != nil && *input.Gender != "" {
		addParam(db.Patient.Gender.Set(*input.Gender))
	}

	// Check ContactNumber
	if input.ContactNumber != nil && *input.ContactNumber != "" {
		addParam(db.Patient.ContactNumber.Set(*input.ContactNumber))
	}

	// Check Address
	if input.Address != nil && *input.Address != "" {
		addParam(db.Patient.Address.Set(*input.Address))
	}

	// Check EmergencyName
	if input.EmergencyName != nil && *input.EmergencyName != "" {
		addParam(db.Patient.EmergencyName.Set(*input.EmergencyName))
	}

	// Check EmergencyRelation
	if input.EmergencyRelation != nil && *input.EmergencyRelation != "" {
		addParam(db.Patient.EmergencyRelation.Set(*input.EmergencyRelation))
	}

	// Check EmergencyPhone
	if input.EmergencyPhone != nil && *input.EmergencyPhone != "" {
		addParam(db.Patient.EmergencyPhone.Set(*input.EmergencyPhone))
	}

	// Check DOB
	if input.DOB != nil {
		addParam(db.Patient.DateOfBirth.Set(*input.DOB))
	}

	// Check Age
	if input.Age != nil && *input.Age > 0 && *input.Age <= 100 {
		addParam(db.Patient.Age.Set(*input.Age))
	}

	return params
}

func prepareVitalCreateParams(input *dto.CreateVitalReq) []db.VitalSetParam {
	var params []db.VitalSetParam

	if input.HeightCm != nil {
		params = append(params, db.Vital.HeightCm.Set(*input.HeightCm))
	}
	if input.WeightKg != nil {
		params = append(params, db.Vital.WeightKg.Set(*input.WeightKg))
	}
	if input.BMI != nil {
		params = append(params, db.Vital.Bmi.Set(*input.BMI))
	}
	if input.TemperatureC != nil {
		params = append(params, db.Vital.TemperatureC.Set(*input.TemperatureC))
	}
	if input.Pulse != nil {
		params = append(params, db.Vital.Pulse.Set(*input.Pulse))
	}
	if input.RespiratoryRate != nil {
		params = append(params, db.Vital.RespiratoryRate.Set(*input.RespiratoryRate))
	}
	if input.BloodPressureSystolic != nil {
		params = append(params, db.Vital.BloodPressureSystolic.Set(*input.BloodPressureSystolic))
	}
	if input.BloodPressureDiastolic != nil {
		params = append(params, db.Vital.BloodPressureDiastolic.Set(*input.BloodPressureDiastolic))
	}
	if input.OxygenSaturation != nil {
		params = append(params, db.Vital.OxygenSaturation.Set(*input.OxygenSaturation))
	}

	fmt.Println(params)

	return params
}

func prepareVitalsUpdateParams(input *dto.UpdateVitalReq) []db.VitalSetParam {
	var params []db.VitalSetParam

	with := func(ok bool, p db.VitalSetParam) {
		if ok {
			params = append(params, p)
		}
	}

	if input.HeightCm != nil {
		with(*input.HeightCm >= 0, db.Vital.HeightCm.Set(*input.HeightCm))
	}
	if input.WeightKg != nil {
		with(*input.WeightKg >= 0, db.Vital.WeightKg.Set(*input.WeightKg))
	}
	if input.BMI != nil {
		with(*input.BMI >= 0, db.Vital.Bmi.Set(*input.BMI))
	}
	if input.TemperatureC != nil {
		with(*input.TemperatureC >= 30 && *input.TemperatureC <= 45, db.Vital.TemperatureC.Set(*input.TemperatureC))
	}
	if input.Pulse != nil {
		with(*input.Pulse >= 0, db.Vital.Pulse.Set(*input.Pulse))
	}
	if input.RespiratoryRate != nil {
		with(*input.RespiratoryRate >= 0, db.Vital.RespiratoryRate.Set(*input.RespiratoryRate))
	}
	if input.BloodPressureSystolic != nil {
		with(*input.BloodPressureSystolic >= 0, db.Vital.BloodPressureSystolic.Set(*input.BloodPressureSystolic))
	}
	if input.BloodPressureDiastolic != nil {
		with(*input.BloodPressureDiastolic >= 0, db.Vital.BloodPressureDiastolic.Set(*input.BloodPressureDiastolic))
	}
	if input.OxygenSaturation != nil {
		with(*input.OxygenSaturation >= 0 && *input.OxygenSaturation <= 100, db.Vital.OxygenSaturation.Set(*input.OxygenSaturation))
	}

	return params
}

func prepareAllergyUpdateParams(input *dto.UpdateAllergyReq) []db.AllergySetParam {
	var params []db.AllergySetParam

	with := func(ok bool, p db.AllergySetParam) {
		if ok {
			params = append(params, p)
		}
	}

	if input.Name != nil {
		trimmed := strings.TrimSpace(*input.Name)
		with(len(trimmed) >= 2 && len(trimmed) <= 100, db.Allergy.Name.Set(trimmed))
	}

	if input.Severity != nil {
		trimmed := strings.ToLower(strings.TrimSpace(*input.Severity))
		switch trimmed {
		case "mild", "moderate", "severe":
			with(true, db.Allergy.Severity.Set(trimmed))
		}
	}

	if input.Reaction != nil {
		trimmed := strings.TrimSpace(*input.Reaction)
		with(len(trimmed) >= 2 && len(trimmed) <= 255, db.Allergy.Reaction.Set(trimmed))
	}

	return params
}
