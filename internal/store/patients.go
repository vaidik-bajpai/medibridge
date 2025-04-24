package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/vaidik-bajpai/medibridge/internal/dto"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

var (
	ErrPatientNotFound = errors.New("patient not found")
)

type Patient struct {
	client *db.PrismaClient
}

func (s *Patient) Create(ctx context.Context, req *dto.RegPatientReq) error {
	_, err := s.client.Patient.CreateOne(
		db.Patient.FullName.Set(req.FullName),
		db.Patient.Age.Set(req.Age),
		db.Patient.Gender.Set(req.Gender),
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
	query := s.client.Patient.FindMany().OrderBy(
		db.Patient.CreatedAt.Order(db.DESC),
	)

	if req.LastID != "" {
		query = query.Cursor(db.Patient.ID.Cursor(req.LastID)).Skip(1).Take(int(req.PageSize))
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

func (s *Patient) Update(ctx context.Context, req *dto.UpdatePatientReq) error {
	update := preparePatientUpdateParams(req)

	fmt.Println("ID", req.ID)

	_, err := s.client.Patient.FindUnique(
		db.Patient.ID.Equals(req.ID),
	).Update(
		update...,
	).Exec(ctx)
	if err != nil {
		if ok := db.IsErrNotFound(err); ok {
			return ErrPatientNotFound
		}
		return err
	}
	return nil
}

func (s *Patient) Delete(ctx context.Context, pID string) error {
	_, err := s.client.Patient.FindUnique(
		db.Patient.ID.Equals(pID),
	).Delete().Exec(ctx)
	if err != nil {
		if ok := db.IsErrNotFound(err); ok {
			return ErrPatientNotFound
		}
		return err
	}
	return nil
}

func (s *Patient) Get(ctx context.Context, pID string) (*dto.Record, error) {
	patient, err := s.client.Patient.FindUnique(
		db.Patient.ID.Equals(pID),
	).With(
		db.Patient.Conditions.Fetch(),
		db.Patient.Diagnoses.Fetch(),
		db.Patient.Allergies.Fetch(),
		db.Patient.Diagnoses.Fetch(),
	).Exec(ctx)
	if err != nil {
		if ok := db.IsErrNotFound(err); ok {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}

	record := &dto.Record{
		Patient: dto.PatientModel{
			ID:                patient.ID,
			FullName:          patient.FullName,
			Age:               patient.Age,
			Gender:            patient.Gender,
			DateOfBirth:       patient.DateOfBirth,
			ContactNumber:     patient.ContactNumber,
			Address:           patient.Address,
			EmergencyName:     patient.EmergencyName,
			EmergencyRelation: patient.EmergencyRelation,
			EmergencyPhone:    patient.EmergencyPhone,
			RegisteredByID:    patient.RegisteredByID,
			CreatedAt:         patient.CreatedAt,
		},
	}

	// Map allergies
	for _, a := range patient.Allergies() {
		allergy := dto.AllergyModel{
			ID:         a.ID,
			PatientID:  a.PatientID,
			Name:       a.Name,
			Reaction:   a.Reaction,
			Severity:   a.Severity,
			RecordedAt: a.RecordedAt,
		}
		updatedAt, ok := a.UpdatedAt()
		if ok {
			allergy.UpdatedAt = &updatedAt
		}
		record.Allergies = append(record.Allergies, allergy)
	}

	// Map conditions
	for _, c := range patient.Conditions() {
		condition := dto.ConditionModel{
			ID:        c.ID,
			PatientID: c.PatientID,
			Name:      c.Name,
			CreatedAt: c.CreatedAt,
		}
		updatedAt, ok := c.UpdatedAt()
		if ok {
			condition.UpdatedAt = &updatedAt
		}
		record.Conditions = append(record.Conditions, condition)
	}

	// Map diagnoses
	for _, d := range patient.Diagnoses() {
		diagnosis := dto.DiagnosesModel{
			ID:        d.ID,
			PatientID: d.PatientID,
			Name:      d.Name,
			CreatedAt: d.CreatedAt,
		}
		updatedAt, ok := d.UpdatedAt()
		if ok {
			diagnosis.UpdatedAt = &updatedAt
		}
		record.Diagnoses = append(record.Diagnoses, diagnosis)
	}

	var vitals dto.VitalModel

	v, ok := patient.Vitals()
	if ok {
		vitals.ID = v.ID
		vitals.PatientID = v.PatientID
		vitals.CreatedAt = v.CreatedAt
		vitals.UpdatedAt = v.UpdatedAt

		if height, ok := v.HeightCm(); ok {
			vitals.HeightCm = &height
		}

		if weight, ok := v.WeightKg(); ok {
			vitals.WeightKg = &weight
		}
		if bmi, ok := v.Bmi(); ok {
			vitals.BMI = &bmi
		}
		if temperature, ok := v.TemperatureC(); ok {
			vitals.TemperatureC = &temperature
		}
		if pulse, ok := v.Pulse(); ok {
			vitals.Pulse = &pulse
		}
		if respiratoryRate, ok := v.RespiratoryRate(); ok {
			vitals.RespiratoryRate = &respiratoryRate
		}
		if systolic, ok := v.BloodPressureSystolic(); ok {
			vitals.BloodPressureSystolic = &systolic
		}
		if diastolic, ok := v.BloodPressureDiastolic(); ok {
			vitals.BloodPressureDiastolic = &diastolic
		}
		if oxygenSaturation, ok := v.OxygenSaturation(); ok {
			vitals.OxygenSaturation = &oxygenSaturation
		}
	}

	record.Vitals = vitals

	return record, nil
}
