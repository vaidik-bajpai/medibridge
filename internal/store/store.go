package store

import (
	"context"

	dto "github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

type UserStorer interface {
	Create(context.Context, *dto.SignupReq) error
	FindViaEmail(ctx context.Context, email string) (*dto.UserModel, error)
}

type PatientStorer interface {
	Create(context.Context, *dto.RegPatientReq) error
	Update(context.Context, *dto.UpdatePatientReq) error
	Delete(ctx context.Context, pID string) error
	List(ctx context.Context, req *dto.Paginate) ([]*dto.PatientListItem, error)
	Get(ctx context.Context, pID string) (*dto.Record, error)
}

type SessionStorer interface {
	Create(context.Context, *dto.CreateSessReq) error
	FindUserByToken(ctx context.Context, token string) (*dto.UserModel, error)
}

type DiagnosesStorer interface {
	Add(ctx context.Context, req *dto.DiagnosesReq) error
	Update(ctx context.Context, req *dto.UpdateDiagnosesReq) error
	Delete(ctx context.Context, pID string) error
}

type VitalsStorer interface {
	Create(ctx context.Context, req *dto.CreateVitalReq) error
	Update(ctx context.Context, req *dto.UpdateVitalReq) error
	Delete(ctx context.Context, pID string) error
}

type ConditionStorer interface {
	Add(ctx context.Context, req *dto.AddConditionReq) error
	Delete(ctx context.Context, pID string) error
}

type AllergyStorer interface {
	Record(ctx context.Context, req *dto.RegAllergyReq) error
	Update(ctx context.Context, req *dto.UpdateAllergyReq) error
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
