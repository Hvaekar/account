package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddSpecialistProfilePatent(c context.Context, specialistID interface{}, req *model.AddPatent) (*model.Patent, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(specialistPatentsTableName).
		Columns(
			"profile_id",
			"number",
			"name",
			"link",
		).
		Values(
			specialistID,
			req.Number,
			req.Name,
			storage.NullString(req.Link),
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return s.GetSpecialistProfilePatentByID(c, id)
}

func (s *PostgresStorage) GetSpecialistProfilePatents(c context.Context, specialistID interface{}) ([]*model.Patent, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.patentResponseColumns()...).
		From(specialistPatentsTableName).
		Where("profile_id = ?", specialistID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	var as []*model.Patent
	for rows.Next() {
		a, err := s.scanPatent(rows)
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

func (s *PostgresStorage) GetSpecialistProfilePatentByID(c context.Context, id interface{}) (*model.Patent, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.patentResponseColumns()...).
		From(specialistPatentsTableName).
		Where("id = ?", id).
		QueryRowContext(c)

	a, err := s.scanPatent(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return a, nil
}

func (s *PostgresStorage) UpdateSpecialistProfilePatent(c context.Context, id interface{}, specialistID interface{}, req *model.UpdatePatent) (*model.Patent, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(specialistPatentsTableName).
		Set("number", req.Number).
		Set("name", req.Name).
		Set("link", storage.NullString(req.Link)).
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

	return s.GetSpecialistProfilePatentByID(c, id)
}

func (s *PostgresStorage) UpdateSpecialistProfilePatentFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdatePatentFields) (*model.Patent, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(specialistPatentsTableName).
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

	return s.GetSpecialistProfilePatentByID(c, id)
}

func (s *PostgresStorage) DeleteSpecialistProfilePatent(c context.Context, id interface{}, specialistID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(specialistPatentsTableName).
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

func (s *PostgresStorage) patentResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "profile_id",
		pre + "number",
		pre + "name",
		pre + "link",
	}
}

func (s *PostgresStorage) scanPatent(row squirrel.RowScanner) (*model.Patent, error) {
	var a model.Patent

	if err := row.Scan(
		&a.ID,
		&a.ProfileID,
		&a.Number,
		&a.Name,
		&a.Link,
	); err != nil {
		return nil, err
	}

	return &a, nil
}
