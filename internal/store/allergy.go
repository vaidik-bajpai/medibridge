package store

import (
	"context"

	dto "github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

type Allergy struct {
	client *db.PrismaClient
}

func (s *Allergy) Record(ctx context.Context, req *dto.RegAllergyReq) error {
	_, err := s.client.Allergy.CreateOne(
		db.Allergy.Name.Set(req.Name),
		db.Allergy.Reaction.Set(req.Reaction),
		db.Allergy.Severity.Set(req.Severity),
		db.Allergy.Patient.Link(
			db.Patient.ID.Equals(req.PatientID),
		),
	).Exec(ctx)
	return err
}

func (s *Allergy) Update(ctx context.Context, req *dto.UpdateAllergyReq) error {
	update := prepareAllergyUpdateParams(req)
	_, err := s.client.Allergy.FindUnique(
		db.Allergy.ID.Equals(req.AllergyID),
	).Update(
		update...,
	).Exec(ctx)
	return err
}

func (s *Allergy) Delete(ctx context.Context, aID string) error {
	_, err := s.client.Allergy.FindUnique(
		db.Allergy.ID.Equals(aID),
	).Delete().Exec(ctx)
	return err
}
