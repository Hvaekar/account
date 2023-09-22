package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddEmail(c context.Context, accountID interface{}, req *model.AddEmail) (*model.Email, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(accountEmailsTableName).
		Columns(
			"account_id",
			"type",
			"email",
			"open",
		).
		Values(
			accountID,
			req.Type,
			req.Email,
			req.Open,
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return s.GetEmailByID(c, id)
}

func (s *PostgresStorage) GetEmails(c context.Context, accountID interface{}) ([]*model.Email, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.emailResponseColumns()...).
		From(accountEmailsTableName).
		Where("account_id = ?", accountID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var emails []*model.Email
	for rows.Next() {
		email, err := s.scanEmail(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		emails = append(emails, email)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return emails, nil
}

func (s *PostgresStorage) GetEmailByID(c context.Context, id interface{}) (*model.Email, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.emailResponseColumns()...).
		From(accountEmailsTableName).
		Where("id = ?", id).
		QueryRowContext(c)

	e, err := s.scanEmail(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return e, nil
}

func (s *PostgresStorage) UpdateEmail(c context.Context, id interface{}, accountID interface{}, req *model.UpdateEmail) (*model.Email, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(accountEmailsTableName).
		//Set("updated_at", time.Now()).
		Set("type", req.Type).
		Set("open", req.Open).
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

	return s.GetEmailByID(c, id)
}

func (s *PostgresStorage) UpdateEmailFields(c context.Context, id interface{}, accountID interface{}, req model.UpdateEmailFields) (*model.Email, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(accountEmailsTableName).
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

	return s.GetEmailByID(c, id)
}

func (s *PostgresStorage) DeleteEmail(c context.Context, id interface{}, accountID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(accountEmailsTableName).
		Where("id = ? AND account_id = ?", id, accountID).
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

func (s *PostgresStorage) emailResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "account_id",
		pre + "type",
		pre + "email",
		pre + "verified",
		pre + "open",
	}
}

func (s *PostgresStorage) scanEmail(row squirrel.RowScanner) (*model.Email, error) {
	var e model.Email

	if err := row.Scan(
		&e.ID,
		&e.AccountID,
		&e.Type,
		&e.Email,
		&e.Verified,
		&e.Open,
	); err != nil {
		return nil, err
	}

	return &e, nil
}
