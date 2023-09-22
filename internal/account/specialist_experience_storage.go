package account

import (
	"context"
	"database/sql"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
)

func (s *PostgresStorage) AddSpecialistProfileExperience(c context.Context, specialistID interface{}, req *model.AddExperience) (*model.Experience, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(specialistExperiencesTableName).
		Columns(
			"profile_id",
			"company_id",
			"company",
			"start",
			"finish",
		).
		Values(
			specialistID,
			storage.NullInt64(req.CompanyID),
			req.Company,
			req.Start,
			storage.NullDatePGX(req.Finish),
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	if err := s.UpdateSpecialistProfileExperienceSpecializations(id, req.Specializations); err != nil {
		return nil, err
	}

	return s.GetSpecialistProfileExperienceByID(c, id)
}

func (s *PostgresStorage) GetSpecialistProfileExperiences(c context.Context, specialistID interface{}) ([]*model.Experience, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.experienceResponseColumns()...).
		From(specialistExperiencesTableName).
		LeftJoin(specialistExperienceSpecializationsTableName+" ON "+specialistExperienceSpecializationsTableName+".experience_id = "+specialistExperiencesTableName+".id").
		Where(specialistExperiencesTableName+".profile_id = ?", specialistID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	e, err := s.scanExperiences(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return e, nil
}

func (s *PostgresStorage) GetSpecialistProfileExperienceByID(c context.Context, id interface{}) (*model.Experience, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.experienceResponseColumns()...).
		From(specialistExperiencesTableName).
		LeftJoin(specialistExperienceSpecializationsTableName+" ON "+specialistExperienceSpecializationsTableName+".experience_id = "+specialistExperiencesTableName+".id").
		Where(specialistExperiencesTableName+".id = ?", id).
		QueryContext(c)

	e, err := s.scanExperiences(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return e[0], nil
}

func (s *PostgresStorage) UpdateSpecialistProfileExperience(c context.Context, id interface{}, specialistID interface{}, req *model.UpdateExperience) (*model.Experience, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(specialistExperiencesTableName).
		Set("company_id", storage.NullInt64(req.CompanyID)).
		Set("company", req.Company).
		Set("start", req.Start).
		Set("finish", storage.NullDatePGX(req.Finish)).
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

	if err := s.UpdateSpecialistProfileExperienceSpecializations(id, req.Specializations); err != nil {
		return nil, err
	}

	return s.GetSpecialistProfileExperienceByID(c, id)
}

func (s *PostgresStorage) UpdateSpecialistProfileExperienceSpecializations(experienceID interface{}, req []int64) error {
	psql := s.SetFormat().RunWith(s.DB)

	_, err := psql.Delete(specialistExperienceSpecializationsTableName).Where("experience_id = ?", experienceID).Exec()
	if err != nil {
		return postgres.ConvertError(err)
	}

	if len(req) > 0 {
		aq := psql.Insert(specialistExperienceSpecializationsTableName).Columns("experience_id", "specialization_id")

		for _, v := range req {
			aq = aq.Values(experienceID, v)
		}

		res, err := aq.Exec()
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
	}

	return nil
}

func (s *PostgresStorage) UpdateSpecialistProfileExperienceFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdateExperienceFields) (*model.Experience, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(specialistExperiencesTableName).
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

	return s.GetSpecialistProfileExperienceByID(c, id)
}

func (s *PostgresStorage) DeleteSpecialistProfileExperience(c context.Context, id interface{}, specialistID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(specialistExperiencesTableName).
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

func (s *PostgresStorage) experienceResponseColumns() []string {
	fields := make([]string, 0)

	fields = append(s.experienceResponseColumnsMain(specialistExperiencesTableName),
		s.experienceSpecializationResponseColumns(specialistExperienceSpecializationsTableName)...)

	return fields
}

func (s *PostgresStorage) experienceResponseColumnsMain(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "profile_id",
		pre + "company_id",
		pre + "company",
		pre + "start",
		pre + "finish",
	}
}

func (s *PostgresStorage) experienceSpecializationResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "specialization_id",
	}
}

func (s *PostgresStorage) scanExperiences(rows *sql.Rows) ([]*model.Experience, error) {
	var es []*model.Experience
	experiences := make(map[int64]*model.Experience)

	for rows.Next() {
		var e model.Experience
		var sp *int64

		if err := rows.Scan(
			&e.ID,
			&e.ProfileID,
			&e.CompanyID,
			&e.Company,
			&e.Start,
			&e.Finish,

			&sp,
		); err != nil {
			return nil, err
		}

		if _, ok := experiences[e.ID]; !ok {
			experiences[e.ID] = &e
		}

		if sp != nil {
			experiences[e.ID].Specializations = append(experiences[e.ID].Specializations, *sp)
		}
	}

	for _, v := range experiences {
		es = append(es, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return es, nil
}
