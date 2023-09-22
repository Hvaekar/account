package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddMetalComponent(c context.Context, patientID interface{}, req *model.AddMetalComponent) (*model.MetalComponent, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(patientMetalComponentsTableName).
		Columns(
			"profile_id",
			"metal",
			"organ_id",
			"description",
		).
		Values(
			patientID,
			storage.NullString(req.Metal),
			req.OrganID,
			storage.NullString(req.Description),
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return s.GetMetalComponentByID(c, id)
}

func (s *PostgresStorage) GetMetalComponents(c context.Context, patientID interface{}) ([]*model.MetalComponent, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.patientMetalComponentsResponseColumns()...).
		From(patientMetalComponentsTableName).
		Where("profile_id = ?", patientID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var mcs []*model.MetalComponent
	for rows.Next() {
		mc, err := s.scanMetalComponent(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		mcs = append(mcs, mc)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return mcs, nil
}

func (s *PostgresStorage) GetMetalComponentByID(c context.Context, id interface{}) (*model.MetalComponent, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.patientMetalComponentsResponseColumns()...).
		From(patientMetalComponentsTableName).
		Where("id = ?", id).
		QueryRowContext(c)

	mc, err := s.scanMetalComponent(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return mc, nil
}

func (s *PostgresStorage) UpdateMetalComponent(c context.Context, id interface{}, patientID interface{}, req *model.UpdateMetalComponent) (*model.MetalComponent, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(patientMetalComponentsTableName).
		//Set("updated_at", time.Now()).
		Set("metal", storage.NullString(req.Metal)).
		Set("organ_id", req.OrganID).
		Set("description", storage.NullString(req.Description)).
		Where("id = ? AND profile_id = ?", id, patientID).
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

	return s.GetMetalComponentByID(c, id)
}

func (s *PostgresStorage) UpdateMetalComponentFields(c context.Context, id interface{}, patientID interface{}, req model.UpdateMetalComponentFields) (*model.MetalComponent, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(patientMetalComponentsTableName).
		SetMap(req).
		Where("id = ? AND profile_id = ?", id, patientID).
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

	return s.GetMetalComponentByID(c, id)
}

func (s *PostgresStorage) DeleteMetalComponent(c context.Context, id interface{}, patientID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(patientMetalComponentsTableName).
		Where("id = ? AND profile_id = ?", id, patientID).
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

func (s *PostgresStorage) patientMetalComponentsResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "profile_id",
		pre + "metal",
		pre + "organ_id",
		pre + "description",
	}
}

func (s *PostgresStorage) scanMetalComponent(row squirrel.RowScanner) (*model.MetalComponent, error) {
	var mc model.MetalComponent

	if err := row.Scan(
		&mc.ID,
		&mc.PatientID,
		&mc.Metal,
		&mc.OrganID,
		&mc.Description,
	); err != nil {
		return nil, err
	}

	return &mc, nil
}
