package store

import (
	"context"

	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

type Patient struct {
	client *db.PrismaClient
}

func (s *Patient) Create(ctx context.Context, req *dto.RegPatientReq) error {
	_, err := s.client.Patient.CreateOne(
		db.Patient.FullName.Set(req.FullName),
		db.Patient.Gender.Set(db.Gender(req.Gender)),
		db.Patient.DateOfBirth.Set(req.DOB),
		db.Patient.ContactNumber.Set(req.ContactNumber),
		db.Patient.Address.Set(req.Address),
		db.Patient.EmergencyName.Set(req.EmergencyName),
		db.Patient.EmergencyRelation.Set(req.EmergencyRelation),
		db.Patient.EmergencyPhone.Set(req.EmergencyPhone),
		db.Patient.RegisteredBy.Link(
			db.User.ID.Equals(req.RegByID),
		),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
