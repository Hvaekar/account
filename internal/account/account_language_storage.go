package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddLanguage(c context.Context, accountID interface{}, req *model.AddLanguage) (*model.Language, error) {
	psql := s.SetFormat().RunWith(s.DB)

	_ = psql.Insert(accountLanguagesTableName).
		Columns(
			"account_id",
			"language",
			"level",
		).
		Values(
			accountID,
			req.Language,
			req.Level,
		).
		QueryRowContext(c)

	return s.GetLanguageByCode(c, req.Language, accountID)
}

func (s *PostgresStorage) GetLanguages(c context.Context, accountID interface{}) ([]*model.Language, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.languageResponseColumns()...).
		From(accountLanguagesTableName).
		Where("account_id = ?", accountID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var languages []*model.Language
	for rows.Next() {
		language, err := s.scanLanguage(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		languages = append(languages, language)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return languages, nil
}

func (s *PostgresStorage) GetLanguageByCode(c context.Context, languageCode interface{}, accountID interface{}) (*model.Language, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.languageResponseColumns()...).
		From(accountLanguagesTableName).
		Where("language = ? AND account_id = ?", languageCode, accountID).
		QueryRowContext(c)

	e, err := s.scanLanguage(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return e, nil
}

func (s *PostgresStorage) UpdateLanguage(c context.Context, languageCode interface{}, accountID interface{}, req *model.UpdateLanguage) (*model.Language, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(accountLanguagesTableName).
		//Set("updated_at", time.Now()).
		Set("level", req.Level).
		Where("language = ? AND account_id = ?", languageCode, accountID).
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

	return s.GetLanguageByCode(c, languageCode, accountID)
}

func (s *PostgresStorage) UpdateLanguageFields(c context.Context, languageCode interface{}, accountID interface{}, req model.UpdateLanguageFields) (*model.Language, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(accountLanguagesTableName).
		SetMap(req).
		Where("language = ? AND account_id = ?", languageCode, accountID).
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

	return s.GetLanguageByCode(c, languageCode, accountID)
}

func (s *PostgresStorage) DeleteLanguage(c context.Context, languageCode interface{}, accountID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(accountLanguagesTableName).
		Where("language = ? AND account_id = ?", languageCode, accountID).
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

func (s *PostgresStorage) languageResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "account_id",
		pre + "language",
		pre + "level",
	}
}

func (s *PostgresStorage) scanLanguage(row squirrel.RowScanner) (*model.Language, error) {
	var l model.Language

	if err := row.Scan(
		&l.AccountID,
		&l.Language,
		&l.Level,
	); err != nil {
		return nil, err
	}

	return &l, nil
}
