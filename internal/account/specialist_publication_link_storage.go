package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddSpecialistProfilePublicationLink(c context.Context, specialistID interface{}, req *model.AddPublicationLink) (*model.PublicationLink, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(specialistPublicationLinksTableName).
		Columns(
			"profile_id",
			"title",
			"link",
		).
		Values(
			specialistID,
			req.Title,
			req.Link,
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return s.GetSpecialistProfilePublicationLinkByID(c, id)
}

func (s *PostgresStorage) GetSpecialistProfilePublicationLinks(c context.Context, specialistID interface{}) ([]*model.PublicationLink, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.publicationLinkResponseColumns()...).
		From(specialistPublicationLinksTableName).
		Where("profile_id = ?", specialistID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	var as []*model.PublicationLink
	for rows.Next() {
		a, err := s.scanPublicationLink(rows)
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

func (s *PostgresStorage) GetSpecialistProfilePublicationLinkByID(c context.Context, id interface{}) (*model.PublicationLink, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.publicationLinkResponseColumns()...).
		From(specialistPublicationLinksTableName).
		Where("id = ?", id).
		QueryRowContext(c)

	a, err := s.scanPublicationLink(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return a, nil
}

func (s *PostgresStorage) UpdateSpecialistProfilePublicationLink(c context.Context, id interface{}, specialistID interface{}, req *model.UpdatePublicationLink) (*model.PublicationLink, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(specialistPublicationLinksTableName).
		Set("title", req.Title).
		Set("link", req.Link).
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

	return s.GetSpecialistProfilePublicationLinkByID(c, id)
}

func (s *PostgresStorage) UpdateSpecialistProfilePublicationLinkFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdatePublicationLinkFields) (*model.PublicationLink, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(specialistPublicationLinksTableName).
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

	return s.GetSpecialistProfilePublicationLinkByID(c, id)
}

func (s *PostgresStorage) DeleteSpecialistProfilePublicationLink(c context.Context, id interface{}, specialistID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(specialistPublicationLinksTableName).
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

func (s *PostgresStorage) publicationLinkResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "profile_id",
		pre + "title",
		pre + "link",
	}
}

func (s *PostgresStorage) scanPublicationLink(row squirrel.RowScanner) (*model.PublicationLink, error) {
	var a model.PublicationLink

	if err := row.Scan(
		&a.ID,
		&a.ProfileID,
		&a.Title,
		&a.Link,
	); err != nil {
		return nil, err
	}

	return &a, nil
}
