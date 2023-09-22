package account

import (
	"context"
	"database/sql"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/Masterminds/squirrel"
	"time"
)

func (s *PostgresStorage) AddSpecialistProfile(c context.Context, accountID interface{}, req *model.AddSpecialistProfile) (*model.Specialist, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(specialistProfilesTableName).
		Columns(
			"account_id",
			"phone_id",
			"email_id",
			"about",
			"medical_category",
			"treats_adults",
			"treats_children",
		).
		Values(
			accountID,
			storage.NullInt64(req.PhoneID),
			storage.NullInt64(req.EmailID),
			storage.NullString(req.About),
			storage.NullString(req.MedicalCategory),
			req.TreatsAdults,
			req.TreatsChildren,
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	if id != 0 {
		if err := s.UpdateSpecialistProfileCuresDiseases(id, req.CuresDiseases); err != nil {
			return nil, err
		}
		if err := s.UpdateSpecialistProfileServices(id, req.Services); err != nil {
			return nil, err
		}
	}

	return s.GetSpecialistByID(c, id)
}

func (s *PostgresStorage) GetSpecialistByID(c context.Context, id interface{}) (*model.Specialist, error) {
	psql := s.SetFormat().RunWith(s.DB)

	rows, err := psql.Select(s.specialistResponseColumns()...).
		From(specialistProfilesTableName).
		LeftJoin(accountsTableName + " ON " + accountsTableName + ".id = " + specialistProfilesTableName + ".account_id").
		LeftJoin(accountPhonesTableName + " ON " + accountPhonesTableName + ".id = " + specialistProfilesTableName + ".phone_id AND " + specialistProfilesTableName + ".phone_id IS NOT NULL").
		LeftJoin(accountEmailsTableName + " ON " + accountEmailsTableName + ".id = " + specialistProfilesTableName + ".email_id AND " + specialistProfilesTableName + ".email_id IS NOT NULL").
		LeftJoin(specialistSpecializationsTableName + " ON " + specialistSpecializationsTableName + ".profile_id = " + specialistProfilesTableName + ".id").
		LeftJoin(specialistCuresDiseasesTableName + " ON " + specialistCuresDiseasesTableName + ".profile_id = " + specialistProfilesTableName + ".id").
		LeftJoin(specialistServicesTableName + " ON " + specialistServicesTableName + ".profile_id = " + specialistProfilesTableName + ".id").
		LeftJoin(specialistEducationsTableName + " ON " + specialistEducationsTableName + ".profile_id = " + specialistProfilesTableName + ".id").
		LeftJoin(specialistEducationFilesTableName + " ON " + specialistEducationFilesTableName + ".education_id = " + specialistEducationsTableName + ".id").
		LeftJoin(accountFilesTableName + " ON " + accountFilesTableName + ".id = " + specialistEducationFilesTableName + ".file_id").
		LeftJoin(specialistExperiencesTableName + " ON " + specialistExperiencesTableName + ".profile_id = " + specialistProfilesTableName + ".id").
		LeftJoin(specialistExperienceSpecializationsTableName + " ON " + specialistExperienceSpecializationsTableName + ".experience_id = " + specialistExperiencesTableName + ".id").
		LeftJoin(specialistAssociationsTableName + " ON " + specialistAssociationsTableName + ".profile_id = " + specialistProfilesTableName + ".id").
		LeftJoin(specialistPatentsTableName + " ON " + specialistPatentsTableName + ".profile_id = " + specialistProfilesTableName + ".id").
		LeftJoin(specialistPublicationLinksTableName + " ON " + specialistPublicationLinksTableName + ".profile_id = " + specialistProfilesTableName + ".id").
		Where(squirrel.Eq{specialistProfilesTableName + ".id": id, accountsTableName + ".deleted_at": nil}).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	p, err := s.scanSpecialist(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	if p.ID == 0 {
		return nil, storage.ErrNotFound
	}

	return p, nil
}

func (s *PostgresStorage) GetSpecialists(c context.Context, req *model.ListSpecialistsRequest) ([]*model.Specialist, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Select(s.specialistsResponseColumns()...).
		From(specialistProfilesTableName).
		LeftJoin(accountsTableName + " ON " + accountsTableName + ".id = " + specialistProfilesTableName + ".account_id").
		LeftJoin(accountPhonesTableName + " ON " + accountPhonesTableName + ".id = " + specialistProfilesTableName + ".phone_id AND " + specialistProfilesTableName + ".phone_id IS NOT NULL").
		LeftJoin(accountEmailsTableName + " ON " + accountEmailsTableName + ".id = " + specialistProfilesTableName + ".email_id AND " + specialistProfilesTableName + ".email_id IS NOT NULL").
		LeftJoin(specialistSpecializationsTableName + " ON " + specialistSpecializationsTableName + ".profile_id = " + specialistProfilesTableName + ".id")

	if len(req.IDList) > 0 {
		q = q.Where(squirrel.Eq{specialistProfilesTableName + ".id": req.IDList})
	}

	rows, err := q.OrderBy(req.OrderBy).
		Limit(req.Limit).
		Offset(req.Offset()).
		QueryContext(c)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}
	defer rows.Close()

	ss, err := s.scanSpecialists(rows)
	if err != nil {
		return nil, postgres.ConvertError(err)
	}

	return ss, nil
}

func (s *PostgresStorage) UpdateSpecialistProfileMain(c context.Context, specialistID interface{}, req *model.UpdateSpecialistProfile) (*model.Specialist, error) {
	psql := s.SetFormat().RunWith(s.DB)

	res, err := psql.Update(specialistProfilesTableName).
		Set("updated_at", time.Now()).
		Set("phone_id", storage.NullInt64(req.PhoneID)).
		Set("email_id", storage.NullInt64(req.EmailID)).
		Set("about", storage.NullString(req.About)).
		Set("medical_category", storage.NullString(req.MedicalCategory)).
		Set("treats_adults", req.TreatsAdults).
		Set("treats_children", req.TreatsChildren).
		Where("id = ?", specialistID).
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

	if err := s.UpdateSpecialistProfileCuresDiseases(specialistID, req.CuresDiseases); err != nil {
		return nil, err
	}
	if err := s.UpdateSpecialistProfileServices(specialistID, req.Services); err != nil {
		return nil, err
	}

	return s.GetSpecialistByID(c, specialistID)
}

func (s *PostgresStorage) UpdateSpecialistProfileCuresDiseases(specialistID interface{}, req []int64) error {
	psql := s.SetFormat().RunWith(s.DB)

	_, err := psql.Delete(specialistCuresDiseasesTableName).Where("profile_id = ?", specialistID).Exec()
	if err != nil {
		return postgres.ConvertError(err)
	}

	if len(req) > 0 {
		aq := psql.Insert(specialistCuresDiseasesTableName).Columns("profile_id", "disease_id")

		for _, v := range req {
			aq = aq.Values(specialistID, v)
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

func (s *PostgresStorage) UpdateSpecialistProfileServices(specialistID interface{}, req []int64) error {
	psql := s.SetFormat().RunWith(s.DB)

	_, err := psql.Delete(specialistServicesTableName).Where("profile_id = ?", specialistID).Exec()
	if err != nil {
		return postgres.ConvertError(err)
	}

	if len(req) > 0 {
		aq := psql.Insert(specialistServicesTableName).Columns("profile_id", "service_id")

		for _, v := range req {
			aq = aq.Values(specialistID, v)
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

func (s *PostgresStorage) UpdateSpecialistProfileFields(c context.Context, id interface{}, req model.UpdateSpecialistProfileFields) (*model.Specialist, error) {
	psql := s.SetFormat().RunWith(s.DB)
	req["updated_at"] = time.Now()

	res, err := psql.Update(specialistProfilesTableName).
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

	return s.GetSpecialistByID(c, id)
}

func (s *PostgresStorage) specialistResponseColumns() []string {
	fields := make([]string, 0)

	fields = append(s.specialistResponseColumnsMain(specialistProfilesTableName), s.specialistAccountResponseColumns(accountsTableName)...)
	fields = append(fields, s.specialistPhoneResponseColumns(accountPhonesTableName)...)
	fields = append(fields, s.specialistEmailResponseColumns(accountEmailsTableName)...)
	fields = append(fields, s.specializationResponseColumns(specialistSpecializationsTableName)...)
	fields = append(fields, s.specialistCuresDiseaseResponseColumns(specialistCuresDiseasesTableName)...)
	fields = append(fields, s.specialistServiceResponseColumns(specialistServicesTableName)...)
	fields = append(fields, s.educationResponseColumnsMain(specialistEducationsTableName)...)
	fields = append(fields, s.fileResponseColumns(accountFilesTableName)...)
	fields = append(fields, s.experienceResponseColumnsMain(specialistExperiencesTableName)...)
	fields = append(fields, s.experienceSpecializationResponseColumns(specialistExperienceSpecializationsTableName)...)
	fields = append(fields, s.associationResponseColumns(specialistAssociationsTableName)...)
	fields = append(fields, s.patentResponseColumns(specialistPatentsTableName)...)
	fields = append(fields, s.publicationLinkResponseColumns(specialistPublicationLinksTableName)...)

	return fields
}

func (s *PostgresStorage) specialistsResponseColumns() []string {
	fields := make([]string, 0)

	fields = append(s.specialistResponseColumnsMain(specialistProfilesTableName), s.specialistAccountResponseColumns(accountsTableName)...)
	fields = append(fields, s.specialistPhoneResponseColumns(accountPhonesTableName)...)
	fields = append(fields, s.specialistEmailResponseColumns(accountEmailsTableName)...)
	fields = append(fields, s.specializationResponseColumns(specialistSpecializationsTableName)...)

	return fields
}

func (s *PostgresStorage) specialistResponseColumnsMain(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "id",
		pre + "updated_at",
		pre + "about",
		pre + "medical_category",
		pre + "treats_adults",
		pre + "treats_children",
	}
}

func (s *PostgresStorage) specialistAccountResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "first_name",
		pre + "father_name",
		pre + "last_name",
		pre + "sex",
		pre + "photo",
	}
}

func (s *PostgresStorage) specialistPhoneResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "code",
		pre + "phone",
	}
}

func (s *PostgresStorage) specialistEmailResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "email",
	}
}

