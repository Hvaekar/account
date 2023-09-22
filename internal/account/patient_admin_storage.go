package account

import (
	"context"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) AddPatientProfileAdmin(c context.Context, patientID interface{}, req *model.AddAdmin) (*model.PatientAdmin, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Insert(accountsPatientProfilesTableName).
		Columns(
			"account_id",
			"patient_profile_id",
			"permission_edit",
		).
		Values(
			req.AdminID,
			patientID,
			req.PermissionEdit,
		).
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

	return s.GetPatientProfileAdminByID(c, patientID, req.AdminID)
}

func (s *PostgresStorage) GetPatientProfileAdmins(c context.Context, patientID interface{}) ([]*model.PatientAdmin, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.patientAdminsResponseColumns()...).
		From(accountsPatientProfilesTableName).
		LeftJoin(accountsTableName+" ON "+accountsTableName+".id = "+accountsPatientProfilesTableName+".account_id").
		Where(accountsPatientProfilesTableName+".patient_profile_id = ?", patientID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var admins []*model.PatientAdmin
	for rows.Next() {
		admin, err := s.scanPatientAdmin(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		admins = append(admins, admin)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return admins, nil
}

func (s *PostgresStorage) GetPatientProfileAdminByID(c context.Context, patientID interface{}, adminID interface{}) (*model.PatientAdmin, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(s.patientAdminsResponseColumns()...).
		From(accountsPatientProfilesTableName).
		LeftJoin(accountsTableName+" ON "+accountsTableName+".id = "+accountsPatientProfilesTableName+".account_id").
		Where(accountsPatientProfilesTableName+".account_id = ? AND "+accountsPatientProfilesTableName+".patient_profile_id = ?", adminID, patientID).
		QueryRowContext(c)

	admin, err := s.scanPatientAdmin(row)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return admin, nil
}

func (s *PostgresStorage) UpdatePatientProfileAdmin(c context.Context, patientID interface{}, adminID interface{}, req *model.UpdateAdmin) (*model.PatientAdmin, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(accountsPatientProfilesTableName).
		Set("permission_edit", req.PermissionEdit).
		Where("account_id = ? AND patient_profile_id = ?", adminID, patientID).
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

	return s.GetPatientProfileAdminByID(c, patientID, adminID)
}

func (s *PostgresStorage) UpdatePatientProfileAdminFields(c context.Context, patientID interface{}, adminID interface{}, req model.UpdateAdminFields) (*model.PatientAdmin, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(accountsPatientProfilesTableName).
		SetMap(req).
		Where("account_id = ? AND patient_profile_id = ?", adminID, patientID).
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

	return s.GetPatientProfileAdminByID(c, patientID, adminID)
}

func (s *PostgresStorage) DeletePatientProfileAdmin(c context.Context, patientID interface{}, accountID interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(accountsPatientProfilesTableName).
		Where("account_id = ? AND patient_profile_id = ?", accountID, patientID).
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

func (s *PostgresStorage) patientAdminsResponseColumns() []string {
	fields := make([]string, 0)

	fields = append(s.patientAdminResponseColumns(accountsPatientProfilesTableName), s.patientAdminAccountResponseColumns(accountsTableName)...)

	return fields
}

func (s *PostgresStorage) patientAdminResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "account_id",
		pre + "patient_profile_id",
		pre + "permission_edit",
		pre + "verified",
	}
}

func (s *PostgresStorage) patientAdminAccountResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "first_name",
		pre + "father_name",
		pre + "last_name",
		pre + "photo",
	}
}

func (s *PostgresStorage) scanPatientAdmin(row squirrel.RowScanner) (*model.PatientAdmin, error) {
	var a model.PatientAdmin
	var aRes model.PatientAdmin

	var patientID *int64
	if err := row.Scan(
		&a.ID,
		&patientID,
		&a.PermissionEdit,
		&a.Verified,
		&a.FirstName,
		&a.FatherName,
		&a.LastName,
		&a.Photo,
	); err != nil {
		return nil, err
	}

	if !a.Verified {
		aRes = model.PatientAdmin{
			ID:             a.ID,
			PermissionEdit: a.PermissionEdit,
		}
	} else {
		aRes = a
	}

	return &aRes, nil
}
