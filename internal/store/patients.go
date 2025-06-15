package store

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/vaidik-bajpai/medibridge/internal/models"
	"github.com/vaidik-bajpai/medibridge/internal/prisma/db"
)

var (
	ErrPatientNotFound = errors.New("patient not found")
)

type Patient struct {
	client *db.PrismaClient
}

func (s *Patient) Create(ctx context.Context, req *models.RegPatientReq) (*models.Patient, error) {
	p, err := s.client.Patient.CreateOne(
		db.Patient.FullName.Set(req.FullName),
		db.Patient.Age.Set(req.Age),
		db.Patient.Gender.Set(req.Gender),
		db.Patient.DateOfBirth.Set(time.Time(req.DOB)),
		db.Patient.ContactNumber.Set(req.ContactNumber),
		db.Patient.Address.Set(req.Address),
		db.Patient.EmergencyName.Set(req.EmergencyName),
		db.Patient.EmergencyRelation.Set(req.EmergencyRelation),
		db.Patient.EmergencyPhone.Set(req.EmergencyPhone),
		db.Patient.RegisteredBy.Link(
			db.User.ID.Equals(req.RegByID),
		),
	).With(
		db.Patient.RegisteredBy.Fetch(),
	).Exec(ctx)
	if err != nil {
		return nil, err
	}

	patient := models.Patient{
		ID:                p.ID,
		FullName:          p.FullName,
		Gender:            p.Gender,
		DOB:               models.DateOnly(p.DateOfBirth),
		Age:               p.Age,
		ContactNumber:     p.ContactNumber,
		Address:           p.Address,
		EmergencyName:     p.EmergencyName,
		EmergencyPhone:    p.EmergencyPhone,
		EmergencyRelation: p.EmergencyRelation,
		RegByID:           p.RegisteredByID,
		Registrar:         p.RegisteredBy().Fullname,
	}

	return &patient, err
}

func (s *Patient) List(ctx context.Context, req *models.Paginate) (*models.ListPatientRes, error) {
	offset := (req.Page - 1) * req.PageSize

	query := `
		SELECT
			id,
			"fullName",
			gender,
			age,
			"dateOfBirth" AS dob,
			COUNT(*) OVER() AS "totalCount"
		FROM
			"Patient"
		WHERE
			"fullName" ILIKE $1
		ORDER BY
			"createdAt" DESC
		LIMIT $2 OFFSET $3;
	`

	var args []interface{}
	args = append(args, "%"+req.SearchTerm+"%")
	args = append(args, req.PageSize)
	args = append(args, offset)

	var queryRes []struct {
		models.ListPatientItem
		TotalCount string `json:"totalCount"`
	}

	err := s.client.Prisma.QueryRaw(query, args...).Exec(ctx, &queryRes)
	if err != nil {
		log.Println("=====ERROR=====")
		return nil, err
	}

	log.Println("=====QUERYRES=====")
	log.Println(queryRes)

	if len(queryRes) == 0 {
		return nil, nil
	}

	res := &models.ListPatientRes{}
	totalItems, err := strconv.ParseInt(queryRes[0].TotalCount, 10, 64)
	if err != nil {
		return nil, err
	}

	for _, p := range queryRes {
		res.Patients = append(res.Patients, &models.ListPatientItem{
			ID:       p.ID,
			FullName: p.FullName,
			Gender:   p.Gender,
			Age:      p.Age,
			DOB:      p.DOB,
		})
	}

	totalPages := (totalItems + req.PageSize - 1) / req.PageSize
	from := offset + 1
	to := offset + int64(len(queryRes))
	if totalItems == 0 {
		from = 0
		to = 0
	}

	res.Meta = &models.ListPatientMetadata{
		CurrentPage: req.Page,
		PageSize:    req.PageSize,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		From:        from,
		To:          to,
		HasNext:     req.Page < totalPages,
		HasPrevious: req.Page != 1,
	}

	return res, nil
}

func (s *Patient) Update(ctx context.Context, req *models.UpdatePatientReq) (*models.Patient, error) {
	update := preparePatientUpdateParams(req)

	fmt.Println("ID", req.ID)

	p, err := s.client.Patient.FindUnique(
		db.Patient.ID.Equals(req.ID),
	).With(
		db.Patient.RegisteredBy.Fetch(),
	).Update(
		update...,
	).Exec(ctx)
	if err != nil {
		if ok := db.IsErrNotFound(err); ok {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}

	patient := models.Patient{
		ID:                p.ID,
		FullName:          p.FullName,
		Gender:            p.Gender,
		DOB:               models.DateOnly(p.DateOfBirth),
		Age:               p.Age,
		ContactNumber:     p.ContactNumber,
		Address:           p.Address,
		EmergencyName:     p.EmergencyName,
		EmergencyPhone:    p.EmergencyPhone,
		EmergencyRelation: p.EmergencyRelation,
		RegByID:           p.RegisteredByID,
		Registrar:         p.RegisteredBy().Fullname,
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         &p.UpdatedAt,
		Version:           p.Version,
	}
	return &patient, nil
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

func (s *Patient) Get(ctx context.Context, pID string) (*models.Record, error) {
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

	var updatedAt *time.Time
	if patient.Version == 0 {
		updatedAt = &patient.UpdatedAt
	}

	record := &models.Record{
		Patient: models.Patient{
			ID:                patient.ID,
			FullName:          patient.FullName,
			Age:               patient.Age,
			Gender:            patient.Gender,
			DOB:               models.DateOnly(patient.DateOfBirth),
			ContactNumber:     patient.ContactNumber,
			Address:           patient.Address,
			EmergencyName:     patient.EmergencyName,
			EmergencyRelation: patient.EmergencyRelation,
			EmergencyPhone:    patient.EmergencyPhone,
			RegByID:           patient.RegisteredByID,
			CreatedAt:         patient.CreatedAt,
			UpdatedAt:         updatedAt,
			Version:           patient.Version,
		},
	}

	for _, a := range patient.Allergies() {
		allergy := models.AllergyModel{
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

	for _, c := range patient.Conditions() {
		condition := models.ConditionModel{
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

	for _, d := range patient.Diagnoses() {
		diagnosis := models.DiagnosesModel{
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

	var vitals models.VitalModel

	v, ok := patient.Vitals()
	if ok {
		vitals.ID = v.ID
		vitals.PatientID = v.PatientID
		vitals.CreatedAt = &v.CreatedAt
		vitals.UpdatedAt = &v.UpdatedAt

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
