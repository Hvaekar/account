package account

import (
	"context"
	"database/sql"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
	"time"
)

func (s *PostgresStorage) GetPatientByID(c context.Context, id interface{}) (*model.Patient, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.patientResponseColumns()...).
		From(patientProfilesTableName).
		LeftJoin(accountsTableName + " ON " + accountsTableName + ".id = " + patientProfilesTableName + ".account_id").
		LeftJoin(accountPhonesTableName + " ON " + accountPhonesTableName + ".id = " + patientProfilesTableName + ".phone_id AND " + patientProfilesTableName + ".phone_id IS NOT NULL").
		LeftJoin(accountEmailsTableName + " ON " + accountEmailsTableName + ".id = " + patientProfilesTableName + ".email_id AND " + patientProfilesTableName + ".email_id IS NOT NULL").
		LeftJoin(patientMetalComponentsTableName + " ON " + patientMetalComponentsTableName + ".profile_id = " + patientProfilesTableName + ".id").
		LeftJoin(patientDisabilityFilesTableName + " ON " + patientDisabilityFilesTableName + ".profile_id = " + patientProfilesTableName + ".id").
		LeftJoin(accountFilesTableName + " ON " + accountFilesTableName + ".id = " + patientDisabilityFilesTableName + ".file_id").
		Where(squirrel.Eq{patientProfilesTableName + ".id": id, accountsTableName + ".deleted_at": nil}).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	p, err := s.scanPatient(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	if p.ID == 0 {
		return nil, storage.ErrNotFound
	}

	p.Admins, err = s.GetPatientProfileAdmins(c, id)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *PostgresStorage) GetPatients(c context.Context, req *model.ListPatientsRequest) ([]*model.Patient, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Select(s.patientsResponseColumns()...).
		From(patientProfilesTableName).
		LeftJoin(accountsTableName + " ON " + accountsTableName + ".id = " + patientProfilesTableName + ".account_id").
		LeftJoin(accountPhonesTableName + " ON " + accountPhonesTableName + ".id = " + patientProfilesTableName + ".phone_id AND " + patientProfilesTableName + ".phone_id IS NOT NULL").
		LeftJoin(accountEmailsTableName + " ON " + accountEmailsTableName + ".id = " + patientProfilesTableName + ".email_id AND " + patientProfilesTableName + ".email_id IS NOT NULL")

	if len(req.IDList) > 0 {
		q = q.Where(squirrel.Eq{patientProfilesTableName + ".id": req.IDList})
	}

	rows, err := q.OrderBy(req.OrderBy).
		Limit(req.Limit).
		Offset(req.Offset()).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	ps, err := s.scanPatients(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return ps, nil
}

func (s *PostgresStorage) UpdatePatientProfile(c context.Context, id interface{}, req model.UpdatePatientProfile) (*model.Patient, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(patientProfilesTableName).
		Set("updated_at", time.Now()).
		Set("phone_id", storage.NullInt64(req.PhoneID)).
		Set("email_id", storage.NullInt64(req.EmailID)).
		Set("height", storage.NullFloat64(req.Height)).
		Set("weight", storage.NullFloat64(req.Weight)).
		Set("body_type", storage.NullString(req.BodyType)).
		Set("blood_type", storage.NullString(req.BloodType)).
		Set("rh", storage.NullBool(req.Rh)).
		Set("left_eye", storage.NullFloat64(req.LeftEye)).
		Set("right_eye", storage.NullFloat64(req.RightEye)).
		Set("disability_group", storage.NullString(req.DisabilityGroup)).
		Set("disability_reason", storage.NullString(req.DisabilityReason)).
		Set("disability_document_num", storage.NullString(req.DisabilityDocumentNum)).
		Set("activity", storage.NullString(req.Activity)).
		Set("nutrition", storage.NullString(req.Nutrition)).
		Set("work", storage.NullString(req.Work)).
		Where("id = ?", id).
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

	if err := s.UpdatePatientProfileDisabilityFiles(id, req.DisabilityFiles); err != nil {
		return nil, err
	}

	return s.GetPatientByID(c, id)
}

func (s *PostgresStorage) UpdatePatientProfileDisabilityFiles(patientID interface{}, req []*model.File) error {
	psql := s.SetFormat().RunWith(s.DB)

	_, err := psql.Delete(patientDisabilityFilesTableName).Where("profile_id = ?", patientID).Exec()
	if err != nil {
		return postgres.ConvertError(err)
	}

	if len(req) > 0 {
		aq := psql.Insert(patientDisabilityFilesTableName).Columns("profile_id", "file_id")

		for _, v := range req {
			aq = aq.Values(patientID, v.ID)
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

func (s *PostgresStorage) UpdatePatientProfileFields(c context.Context, id interface{}, req model.UpdatePatientProfileFields) (*model.Patient, error) {
	psql := s.SetFormat().RunWith(s.DB)
	req["updated_at"] = time.Now()

	res, err := psql.Update(patientProfilesTableName).
		SetMap(req).
		Where("id = ?", id).
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

	return s.GetPatientByID(c, id)
}

func (s *PostgresStorage) patientResponseColumns() []string {
	fields := make([]string, 0)

	fields = append(s.patientResponseColumnsMain(patientProfilesTableName), s.patientAccountResponseColumns(accountsTableName)...)
	fields = append(fields, s.patientPhoneResponseColumns(accountPhonesTableName)...)
	fields = append(fields, s.patientEmailResponseColumns(accountEmailsTableName)...)
	fields = append(fields, s.patientMetalComponentsResponseColumns(patientMetalComponentsTableName)...)
	fields = append(fields, s.fileResponseColumns(accountFilesTableName)...)

	return fields
}

func (s *PostgresStorage) patientsResponseColumns() []string {
	fields := make([]string, 0)

	fields = append(s.patientsResponseColumnsMain(patientProfilesTableName), s.patientAccountResponseColumns(accountsTableName)...)
	fields = append(fields, s.patientPhoneResponseColumns(accountPhonesTableName)...)
	fields = append(fields, s.patientEmailResponseColumns(accountEmailsTableName)...)

	return fields
}

func (s *PostgresStorage) patientResponseColumnsMain(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "account_id",
		//pre + "phone_id",
		//pre + "email_id",
		pre + "height",
		pre + "weight",
		pre + "body_type",
		pre + "blood_type",
		pre + "rh",
		pre + "left_eye",
		pre + "right_eye",
		pre + "disability_group",
		pre + "disability_reason",
		pre + "disability_document_num",
		pre + "activity",
		pre + "nutrition",
		pre + "work",
	}
}

func (s *PostgresStorage) patientsResponseColumnsMain(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
	}
}

func (s *PostgresStorage) patientAccountResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "first_name",
		pre + "father_name",
		pre + "last_name",
		pre + "sex",
		pre + "photo",
		pre + "birthday",
	}
}

