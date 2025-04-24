package store

import (
	"context"
	"errors"

	dto "github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

var (
	ErrUniqueConstraintViolated = errors.New("unique constraint violated")
)

type Vitals struct {
	client *db.PrismaClient
}

func (s *Vitals) Create(ctx context.Context, req *dto.CreateVitalReq) error {
	create := prepareVitalCreateParams(req)
	_, err := s.client.Vital.CreateOne(
		db.Vital.Patient.Link(
			db.Patient.ID.Equals(req.PatientID),
		),
		create...,
	).Exec(ctx)
	if err != nil {
		if _, ok := db.IsErrUniqueConstraint(err); ok {
			return ErrUniqueConstraintViolated
		}
		return err
	}
	return nil
}

func (s *Vitals) Update(ctx context.Context, req *dto.UpdateVitalReq) error {
	update := prepareVitalsUpdateParams(req)
	_, err := s.client.Vital.FindUnique(
		db.Vital.PatientID.Equals(req.PatientID),
	).Update(
		update...,
	).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Vitals) Delete(ctx context.Context, pID string) error {
	_, err := s.client.Vital.FindUnique(
		db.Vital.PatientID.Equals(pID),
	).Delete().Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
