package store

import (
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
