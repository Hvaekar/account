package account

import (
	"context"
	"database/sql"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type Storage interface {
	Register(c *gin.Context, req *model.RegisterRequest) (*model.Account, error)
	Login(c *gin.Context, req *model.LoginRequest) (*model.Account, error)
	GetAccountByID(c context.Context, id interface{}) (*model.Account, error)
	GetAccountByLogin(c context.Context, login interface{}) (*model.Account, error)
	GetAccounts(c context.Context, req *model.ListAccountsRequest) ([]*model.Account, error)
	DeleteAccount(c context.Context, id interface{}) error
	UpdateAccountMain(c context.Context, id interface{}, req *model.UpdateAccount) (*model.Account, error)
	UpdateAccountFields(c context.Context, id interface{}, req model.UpdateAccountFields) (*model.Account, error)

	AddFile(c context.Context, accountID interface{}, fileName *string) (*model.File, error)
	GetFiles(c context.Context, accountID interface{}) ([]*model.File, error)
	GetPatientDisabilityFiles(c context.Context, patientID interface{}) ([]*model.File, error)
	GetFileByID(c context.Context, id interface{}) (*model.File, error)
	GetFileByName(c context.Context, name interface{}) (*model.File, error)
	UpdateFile(c context.Context, id interface{}, accountID interface{}, req *model.UpdateFile) (*model.File, error)
	UpdateFileFields(c context.Context, id interface{}, accountID interface{}, req model.UpdateFileFields) (*model.File, error)
	DeleteFile(c context.Context, id interface{}, accountID interface{}) (*string, error)

	AddEmail(c context.Context, accountID interface{}, req *model.AddEmail) (*model.Email, error)
	GetEmails(c context.Context, accountID interface{}) ([]*model.Email, error)
	GetEmailByID(c context.Context, id interface{}) (*model.Email, error)
	UpdateEmail(c context.Context, id interface{}, accountID interface{}, req *model.UpdateEmail) (*model.Email, error)
	UpdateEmailFields(c context.Context, id interface{}, accountID interface{}, req model.UpdateEmailFields) (*model.Email, error)
	DeleteEmail(c context.Context, id interface{}, accountID interface{}) error

	AddPhone(c context.Context, accountID interface{}, req *model.AddPhone) (*model.Phone, error)
	GetPhones(c context.Context, accountID interface{}) ([]*model.Phone, error)
	GetPhoneByID(c context.Context, id interface{}) (*model.Phone, error)
	UpdatePhone(c context.Context, id interface{}, accountID interface{}, req *model.UpdatePhone) (*model.Phone, error)
	UpdatePhoneFields(c context.Context, id interface{}, accountID interface{}, req model.UpdatePhoneFields) (*model.Phone, error)
	DeletePhone(c context.Context, id interface{}, accountID interface{}) error

	AddAddress(c context.Context, accountID interface{}, req *model.AddAddress) (*model.Address, error)
	GetAddresses(c context.Context, accountID interface{}) ([]*model.Address, error)
	GetAddressByID(c context.Context, id interface{}) (*model.Address, error)
	UpdateAddress(c context.Context, id interface{}, accountID interface{}, req *model.UpdateAddress) (*model.Address, error)
	UpdateAddressFields(c context.Context, id interface{}, accountID interface{}, req model.UpdateAddressFields) (*model.Address, error)
	DeleteAddress(c context.Context, id interface{}, accountID interface{}) error

	AddLanguage(c context.Context, accountID interface{}, req *model.AddLanguage) (*model.Language, error)
	GetLanguages(c context.Context, accountID interface{}) ([]*model.Language, error)
	GetLanguageByCode(c context.Context, languageCode interface{}, accountID interface{}) (*model.Language, error)
	UpdateLanguage(c context.Context, languageCode interface{}, accountID interface{}, req *model.UpdateLanguage) (*model.Language, error)
	UpdateLanguageFields(c context.Context, languageCode interface{}, accountID interface{}, req model.UpdateLanguageFields) (*model.Language, error)
	DeleteLanguage(c context.Context, languageCode interface{}, accountID interface{}) error

	GetPatientProfileID(c context.Context, accountID interface{}) (*int64, error)
	GetSpecialistProfileID(c context.Context, accountID interface{}) (*int64, error)
	GetPatientProfiles(c context.Context, accountID interface{}) ([]*model.AccountPatient, error)
	VerifyPatientProfile(c context.Context, accountID interface{}, profileID interface{}) (*model.Patient, error)
	DeletePatientProfile(c context.Context, accountID interface{}, profileID interface{}) error

	AddSpecialistProfile(c context.Context, accountID interface{}, req *model.AddSpecialistProfile) (*model.Specialist, error)

	GetPatientByID(c context.Context, id interface{}) (*model.Patient, error)
	GetPatients(c context.Context, req *model.ListPatientsRequest) ([]*model.Patient, error)
	UpdatePatientProfile(c context.Context, id interface{}, req model.UpdatePatientProfile) (*model.Patient, error)
	UpdatePatientProfileFields(c context.Context, id interface{}, req model.UpdatePatientProfileFields) (*model.Patient, error)

	AddPatientProfileAdmin(c context.Context, patientID interface{}, req *model.AddAdmin) (*model.PatientAdmin, error)
	GetPatientProfileAdmins(c context.Context, patientID interface{}) ([]*model.PatientAdmin, error)
	GetPatientProfileAdminByID(c context.Context, patientID interface{}, adminID interface{}) (*model.PatientAdmin, error)
	UpdatePatientProfileAdmin(c context.Context, patientID interface{}, adminID interface{}, req *model.UpdateAdmin) (*model.PatientAdmin, error)
	UpdatePatientProfileAdminFields(c context.Context, patientID interface{}, adminID interface{}, req model.UpdateAdminFields) (*model.PatientAdmin, error)
	DeletePatientProfileAdmin(c context.Context, patientID interface{}, accountID interface{}) error

	AddMetalComponent(c context.Context, patientID interface{}, req *model.AddMetalComponent) (*model.MetalComponent, error)
	GetMetalComponents(c context.Context, patientID interface{}) ([]*model.MetalComponent, error)
	GetMetalComponentByID(c context.Context, id interface{}) (*model.MetalComponent, error)
	UpdateMetalComponent(c context.Context, id interface{}, patientID interface{}, req *model.UpdateMetalComponent) (*model.MetalComponent, error)
	UpdateMetalComponentFields(c context.Context, id interface{}, patientID interface{}, req model.UpdateMetalComponentFields) (*model.MetalComponent, error)
	DeleteMetalComponent(c context.Context, id interface{}, patientID interface{}) error

	GetSpecialistByID(c context.Context, id interface{}) (*model.Specialist, error)
	GetSpecialists(c context.Context, req *model.ListSpecialistsRequest) ([]*model.Specialist, error)
	UpdateSpecialistProfileMain(c context.Context, specialistID interface{}, req *model.UpdateSpecialistProfile) (*model.Specialist, error)
	UpdateSpecialistProfileFields(c context.Context, id interface{}, req model.UpdateSpecialistProfileFields) (*model.Specialist, error)

	AddSpecialistProfileSpecialization(c context.Context, specialistID interface{}, req *model.AddSpecialization) (*model.Specialization, error)
	GetSpecialistProfileSpecializations(c context.Context, specialistID interface{}) ([]*model.Specialization, error)
	GetSpecialistProfileSpecialization(c context.Context, id interface{}, specialistID interface{}) (*model.Specialization, error)
	UpdateSpecialistProfileSpecialization(c context.Context, id interface{}, specialistID interface{}, req *model.UpdateSpecialization) (*model.Specialization, error)
	UpdateSpecialistProfileSpecializationFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdateSpecializationFields) (*model.Specialization, error)
	DeleteSpecialistProfileSpecialization(c context.Context, id interface{}, specialistID interface{}) error

	AddSpecialistProfileEducation(c context.Context, specialistID interface{}, req *model.AddEducation) (*model.Education, error)
	GetSpecialistProfileEducations(c context.Context, specialistID interface{}) ([]*model.Education, error)
	GetSpecialistProfileEducationByID(c context.Context, id interface{}) (*model.Education, error)
	UpdateSpecialistProfileEducation(c context.Context, id interface{}, specialistID interface{}, req *model.UpdateEducation) (*model.Education, error)
	UpdateSpecialistProfileEducationFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdateEducationFields) (*model.Education, error)
	DeleteSpecialistProfileEducation(c context.Context, id interface{}, specialistID interface{}) error

	AddSpecialistProfileExperience(c context.Context, specialistID interface{}, req *model.AddExperience) (*model.Experience, error)
	GetSpecialistProfileExperiences(c context.Context, specialistID interface{}) ([]*model.Experience, error)
	GetSpecialistProfileExperienceByID(c context.Context, id interface{}) (*model.Experience, error)
	UpdateSpecialistProfileExperience(c context.Context, id interface{}, specialistID interface{}, req *model.UpdateExperience) (*model.Experience, error)
	UpdateSpecialistProfileExperienceFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdateExperienceFields) (*model.Experience, error)
	DeleteSpecialistProfileExperience(c context.Context, id interface{}, specialistID interface{}) error

	AddSpecialistProfileAssociation(c context.Context, specialistID interface{}, req *model.AddAssociation) (*model.Association, error)
	GetSpecialistProfileAssociations(c context.Context, specialistID interface{}) ([]*model.Association, error)
	GetSpecialistProfileAssociationByID(c context.Context, id interface{}) (*model.Association, error)
	UpdateSpecialistProfileAssociation(c context.Context, id interface{}, specialistID interface{}, req *model.UpdateAssociation) (*model.Association, error)
	UpdateSpecialistProfileAssociationFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdateAssociationFields) (*model.Association, error)
	DeleteSpecialistProfileAssociation(c context.Context, id interface{}, specialistID interface{}) error

	AddSpecialistProfilePatent(c context.Context, specialistID interface{}, req *model.AddPatent) (*model.Patent, error)
	GetSpecialistProfilePatents(c context.Context, specialistID interface{}) ([]*model.Patent, error)
	GetSpecialistProfilePatentByID(c context.Context, id interface{}) (*model.Patent, error)
	UpdateSpecialistProfilePatent(c context.Context, id interface{}, specialistID interface{}, req *model.UpdatePatent) (*model.Patent, error)
	UpdateSpecialistProfilePatentFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdatePatentFields) (*model.Patent, error)
	DeleteSpecialistProfilePatent(c context.Context, id interface{}, specialistID interface{}) error

	AddSpecialistProfilePublicationLink(c context.Context, specialistID interface{}, req *model.AddPublicationLink) (*model.PublicationLink, error)
	GetSpecialistProfilePublicationLinks(c context.Context, specialistID interface{}) ([]*model.PublicationLink, error)
	GetSpecialistProfilePublicationLinkByID(c context.Context, id interface{}) (*model.PublicationLink, error)
	UpdateSpecialistProfilePublicationLink(c context.Context, id interface{}, specialistID interface{}, req *model.UpdatePublicationLink) (*model.PublicationLink, error)
	UpdateSpecialistProfilePublicationLinkFields(c context.Context, id interface{}, specialistID interface{}, req model.UpdatePublicationLinkFields) (*model.PublicationLink, error)
	DeleteSpecialistProfilePublicationLink(c context.Context, id interface{}, specialistID interface{}) error
}

