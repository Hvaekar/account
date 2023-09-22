package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddSpecialistProfileSpecialization(c context.Context, specialistID interface{}, req *model.AddSpecialization) (*model.Specialization, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Insert(specialistSpecializationsTableName).
		Columns("profile_id", "specialization_id", "start").
		Values(specialistID, req.SpecializationID, req.Start).
		ExecContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	if ra == 0 {
		return nil, storage.ErrNotFound
	}

	return s.GetSpecialistProfileSpecialization(c, req.SpecializationID, specialistID)
}

func (s *PostgresStorage) GetSpecialistProfileSpecializations(c context.Context, specialistID interface{}) ([]*model.Specialization, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.specializationResponseColumns()...).
		From(specialistSpecializationsTableName).
		Where("profile_id = ?", specialistID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	var specializations []*model.Specialization
	for rows.Next() {
		sp, err := s.scanSpecialization(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		specializations = append(specializations, sp)
	}

	return specializations, nil
}

func (s *PostgresStorage) GetSpecialistProfileSpecialization(c context.Context, id interface{}, specialistID interface{}) (*model.Specialization, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.specializationResponseColumns()...).
		From(specialistSpecializationsTableName).
		Where("profile_id = ? AND specialization_id = ?", specialistID, id).
		QueryRowContext(c)

	sp, err := s.scanSpecialization(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return sp, nil
}

func (s *PostgresStorage) UpdateSpecialistProfileSpecialization(c context.Context, id interface{}, specialistID interface{}, req *model.UpdateSpecialization) (*model.Specialization, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(specialistSpecializationsTableName).
		Set("start", req.Start).
		Where("profile_id = ? AND specialization_id = ?", specialistID, id).
		ExecContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	if ra == 0 {
		return nil, storage.ErrNotFound
	}

	return s.GetSpecialistProfileSpecialization(c, id, specialistID)
}

func (s *PostgresStorage) UpdateSpecialistProfileSpecializationFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdateSpecializationFields) (*model.Specialization, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(specialistSpecializationsTableName).
		SetMap(req).
		Where("id = ?", id).
		ExecContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	if ra == 0 {
		return nil, storage.ErrNotFound
	}

	return s.GetSpecialistProfileSpecialization(c, id, specialistID)
}

func (s *PostgresStorage) DeleteSpecialistProfileSpecialization(c context.Context, id interface{}, specialistID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(specialistSpecializationsTableName).
		Where("profile_id = ? AND specialization_id = ?", specialistID, id).
		ExecContext(c)
	if err != nil {
		return postgres.ConvertError(err)
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return postgres.ConvertError(err)
	}
	if ra == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (s *PostgresStorage) specializationResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "specialization_id",
		pre + "start",
	}
}

func (s *PostgresStorage) scanSpecialization(row squirrel.RowScanner) (*model.Specialization, error) {
	var sp model.Specialization

	if err := row.Scan(
		&sp.SpecializationID,
		&sp.Start,
	); err != nil {
		return nil, err
	}

	return &sp, nil
}
