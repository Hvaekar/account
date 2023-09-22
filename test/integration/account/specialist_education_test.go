package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type EducationTestSuite struct {
	TestSuite
}

func TestEducationSuite(t *testing.T) {
	suite.Run(t, new(EducationTestSuite))
}

func (s *EducationTestSuite) TestAddEducation() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var facultyID int64 = 2
	var departmentID int64 = 3
	var formID int64 = 4
	var degreeID int64 = 5
	graduation, err := time.ParseInLocation("2006-01-02", "2017-06-01", time.UTC)
	s.Require().NoError(err)
	req := model.AddEducation{
		InstitutionID: 1,
		FacultyID:     &facultyID,
		DepartmentID:  &departmentID,
		FormID:        &formID,
		DegreeID:      &degreeID,
		Graduation:    pgtype.Date{Time: graduation, Valid: true},
		Files: []*model.File{
			{
				ID: 1,
			},
			{
				ID: 2,
			},
		},
	}

	education, err := s.client.AddEducation(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.NotEmpty(education.ID)
	s.Equal(req.InstitutionID, education.InstitutionID)
	s.Equal(*req.FacultyID, *education.FacultyID)
	s.Equal(*req.DepartmentID, *education.DepartmentID)
	s.Equal(*req.FormID, *education.FormID)
	s.Equal(*req.DegreeID, *education.DegreeID)
	s.Equal(req.Graduation.Time, education.Graduation.Time)
	s.Len(education.Files, 2)
}

func (s *EducationTestSuite) TestGetEducations() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetEducations(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Educations, 1)
}

func (s *EducationTestSuite) TestUpdateEducation() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var facultyID int64 = 3
	var departmentID int64 = 4
	var formID int64 = 5
	var degreeID int64 = 6
	graduation, err := time.ParseInLocation("2006-01-02", "2017-06-01", time.UTC)
	s.Require().NoError(err)
	req := model.UpdateEducation{
		InstitutionID: 2,
		FacultyID:     &facultyID,
		DepartmentID:  &departmentID,
		FormID:        &formID,
		DegreeID:      &degreeID,
		Graduation:    pgtype.Date{Time: graduation, Valid: true},
		Files: []*model.File{
			{
				ID: 2,
			},
		},
	}

	education, err := s.client.UpdateEducation(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(req.InstitutionID, education.InstitutionID)
	s.Equal(*req.FacultyID, *education.FacultyID)
	s.Equal(*req.DepartmentID, *education.DepartmentID)
	s.Equal(*req.FormID, *education.FormID)
	s.Equal(*req.DegreeID, *education.DegreeID)
	s.Equal(req.Graduation.Time, education.Graduation.Time)
	s.Len(education.Files, 1)
	s.Equal(int64(2), education.Files[0].ID)
}

func (s *EducationTestSuite) TestUpdateEducationNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var facultyID int64 = 3
	var departmentID int64 = 4
	var formID int64 = 5
	var degreeID int64 = 6
	graduation, err := time.ParseInLocation("2006-01-02", "2017-06-01", time.UTC)
	s.Require().NoError(err)
	req := model.UpdateEducation{
		InstitutionID: 2,
		FacultyID:     &facultyID,
		DepartmentID:  &departmentID,
		FormID:        &formID,
		DegreeID:      &degreeID,
		Graduation:    pgtype.Date{Time: graduation, Valid: true},
		Files: []*model.File{
			{
				ID: 2,
			},
		},
	}

	_, err = s.client.UpdateEducation(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *EducationTestSuite) TestDeleteEducation() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteEducation(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetEducations(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Educations, 0)
	//for _, e := range list.Educations {
	//	s.NotEqual(int64(1), e.ID)
	//}
}
