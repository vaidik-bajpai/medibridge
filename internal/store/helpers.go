package store

import (
	"strings"

	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

func preparePatientUpdateParams(input *dto.UpdatePatientReq) []db.PatientSetParam {
	var params []db.PatientSetParam

	with := func(ok bool, p db.PatientSetParam) {
		if ok {
			params = append(params, p)
		}
	}

	with(input.FullName != "", db.Patient.FullName.Set(input.FullName))
	with(input.Gender != "", db.Patient.Gender.Set(db.Gender(input.Gender)))
	with(!input.DOB.IsZero(), db.Patient.DateOfBirth.Set(input.DOB))
	with(input.Age != 0, db.Patient.Age.Set(input.Age))

	with(input.ContactNumber != "", db.Patient.ContactNumber.Set(input.ContactNumber))
	with(input.Address != "", db.Patient.Address.Set(input.Address))

	with(input.EmergencyName != "", db.Patient.EmergencyName.Set(input.EmergencyName))
	with(input.EmergencyRelation != "", db.Patient.EmergencyRelation.Set(input.EmergencyRelation))
	with(input.EmergencyPhone != "", db.Patient.EmergencyPhone.Set(input.EmergencyPhone))

	return params
}

func prepareVitalCreateParams(input *dto.CreateVitalReq) []db.VitalSetParam {
	var params []db.VitalSetParam

	with := func(ok bool, p db.VitalSetParam) {
		if ok {
			params = append(params, p)
		}
	}

	with(input.HeightCm >= 0, db.Vital.HeightCm.Set(input.HeightCm))
	with(input.WeightKg >= 0, db.Vital.WeightKg.Set(input.WeightKg))
	with(input.BMI >= 0, db.Vital.Bmi.Set(input.BMI))
	with(input.TemperatureC >= 30 && input.TemperatureC <= 45, db.Vital.TemperatureC.Set(input.TemperatureC))
	with(input.Pulse >= 0, db.Vital.Pulse.Set(input.Pulse))
	with(input.RespiratoryRate >= 0, db.Vital.RespiratoryRate.Set(input.RespiratoryRate))
	with(input.BloodPressureSystolic >= 0, db.Vital.BloodPressureSystolic.Set(input.BloodPressureSystolic))
	with(input.BloodPressureDiastolic >= 0, db.Vital.BloodPressureDiastolic.Set(input.BloodPressureDiastolic))
	with(input.OxygenSaturation >= 0 && input.OxygenSaturation <= 100, db.Vital.OxygenSaturation.Set(input.OxygenSaturation))

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
