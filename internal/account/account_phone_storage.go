package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddPhone(c context.Context, accountID interface{}, req *model.AddPhone) (*model.Phone, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(accountPhonesTableName).
		Columns(
			"account_id",
			"type",
			"code",
			"phone",
			"open",
		).
		Values(
			accountID,
			req.Type,
			req.Code,
			req.Phone,
			req.Open,
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return s.GetPhoneByID(c, id)
}

func (s *PostgresStorage) GetPhones(c context.Context, accountID interface{}) ([]*model.Phone, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.phoneResponseColumns()...).
		From(accountPhonesTableName).
		Where("account_id = ?", accountID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var phones []*model.Phone
	for rows.Next() {
		phone, err := s.scanPhone(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		phones = append(phones, phone)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return phones, nil
}

func (s *PostgresStorage) GetPhoneByID(c context.Context, id interface{}) (*model.Phone, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.phoneResponseColumns()...).
		From(accountPhonesTableName).
		Where("id = ?", id).
		QueryRowContext(c)

	e, err := s.scanPhone(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return e, nil
}

func (s *PostgresStorage) UpdatePhone(c context.Context, id interface{}, accountID interface{}, req *model.UpdatePhone) (*model.Phone, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(accountPhonesTableName).
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

	return s.GetPhoneByID(c, id)
}

func (s *PostgresStorage) UpdatePhoneFields(c context.Context, id interface{}, accountID interface{}, req model.UpdatePhoneFields) (*model.Phone, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(accountPhonesTableName).
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

	return s.GetPhoneByID(c, id)
}

func (s *PostgresStorage) DeletePhone(c context.Context, id interface{}, accountID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(accountPhonesTableName).
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

func (s *PostgresStorage) phoneResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "account_id",
		pre + "type",
		pre + "code",
		pre + "phone",
		pre + "verified",
		pre + "open",
	}
}

func (s *PostgresStorage) scanPhone(row squirrel.RowScanner) (*model.Phone, error) {
	var p model.Phone

	if err := row.Scan(
		&p.ID,
		&p.AccountID,
		&p.Type,
		&p.Code,
		&p.Phone,
		&p.Verified,
		&p.Open,
	); err != nil {
		return nil, err
	}

	return &p, nil
}
