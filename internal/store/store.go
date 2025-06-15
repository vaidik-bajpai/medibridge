package store

import (
	"context"

	"github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

type UserStorer interface {
	Create(context.Context, *models.SignupReq) error
	FindViaEmail(ctx context.Context, email string) (*models.UserModel, error)
}

type PatientStorer interface {
	Create(context.Context, *models.RegPatientReq) (*models.Patient, error)
	Update(context.Context, *models.UpdatePatientReq) (*models.Patient, error)
	Delete(ctx context.Context, pID string) error
	List(ctx context.Context, req *models.Paginate) (*models.ListPatientRes, error)
	Get(ctx context.Context, pID string) (*models.Record, error)
}

type SessionStorer interface {
	Create(context.Context, *models.CreateSessReq) error
	FindUserByToken(ctx context.Context, token string) (*models.UserModel, error)
}

type DiagnosesStorer interface {
	Add(ctx context.Context, req *models.DiagnosesReq) error
	Update(ctx context.Context, req *models.UpdateDiagnosesReq) error
	Delete(ctx context.Context, pID string) error
}

type VitalsStorer interface {
	Create(ctx context.Context, req *models.CreateVitalReq) error
	Update(ctx context.Context, req *models.UpdateVitalReq) error
	Delete(ctx context.Context, pID string) error
}

type ConditionStorer interface {
	Add(ctx context.Context, req *models.AddConditionReq) (*models.Condition, error)
	Delete(ctx context.Context, pID string) error
}

type AllergyStorer interface {
	Record(ctx context.Context, req *models.RegAllergyReq) (*models.Allergy, error)
	Update(ctx context.Context, req *models.UpdateAllergyReq) (*models.Allergy, error)
	Delete(ctx context.Context, aID string) error
}

type Store struct {
	User       UserStorer
	Patient    PatientStorer
	Session    SessionStorer
	Diagnoses  DiagnosesStorer
	Vitals     VitalsStorer
	Conditions ConditionStorer
	Allergy    AllergyStorer
}

func NewStore(client *db.PrismaClient) *Store {
	return &Store{
		User:       &User{client: client},
		Patient:    &Patient{client: client},
		Session:    &Session{client: client},
		Diagnoses:  &Diagnoses{client: client},
		Vitals:     &Vitals{client: client},
		Conditions: &Conditions{client: client},
		Allergy:    &Allergy{client: client},
	}
}
