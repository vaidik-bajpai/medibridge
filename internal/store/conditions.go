package store

import (
	"context"

	dto "github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

type Conditions struct {
	client *db.PrismaClient
}

func (s *Conditions) Add(ctx context.Context, req *dto.AddConditionReq) error {
	_, err := s.client.Condition.CreateOne(
		db.Condition.Patient.Link(
			db.Patient.ID.Equals(req.PatientID),
		),
		db.Condition.Name.Set(req.Condition),
	).Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Conditions) Delete(ctx context.Context, pID string) error {
	_, err := s.client.Condition.FindUnique(
		db.Condition.ID.Equals(pID),
	).Delete().Exec(ctx)
	if err != nil {
		return err
	}

	return nil
}
