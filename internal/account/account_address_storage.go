package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddAddress(c context.Context, accountID interface{}, req *model.AddAddress) (*model.Address, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(accountAddressesTableName).
		Columns(
			"account_id",
			"type",
			"city_id",
			"address",
			"open",
		).
		Values(
			accountID,
			req.Type,
			req.CityID,
			req.Address,
			req.Open,
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return s.GetAddressByID(c, id)
}

func (s *PostgresStorage) GetAddresses(c context.Context, accountID interface{}) ([]*model.Address, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.addressResponseColumns()...).
		From(accountAddressesTableName).
		Where("account_id = ?", accountID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var addresses []*model.Address
	for rows.Next() {
		address, err := s.scanAddress(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		addresses = append(addresses, address)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return addresses, nil
}

func (s *PostgresStorage) GetAddressByID(c context.Context, id interface{}) (*model.Address, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.addressResponseColumns()...).
		From(accountAddressesTableName).
		Where("id = ?", id).
		QueryRowContext(c)

	e, err := s.scanAddress(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return e, nil
}

func (s *PostgresStorage) UpdateAddress(c context.Context, id interface{}, accountID interface{}, req *model.UpdateAddress) (*model.Address, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(accountAddressesTableName).
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

	return s.GetAddressByID(c, id)
}

func (s *PostgresStorage) UpdateAddressFields(c context.Context, id interface{}, accountID interface{}, req model.UpdateAddressFields) (*model.Address, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(accountAddressesTableName).
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

	return s.GetAddressByID(c, id)
}

func (s *PostgresStorage) DeleteAddress(c context.Context, id interface{}, accountID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(accountAddressesTableName).
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

func (s *PostgresStorage) addressResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "account_id",
		pre + "type",
		pre + "city_id",
		pre + "address",
		pre + "open",
	}
}

func (s *PostgresStorage) scanAddress(row squirrel.RowScanner) (*model.Address, error) {
	var a model.Address

	if err := row.Scan(
		&a.ID,
		&a.AccountID,
		&a.Type,
		&a.CityID,
		&a.Address,
		&a.Open,
	); err != nil {
		return nil, err
	}

	return &a, nil
}