func (s *PostgresStorage) patientPhoneResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "code",
		pre + "phone",
	}
}

func (s *PostgresStorage) patientEmailResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "email",
	}
}

func (s *PostgresStorage) scanPatient(rows *sql.Rows) (*model.Patient, error) {
	var p model.Patient
	mcs := make(map[int64]*model.MetalComponent)
	files := make(map[int64]*model.File)

	for rows.Next() {
		var phone model.PhoneJoin
		var email model.EmailJoin
		var mc model.MetalComponentJoin
		var df model.FileJoin

		if err := rows.Scan(
			&p.ID,
			&p.AccountID,
			//&p.Phone,
			//&p.Email,
			&p.Height,
			&p.Weight,
			&p.BodyType,
			&p.BloodType,
			&p.Rh,
			&p.LeftEye,
			&p.RightEye,
			&p.Disability.Group,
			&p.Disability.Reason,
			&p.Disability.DocumentNum,
			&p.Activity,
			&p.Nutrition,
			&p.Work,

			&p.FirstName,
			&p.FatherName,
			&p.LastName,
			&p.Sex,
			&p.Photo,
			&p.Birthday,

			&phone.Code,
			&phone.Phone,

			&email.Email,

			&mc.ID,
			&mc.PatientID,
			&mc.Metal,
			&mc.OrganID,
			&mc.Description,

			&df.ID,
			&df.CreatedAt,
			&df.UpdatedAt,
			&df.AccountID,
			&df.Name,
			&df.Description,
		); err != nil {
			return nil, err
		}

		if phone.Phone != nil && p.Phone == nil {
			v := *phone.Code + *phone.Phone
			p.Phone = &v
		}

		if email.Email != nil && p.Email == nil {
			v := *email.Email
			p.Email = &v
		}

		if mc.ID != nil {
			if _, ok := mcs[*mc.ID]; !ok {
				metalComponent := mc.ConvertToMetalComponent()
				mcs[*mc.ID] = &metalComponent
				p.MetalComponents = append(p.MetalComponents, &metalComponent)
			}
		}

		if df.ID != nil {
			if _, ok := files[*df.ID]; !ok {
				disabilityFile := df.ConvertToFile()
				files[*df.ID] = &disabilityFile
				p.Disability.Files = append(p.Disability.Files, &disabilityFile)
			}
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *PostgresStorage) scanPatients(rows *sql.Rows) ([]*model.Patient, error) {
	ps := make([]*model.Patient, 0)

	for rows.Next() {
		var p model.Patient
		var phone model.PhoneJoin
		var email model.EmailJoin

		if err := rows.Scan(
			&p.ID,

			&p.FirstName,
			&p.FatherName,
			&p.LastName,
			&p.Sex,
			&p.Photo,
			&p.Birthday,

			&phone.Code,
			&phone.Phone,

			&email.Email,
		); err != nil {
			return nil, err
		}

		if phone.Phone != nil {
			v := *phone.Code + *phone.Phone
			p.Phone = &v
		}

		if email.Email != nil {
			v := *email.Email
			p.Email = &v
		}

		ps = append(ps, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ps, nil
}
