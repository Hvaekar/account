package account

import (
	"context"
	"database/sql"
	"errors"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
)

func (s *PostgresStorage) GetPatientProfileID(c context.Context, accountID interface{}) (*int64, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select("id").
		From(patientProfilesTableName).
		Where("account_id = ?", accountID).
		QueryRowContext(c)

	var id int64
	if err := row.Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return &id, nil
}

func (s *PostgresStorage) GetSpecialistProfileID(c context.Context, accountID interface{}) (*int64, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select("id").
		From(specialistProfilesTableName).
		Where("account_id = ?", accountID).
		QueryRowContext(c)

	var id *int64
	if err := row.Scan(&id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, postgres.ConvertError(err)
	}

	return id, nil
}

func (s *PostgresStorage) GetPatientProfiles(c context.Context, accountID interface{}) ([]*model.AccountPatient, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.accountPatientResponseColumns()...).
		From(accountsPatientProfilesTableName).
		LeftJoin(accountsTableName+" ON "+accountsTableName+".id = "+accountsPatientProfilesTableName+".patient_profile_id").
		Where(accountsPatientProfilesTableName+".account_id = ?", accountID).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var profiles []*model.AccountPatient
	for rows.Next() {
		profile, err := s.scanAccountPatientProfile(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		profiles = append(profiles, profile)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return profiles, nil
}

func (s *PostgresStorage) VerifyPatientProfile(c context.Context, accountID interface{}, profileID interface{}) (*model.Patient, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(accountsPatientProfilesTableName).
		//Set("updated_at", time.Now()).
		Set("verified", true).
		Where("account_id = ? AND patient_profile_id = ?", accountID, profileID).
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

	return s.GetPatientByID(c, profileID)
}

func (s *PostgresStorage) DeletePatientProfile(c context.Context, accountID interface{}, patientID interface{}) error {
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

func (s *PostgresStorage) accountPatientResponseColumns() []string {
	fields := make([]string, 0)

	fields = append(s.accountsPatientProfilesResponseColumns(accountsPatientProfilesTableName), s.accountPatientProfileResponseColumns(accountsTableName)...)

	return fields
}

func (s *PostgresStorage) accountsPatientProfilesResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "account_id",
		pre + "patient_profile_id",
		pre + "permission_edit",
		pre + "verified",
	}
}

func (s *PostgresStorage) accountPatientProfileResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "first_name",
		pre + "father_name",
		pre + "last_name",
		pre + "photo",
	}
}

func (s *PostgresStorage) scanAccountPatientProfile(row squirrel.RowScanner) (*model.AccountPatient, error) {
	var p model.AccountPatient

	var accountID *int64
	if err := row.Scan(
		&accountID,
		&p.ID,
		&p.PermissionEdit,
		&p.Verified,
		&p.FirstName,
		&p.FatherName,
		&p.LastName,
		&p.Photo,
	); err != nil {
		return nil, err
	}

	return &p, nil
}
