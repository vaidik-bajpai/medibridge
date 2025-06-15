package store

import (
	"context"

	"github.com/vaidik-bajpai/medibridge/internal/models"
	dto "github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

type Allergy struct {
	client *db.PrismaClient
}

func (s *Allergy) Record(ctx context.Context, req *dto.RegAllergyReq) (*models.Allergy, error) {
	allergy, err := s.client.Allergy.CreateOne(
		db.Allergy.Name.Set(req.Name),
		db.Allergy.Reaction.Set(req.Reaction),
		db.Allergy.Severity.Set(req.Severity),
		db.Allergy.Patient.Link(
			db.Patient.ID.Equals(req.PatientID),
		),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	return &models.Allergy{
		ID:        allergy.ID,
		PatientID: allergy.PatientID,
		Name:      allergy.Name,
		Severity:  allergy.Severity,
		Reaction:  allergy.Reaction,
		CreatedAt: allergy.RecordedAt,
	}, nil
}

func (s *Allergy) Update(ctx context.Context, req *dto.UpdateAllergyReq) (*models.Allergy, error) {
	update := prepareAllergyUpdateParams(req)

	allergy, err := s.client.Allergy.FindUnique(
		db.Allergy.ID.Equals(req.AllergyID),
	).Update(
		update...,
	).Exec(ctx)
	if err != nil {
		return nil, err
	}
	return &models.Allergy{
		ID:        allergy.ID,
		PatientID: allergy.PatientID,
		Name:      allergy.Name,
		Severity:  allergy.Severity,
		Reaction:  allergy.Reaction,
		CreatedAt: allergy.RecordedAt,
	}, err
}

func (s *Allergy) Delete(ctx context.Context, aID string) error {
	_, err := s.client.Allergy.FindUnique(
		db.Allergy.ID.Equals(aID),
	).Delete().Exec(ctx)
	return err
}
