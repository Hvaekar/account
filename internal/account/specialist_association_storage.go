package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddSpecialistProfileAssociation(c context.Context, specialistID interface{}, req *model.AddAssociation) (*model.Association, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(specialistAssociationsTableName).
		Columns(
			"profile_id",
			"association_id",
			"name",
			"job_title",
		).
		Values(
			specialistID,
			storage.NullInt64(req.AssociationID),
			req.Name,
			storage.NullString(req.JobTitle),
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return s.GetSpecialistProfileAssociationByID(c, id)
}

func (s *PostgresStorage) GetSpecialistProfileAssociations(c context.Context, specialistID interface{}) ([]*model.Association, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.associationResponseColumns()...).
		From(specialistAssociationsTableName).
		Where("profile_id = ?", specialistID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	var as []*model.Association
	for rows.Next() {
		a, err := s.scanAssociation(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}

		as = append(as, a)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return as, nil
}

func (s *PostgresStorage) GetSpecialistProfileAssociationByID(c context.Context, id interface{}) (*model.Association, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.associationResponseColumns()...).
		From(specialistAssociationsTableName).
		Where("id = ?", id).
		QueryRowContext(c)

	a, err := s.scanAssociation(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return a, nil
}

func (s *PostgresStorage) UpdateSpecialistProfileAssociation(c context.Context, id interface{}, specialistID interface{}, req *model.UpdateAssociation) (*model.Association, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(specialistAssociationsTableName).
		Set("association_id", storage.NullInt64(req.AssociationID)).
		Set("name", req.Name).
		Set("job_title", storage.NullString(req.JobTitle)).
		Where("id = ? AND profile_id = ?", id, specialistID).
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

	return s.GetSpecialistProfileAssociationByID(c, id)
}

func (s *PostgresStorage) UpdateSpecialistProfileAssociationFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdateAssociationFields) (*model.Association, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(specialistAssociationsTableName).
		SetMap(req).
		Where("id = ? AND profile_id = ?", id, specialistID).
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

	return s.GetSpecialistProfileAssociationByID(c, id)
}

func (s *PostgresStorage) DeleteSpecialistProfileAssociation(c context.Context, id interface{}, specialistID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(specialistAssociationsTableName).
		Where("id = ? AND profile_id = ?", id, specialistID).
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

func (s *PostgresStorage) associationResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "profile_id",
		pre + "association_id",
		pre + "name",
		pre + "job_title",
	}
}

func (s *PostgresStorage) scanAssociation(row squirrel.RowScanner) (*model.Association, error) {
	var a model.Association

	if err := row.Scan(
		&a.ID,
		&a.ProfileID,
		&a.AssociationID,
		&a.Name,
		&a.JobTitle,
	); err != nil {
		return nil, err
	}

	return &a, nil
}
