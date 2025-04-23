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
		db.Patient.Age.Set(req.Age),
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

func (s *Patient) List(ctx context.Context, req *dto.Paginate) ([]*dto.PatientListItem, error) {
	query := s.client.Patient.FindMany(
		db.Patient.FullName.Contains((req.SearchTerm)),
		db.Patient.FullName.Mode(db.QueryModeInsensitive),
	).OrderBy(
		db.Patient.CreatedAt.Order(db.DESC),
	)

	if req.LastID != "" {
		query = query.Cursor(db.Patient.ID.Cursor(req.LastID))
	}

	list, err := query.Exec(ctx)
	if err != nil {
		return nil, err
	}

	var res []*dto.PatientListItem
	for _, patient := range list {
		res = append(res, &dto.PatientListItem{
			ID:       patient.ID,
			Username: patient.FullName,
			Gender:   string(patient.Gender),
			Age:      patient.Age,
			DOB:      patient.DateOfBirth,
		})
	}

	return res, nil
}
