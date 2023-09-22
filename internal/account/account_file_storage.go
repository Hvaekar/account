package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
	"time"
)

func (s *PostgresStorage) AddFile(c context.Context, accountID interface{}, fileName *string) (*model.File, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(accountFilesTableName).
		Columns("account_id", "name").
		Values(accountID, fileName).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return s.GetFileByID(c, id)
}

func (s *PostgresStorage) GetFiles(c context.Context, accountID interface{}) ([]*model.File, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.fileResponseColumns()...).
		From(accountFilesTableName).
		Where("account_id = ?", accountID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var files []*model.File
	for rows.Next() {
		file, err := s.scanFile(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return files, nil
}

func (s *PostgresStorage) GetPatientDisabilityFiles(c context.Context, patientID interface{}) ([]*model.File, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.fileResponseColumns(accountFilesTableName)...).
		From(patientDisabilityFilesTableName).
		LeftJoin(accountFilesTableName+" ON "+accountFilesTableName+".id = "+patientDisabilityFilesTableName+".file_id").
		Where(patientDisabilityFilesTableName+".profile_id = ?", patientID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var files []*model.File
	for rows.Next() {
		file, err := s.scanFile(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return files, nil
}

func (s *PostgresStorage) GetFileByID(c context.Context, id interface{}) (*model.File, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.fileResponseColumns()...).
		From(accountFilesTableName).
		Where("id = ?", id).
		QueryRowContext(c)

	f, err := s.scanFile(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return f, nil
}

func (s *PostgresStorage) GetFileByName(c context.Context, name interface{}) (*model.File, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.fileResponseColumns()...).
		From(accountFilesTableName).
		Where("name = ?", name).
		QueryRowContext(c)

	f, err := s.scanFile(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return f, nil
}

func (s *PostgresStorage) UpdateFile(c context.Context, id interface{}, accountID interface{}, req *model.UpdateFile) (*model.File, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(accountFilesTableName).
		Set("updated_at", time.Now()).
		Set("description", storage.NullString(req.Description)).
		Where("id = ? AND account_id = ?", id, accountID).
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

	return s.GetFileByID(c, id)
}

func (s *PostgresStorage) UpdateFileFields(c context.Context, id interface{}, accountID interface{}, req model.UpdateFileFields) (*model.File, error) {
	psql := s.SetFormat().RunWith(s.DB)
	req["updated_at"] = time.Now()

	res, err := psql.Update(accountFilesTableName).
		SetMap(req).
		Where("id = ? AND account_id = ?", id, accountID).
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

	return s.GetFileByID(c, id)
}

func (s *PostgresStorage) DeleteFile(c context.Context, id interface{}, accountID interface{}) (*string, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Delete(accountFilesTableName).
		Where("id = ? AND account_id = ?", id, accountID).
		Suffix("RETURNING \"name\"")

	var name string
	if err := q.QueryRowContext(c).Scan(&name); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return &name, nil
}

func (s *PostgresStorage) fileResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "created_at",
		pre + "updated_at",
		pre + "account_id",
		pre + "name",
		pre + "description",
	}
}

func (s *PostgresStorage) scanFile(row squirrel.RowScanner) (*model.File, error) {
	var f model.File

	if err := row.Scan(
		&f.ID,
		&f.CreatedAt,
		&f.UpdatedAt,
		&f.AccountID,
		&f.Name,
		&f.Description,
	); err != nil {
		return nil, err
	}

	return &f, nil
}
