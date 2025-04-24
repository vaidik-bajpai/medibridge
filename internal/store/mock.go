package store

import (
	"testing"

	"github.com/vaidik-bajpai/medibridge/internal/mocks"
)

func NewMockStore(t *testing.T) *Store {
	return &Store{
		Patient:    mocks.NewPatientStorer(t),
		Session:    mocks.NewSessionStorer(t),
		User:       mocks.NewUserStorer(t),
		Diagnoses:  mocks.NewDiagnosesStorer(t),
		Vitals:     mocks.NewVitalsStorer(t),
		Conditions: mocks.NewConditionStorer(t),
		Allergy:    mocks.NewAllergyStorer(t),
	}
}