type PostgresStorage struct {
	*postgres.Postgres
}

func NewPostgresStorage(postgres *postgres.Postgres) *PostgresStorage {
	return &PostgresStorage{Postgres: postgres}
}

func (s *PostgresStorage) SetFormat() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

func (s *PostgresStorage) GetAccountByID(c context.Context, id interface{}) (*model.Account, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.accountResponseColumns()...).
		From(accountsTableName).
		LeftJoin(accountEmailsTableName + " ON " + accountEmailsTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(accountPhonesTableName + " ON " + accountPhonesTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(accountAddressesTableName + " ON " + accountAddressesTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(accountLanguagesTableName + " ON " + accountLanguagesTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(patientProfilesTableName + " ON " + patientProfilesTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(specialistProfilesTableName + " ON " + specialistProfilesTableName + ".account_id = " + accountsTableName + ".id").
		Where(squirrel.Eq{accountsTableName + ".id": id, accountsTableName + ".deleted_at": nil}).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	a, err := s.scanAccount(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	if a.ID == 0 {
		return nil, storage.ErrNotFound
	}

	a.Profiles.Patients, err = s.GetPatientProfiles(c, a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *PostgresStorage) GetAccountByLogin(c context.Context, login interface{}) (*model.Account, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.accountResponseColumns()...).
		From(accountsTableName).
		LeftJoin(accountEmailsTableName + " ON " + accountEmailsTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(accountPhonesTableName + " ON " + accountPhonesTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(accountAddressesTableName + " ON " + accountAddressesTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(accountLanguagesTableName + " ON " + accountLanguagesTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(patientProfilesTableName + " ON " + patientProfilesTableName + ".account_id = " + accountsTableName + ".id").
		LeftJoin(specialistProfilesTableName + " ON " + specialistProfilesTableName + ".account_id = " + accountsTableName + ".id").
		Where(squirrel.Eq{accountsTableName + ".login": login, accountsTableName + ".deleted_at": nil}).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	a, err := s.scanAccount(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	if a.ID == 0 {
		return nil, storage.ErrNotFound
	}

	a.Profiles.Patients, err = s.GetPatientProfiles(c, a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (s *PostgresStorage) GetAccounts(c context.Context, req *model.ListAccountsRequest) ([]*model.Account, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.accountResponseColumnsMain()...).
		From(accountsTableName).
		OrderBy(req.OrderBy).
		Limit(req.Limit).
		Offset(req.Offset()).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	var accs []*model.Account
	for rows.Next() {
		acc, err := s.scanAccountMain(rows)
		if err != nil {
			return nil, postgres.ConvertError(err)
		}
		accs = append(accs, acc)
	}

	if err := rows.Err(); err != nil {
		return nil, postgres.ConvertError(err)
	}

	return accs, nil
}

func (s *PostgresStorage) DeleteAccount(c context.Context, id interface{}) error {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Delete(accountsTableName).
		Where("id = ?", id).
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

func (s *PostgresStorage) UpdateAccountMain(c context.Context, id interface{}, req *model.UpdateAccount) (*model.Account, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(accountsTableName).
		Set("updated_at", time.Now()).
		Set("login", req.Login).
		Set("first_name", storage.NullString(req.FirstName)).
		Set("father_name", storage.NullString(req.FatherName)).
		Set("last_name", storage.NullString(req.LastName)).
		Set("sex", storage.NullString(req.Sex)).
		Set("birthday", storage.NullDatePGX(req.Birthday)).
		Set("language", storage.NullString(req.Language)).
		Set("country", storage.NullString(req.Country)).
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

	return s.GetAccountByID(c, id)
}

func (s *PostgresStorage) UpdateAccountFields(c context.Context, id interface{}, req model.UpdateAccountFields) (*model.Account, error) {
	psql := s.SetFormat().RunWith(s.DB)
	req["updated_at"] = time.Now()

	res, err := psql.Update(accountsTableName).
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

	if _, ok := req["deleted_at"]; !ok {
		return s.GetAccountByID(c, id)
	}

	return nil, nil
}

func (s *PostgresStorage) GetAccountFields(c context.Context, id interface{}, fields ...string) (map[string]interface{}, error) {
	psql := s.SetFormat().RunWith(s.DB)

	row := psql.Select(fields...).
		From(accountsTableName).
		Where("id = ?", id).
		QueryRowContext(c)

	cols := make([]interface{}, len(fields))
	f := make([]interface{}, len(fields))
	for i := range fields {
		f[i] = &cols[i]
	}

	if err := row.Scan(f...); err != nil {
		return nil, postgres.ConvertError(err)
	}

	m := make(map[string]interface{})
	for i := range f {
		value := f[i].(*interface{})
		m[fields[i]] = *value
	}

	return m, nil
}

func (s *PostgresStorage) accountResponseColumns() []string {
	fields := make([]string, 0)

	fields = append(s.accountResponseColumnsMain(accountsTableName), s.emailResponseColumns(accountEmailsTableName)...)
	fields = append(fields, s.phoneResponseColumns(accountPhonesTableName)...)
	fields = append(fields, s.addressResponseColumns(accountAddressesTableName)...)
	fields = append(fields, s.languageResponseColumns(accountLanguagesTableName)...)
	fields = append(fields, s.patientProfileResponseColumns(patientProfilesTableName)...)
	fields = append(fields, s.specialistProfileResponseColumns(specialistProfilesTableName)...)

	return fields
}

func (s *PostgresStorage) accountResponseColumnsMain(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "created_at",
		pre + "updated_at",
		pre + "deleted_at",
		pre + "login",
		pre + "password",
		pre + "first_name",
		pre + "father_name",
		pre + "last_name",
		pre + "sex",
		pre + "photo",
		pre + "birthday",
		pre + "language",
		pre + "country",
	}
}

func (s *PostgresStorage) patientProfileResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
	}
}

func (s *PostgresStorage) specialistProfileResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
	}
}

func (s *PostgresStorage) scanAccount(rows *sql.Rows) (*model.Account, error) {
	var a model.Account
	emails := make(map[string]bool, 0)
	phones := make(map[string]bool, 0)
	addresses := make(map[string]bool, 0)
	languages := make(map[string]bool, 0)

	for rows.Next() {
		var e model.EmailJoin
		var p model.PhoneJoin
		var add model.AddressJoin
		var l model.LanguageJoin
		var specialistID *int64

		if err := rows.Scan(
			&a.ID,
			&a.CreatedAt,
			&a.UpdatedAt,
			&a.DeletedAt,
			&a.Login,
			&a.Password,
			&a.FirstName,
			&a.FatherName,
			&a.LastName,
			&a.Sex,
			&a.Photo,
			&a.Birthday,
			&a.Language,
			&a.Country,

			&e.ID,
			&e.AccountID,
			&e.Type,
			&e.Email,
			&e.Verified,
			&e.Open,

			&p.ID,
			&p.AccountID,
			&p.Type,
			&p.Code,
			&p.Phone,
			&p.Verified,
			&p.Open,

			&add.ID,
			&add.AccountID,
			&add.Type,
			&add.CityID,
			&add.Address,
			&add.Open,

			&l.AccountID,
			&l.Language,
			&l.Level,

			&a.Profiles.PatientProfileID,
			&specialistID,
		); err != nil {
			return nil, err
		}

		if e.ID != nil {
			if _, ok := emails[*e.Email]; !ok {
				emails[*e.Email] = true
				val := e.ConvertToEmail()
				a.Emails = append(a.Emails, &val)
			}
		}

		if p.ID != nil {
			if _, ok := phones[*p.Code+*p.Phone]; !ok {
				phones[*p.Code+*p.Phone] = true
				val := p.ConvertToPhone()
				a.Phones = append(a.Phones, &val)
			}
		}

		if add.ID != nil {
			if _, ok := addresses[strconv.Itoa(int(*add.CityID))+*add.Address]; !ok {
				addresses[strconv.Itoa(int(*add.CityID))+*add.Address] = true
				val := add.ConvertToAddress()
				a.Addresses = append(a.Addresses, &val)
			}
		}

		if l.Language != nil {
			if _, ok := languages[*l.Language]; !ok {
				languages[*l.Language] = true
				val := l.ConvertToLanguage()
				a.Languages = append(a.Languages, &val)
			}
		}

		if specialistID != nil && a.Profiles.SpecialistProfileID == 0 {
			a.Profiles.SpecialistProfileID = *specialistID
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &a, nil
}

func (s *PostgresStorage) scanAccountMain(row squirrel.RowScanner) (*model.Account, error) {
	var a model.Account

	if err := row.Scan(
		&a.ID,
		&a.CreatedAt,
		&a.UpdatedAt,
		&a.DeletedAt,
		&a.Login,
		&a.Password,
		&a.FirstName,
		&a.FatherName,
		&a.LastName,
		&a.Sex,
		&a.Photo,
		&a.Birthday,
		&a.Language,
		&a.Country,
	); err != nil {
		return nil, err
	}

	return &a, nil
}
