package account

import (
	"context"
	"database/sql"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
)

func (s *PostgresStorage) AddSpecialistProfileEducation(c context.Context, specialistID interface{}, req *model.AddEducation) (*model.Education, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(specialistEducationsTableName).
		Columns(
			"profile_id",
			"institution_id",
			"faculty_id",
			"department_id",
			"form_id",
			"degree_id",
			"graduation",
		).
		Values(
			specialistID,
			req.InstitutionID,
			storage.NullInt64(req.FacultyID),
			storage.NullInt64(req.DepartmentID),
			storage.NullInt64(req.FormID),
			storage.NullInt64(req.DegreeID),
			req.Graduation,
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	if err := s.UpdateSpecialistProfileEducationFiles(id, req.Files); err != nil {
		return nil, err
	}

	return s.GetSpecialistProfileEducationByID(c, id)
}

func (s *PostgresStorage) GetSpecialistProfileEducations(c context.Context, specialistID interface{}) ([]*model.Education, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.educationResponseColumns()...).
		From(specialistEducationsTableName).
		LeftJoin(specialistEducationFilesTableName+" ON "+specialistEducationFilesTableName+".education_id = "+specialistEducationsTableName+".id").
		LeftJoin(accountFilesTableName+" ON "+accountFilesTableName+".id = "+specialistEducationFilesTableName+".file_id").
		Where(specialistEducationsTableName+".profile_id = ?", specialistID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	educations, err := s.scanEducations(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return educations, nil
}

func (s *PostgresStorage) GetSpecialistProfileEducationByID(c context.Context, id interface{}) (*model.Education, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.educationResponseColumns()...).
		From(specialistEducationsTableName).
		LeftJoin(specialistEducationFilesTableName+" ON "+specialistEducationFilesTableName+".education_id = "+specialistEducationsTableName+".id").
		LeftJoin(accountFilesTableName+" ON "+accountFilesTableName+".id = "+specialistEducationFilesTableName+".file_id").
		Where(specialistEducationsTableName+".id = ?", id).
		QueryContext(c)

	e, err := s.scanEducations(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return e[0], nil
}

func (s *PostgresStorage) UpdateSpecialistProfileEducation(c context.Context, id interface{}, specialistID interface{}, req *model.UpdateEducation) (*model.Education, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(specialistEducationsTableName).
		Set("institution_id", req.InstitutionID).
		Set("faculty_id", storage.NullInt64(req.FacultyID)).
		Set("department_id", storage.NullInt64(req.DepartmentID)).
		Set("form_id", storage.NullInt64(req.FormID)).
		Set("degree_id", storage.NullInt64(req.DegreeID)).
		Set("graduation", req.Graduation).
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

	if err := s.UpdateSpecialistProfileEducationFiles(id, req.Files); err != nil {
		return nil, err
	}

	return s.GetSpecialistProfileEducationByID(c, id)
}

func (s *PostgresStorage) UpdateSpecialistProfileEducationFiles(educationID interface{}, req []*model.File) error {
	psql := s.SetFormat().RunWith(s.DB)

	_, err := psql.Delete(specialistEducationFilesTableName).Where("education_id = ?", educationID).Exec()
	if err != nil {
		return postgres.ConvertError(err)
	}

	if len(req) > 0 {
		aq := psql.Insert(specialistEducationFilesTableName).Columns("education_id", "file_id")

		for _, v := range req {
			aq = aq.Values(educationID, v.ID)
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

func (s *PostgresStorage) UpdateSpecialistProfileEducationFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdateEducationFields) (*model.Education, error) {
	psql := s.SetFormat().RunWith(s.DB)
	//req["updated_at"] = time.Now()

	res, err := psql.Update(specialistEducationsTableName).
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

	return s.GetSpecialistProfileEducationByID(c, id)
}

func (s *PostgresStorage) DeleteSpecialistProfileEducation(c context.Context, id interface{}, specialistID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(specialistEducationsTableName).
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

func (s *PostgresStorage) educationResponseColumns() []string {
	fields := make([]string, 0)

	fields = append(s.educationResponseColumnsMain(specialistEducationsTableName), s.fileResponseColumns(accountFilesTableName)...)

	return fields
}

func (s *PostgresStorage) educationResponseColumnsMain(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "profile_id",
		pre + "institution_id",
		pre + "faculty_id",
		pre + "department_id",
		pre + "form_id",
		pre + "degree_id",
		pre + "graduation",
		pre + "verified",
	}
}

func (s *PostgresStorage) scanEducations(rows *sql.Rows) ([]*model.Education, error) {
	var es []*model.Education
	educations := make(map[int64]*model.Education)

	for rows.Next() {
		var edu model.Education
		var ef model.FileJoin

		if err := rows.Scan(
			&edu.ID,
			&edu.ProfileID,
			&edu.InstitutionID,
			&edu.FacultyID,
			&edu.DepartmentID,
			&edu.FormID,
			&edu.DegreeID,
			&edu.Graduation,
			&edu.Verified,

			&ef.ID,
			&ef.CreatedAt,
			&ef.UpdatedAt,
			&ef.AccountID,
			&ef.Name,
			&ef.Description,
		); err != nil {
			return nil, err
		}

		if _, ok := educations[edu.ID]; !ok {
			educations[edu.ID] = &edu
		}

		if ef.ID != nil {
			eduFile := ef.ConvertToFile()
			educations[edu.ID].Files = append(educations[edu.ID].Files, &eduFile)
		}
	}

	for _, v := range educations {
		es = append(es, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return es, nil
}