func (s *PostgresStorage) specialistCuresDiseaseResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "disease_id",
	}
}

func (s *PostgresStorage) specialistServiceResponseColumns(prefixes ...string) []string {
	pre := joinPrefixes(prefixes, ".")

	return []string{
		pre + "service_id",
	}
}

func (s *PostgresStorage) scanSpecialist(rows *sql.Rows) (*model.Specialist, error) {
	var sp model.Specialist
	specializations := make(map[int64]*model.Specialization)
	curesDiseases := make(map[int64]struct{})
	services := make(map[int64]struct{})
	educations := make(map[int64]*model.Education)
	experiences := make(map[int64]*model.Experience)
	associations := make(map[int64]*model.Association)
	patents := make(map[int64]*model.Patent)
	publicationLinks := make(map[int64]*model.PublicationLink)

	for rows.Next() {
		var phone model.PhoneJoin
		var email model.EmailJoin
		var specialization model.SpecializationJoin
		var cd *int64
		var svc *int64
		var edu model.EducationJoin
		var ef model.FileJoin
		var exp model.ExperienceJoin
		var expSp *int64
		var ass model.AssociationJoin
		var pat model.PatentJoin
		var pl model.PublicationLinkJoin

		if err := rows.Scan(
			&sp.ID,
			&sp.UpdatedAt,
			&sp.About,
			&sp.MedicalCategory,
			&sp.TreatsAdults,
			&sp.TreatsChildren,

			&sp.FirstName,
			&sp.FatherName,
			&sp.LastName,
			&sp.Sex,
			&sp.Photo,

			&phone.Code,
			&phone.Phone,

			&email.Email,

			&specialization.SpecializationID,
			&specialization.Start,

			&cd,

			&svc,

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

			&exp.ID,
			&exp.ProfileID,
			&exp.CompanyID,
			&exp.Company,
			&exp.Start,
			&exp.Finish,

			&expSp,

			&ass.ID,
			&ass.ProfileID,
			&ass.AssociationID,
			&ass.Name,
			&ass.JobTitle,

			&pat.ID,
			&pat.ProfileID,
			&pat.Number,
			&pat.Name,
			&pat.Link,

			&pl.ID,
			&pl.ProfileID,
			&pl.Title,
			&pl.Link,
		); err != nil {
			return nil, err
		}

		if phone.Phone != nil && sp.Phone == nil {
			v := *phone.Code + *phone.Phone
			sp.Phone = &v
		}

		if email.Email != nil && sp.Email == nil {
			v := *email.Email
			sp.Email = &v
		}

		if specialization.SpecializationID != nil {
			if _, ok := specializations[*specialization.SpecializationID]; !ok {
				val := specialization.ConvertToSpecialization()
				specializations[*specialization.SpecializationID] = &val
				sp.Specializations = append(sp.Specializations, &val)
			}
		}

		if cd != nil {
			if _, ok := curesDiseases[*cd]; !ok {
				curesDiseases[*cd] = struct{}{}
				sp.CuresDiseases = append(sp.CuresDiseases, *cd)
			}
		}

		if svc != nil {
			if _, ok := services[*svc]; !ok {
				services[*svc] = struct{}{}
				sp.Services = append(sp.Services, *svc)
			}
		}

		if edu.ID != nil {
			if _, ok := educations[*edu.ID]; !ok {
				val := edu.ConvertToEducation()
				educations[*edu.ID] = &val
			}
			if ef.ID != nil {
				eduFile := ef.ConvertToFile()
				educations[*edu.ID].Files = append(educations[*edu.ID].Files, &eduFile)
			}
		}

		if exp.ID != nil {
			if _, ok := experiences[*exp.ID]; !ok {
				val := exp.ConvertToExperience()
				experiences[*exp.ID] = &val
			}
			if expSp != nil {
				experiences[*exp.ID].Specializations = append(experiences[*exp.ID].Specializations, *expSp)
			}
		}

		if ass.ID != nil {
			if _, ok := associations[*ass.ID]; !ok {
				val := ass.ConvertToAssociation()
				associations[*ass.ID] = &val
				sp.Associations = append(sp.Associations, &val)
			}
		}

		if pat.ID != nil {
			if _, ok := patents[*pat.ID]; !ok {
				val := pat.ConvertToPatent()
				patents[*pat.ID] = &val
				sp.Patents = append(sp.Patents, &val)
			}
		}

		if pl.ID != nil {
			if _, ok := publicationLinks[*pl.ID]; !ok {
				val := pl.ConvertToPublicationLink()
				publicationLinks[*pl.ID] = &val
				sp.PublicationLinks = append(sp.PublicationLinks, &val)
			}
		}
	}

	for _, v := range educations {
		v.Files = model.MatchingUniqueFiles(v.Files)
		sp.Educations = append(sp.Educations, v)
	}

	for _, v := range experiences {
		v.Specializations = utils.MatchingUniqueInt64(v.Specializations)
		sp.Experiences = append(sp.Experiences, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &sp, nil
}

func (s *PostgresStorage) scanSpecialists(rows *sql.Rows) ([]*model.Specialist, error) {
	ss := make([]*model.Specialist, 0)
	ssM := make(map[int64]*model.Specialist)

	for rows.Next() {
		var sp model.Specialist
		var phone model.PhoneJoin
		var email model.EmailJoin
		var specialization model.SpecializationJoin

		if err := rows.Scan(
			&sp.ID,
			&sp.UpdatedAt,
			&sp.About,
			&sp.MedicalCategory,
			&sp.TreatsAdults,
			&sp.TreatsChildren,

			&sp.FirstName,
			&sp.FatherName,
			&sp.LastName,
			&sp.Sex,
			&sp.Photo,

			&phone.Code,
			&phone.Phone,

			&email.Email,

			&specialization.SpecializationID,
			&specialization.Start,
		); err != nil {
			return nil, err
		}

		if phone.Phone != nil {
			v := *phone.Code + *phone.Phone
			sp.Phone = &v
		}

		if email.Email != nil {
			v := *email.Email
			sp.Email = &v
		}

		if specialization.SpecializationID != nil {
			val := specialization.ConvertToSpecialization()

			if _, ok := ssM[sp.ID]; ok {
				ssM[sp.ID].Specializations = append(ssM[sp.ID].Specializations, &val)
			} else {
				sp.Specializations = append(sp.Specializations, &val)
				ssM[sp.ID] = &sp
			}
		}

	}

	for _, v := range ssM {
		ss = append(ss, v)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ss, nil
}
