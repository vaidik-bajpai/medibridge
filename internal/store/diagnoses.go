package store

import (
	"context"

	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

type Diagnoses struct {
	client *db.PrismaClient
}

func (s *Diagnoses) Add(ctx context.Context, req *dto.DiagnosesReq) error {
	_, err := s.client.Diagnosis.CreateOne(
		db.Diagnosis.Patient.Link(
			db.Patient.ID.Equals(req.PID),
		),
		db.Diagnosis.Name.Set(req.Name),
	).Exec(ctx)
	if err != nil {
		if ok := db.IsErrNotFound(err); ok {
			return ErrPatientNotFound
		}
		return err
	}
	return nil
}

func (s *Diagnoses) Update(ctx context.Context, req *dto.UpdateDiagnosesReq) error {
	_, err := s.client.Diagnosis.FindUnique(
		db.Diagnosis.ID.Equals(req.DID),
	).Update(
		db.Diagnosis.Name.Set(req.Name),
	).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Diagnoses) Delete(ctx context.Context, pID string) error {
	_, err := s.client.Diagnosis.FindUnique(
		db.Diagnosis.ID.Equals(pID),
	).Delete().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
