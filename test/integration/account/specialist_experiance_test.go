package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ExperienceTestSuite struct {
	TestSuite
}

func TestExperienceSuite(t *testing.T) {
	suite.Run(t, new(ExperienceTestSuite))
}

func (s *ExperienceTestSuite) TestAddExperience() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var companyID int64 = 100
	start, err := time.ParseInLocation("2006-01-02", "2020-08-23", time.UTC)
	s.Require().NoError(err)
	finish, err := time.ParseInLocation("2006-01-02", "2023-09-01", time.UTC)
	s.Require().NoError(err)
	fin := pgtype.Date{Time: finish, Valid: true}
	req := model.AddExperience{
		CompanyID:       &companyID,
		Company:         "Some company name",
		Start:           pgtype.Date{Time: start, Valid: true},
		Finish:          &fin,
		Specializations: []int64{55, 66, 44},
	}

	experience, err := s.client.AddExperience(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.NotEmpty(experience.ID)
	s.Equal(*req.CompanyID, *experience.CompanyID)
	s.Equal(req.Company, experience.Company)
	s.Equal(req.Start, experience.Start)
	s.Equal(req.Finish.Time, experience.Finish.Time)
	s.Len(experience.Specializations, 3)
	for _, v := range req.Specializations {
		s.Contains(experience.Specializations, v)
	}

	specializations, err := s.client.GetSpecializations(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(specializations.Specializations, 6)
}

func (s *ExperienceTestSuite) TestGetExperiences() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetExperiences(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Experiences, 2)
}

func (s *ExperienceTestSuite) TestUpdateExperience() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var companyID int64 = 100
	start, err := time.ParseInLocation("2006-01-02", "2020-08-23", time.UTC)
	s.Require().NoError(err)
	finish, err := time.ParseInLocation("2006-01-02", "2023-09-01", time.UTC)
	s.Require().NoError(err)
	fin := pgtype.Date{Time: finish, Valid: true}
	req := model.UpdateExperience{
		CompanyID:       &companyID,
		Company:         "Some another company name",
		Start:           pgtype.Date{Time: start, Valid: true},
		Finish:          &fin,
		Specializations: []int64{55, 66, 44},
	}

	experience, err := s.client.UpdateExperience(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(*req.CompanyID, *experience.CompanyID)
	s.Equal(req.Company, experience.Company)
	s.Equal(req.Start, experience.Start)
	s.Equal(req.Finish.Time, experience.Finish.Time)
	s.Len(experience.Specializations, 3)
	for _, v := range req.Specializations {
		s.Contains(experience.Specializations, v)
	}
}

func (s *ExperienceTestSuite) TestUpdateExperienceNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var companyID int64 = 100
	start, err := time.ParseInLocation("2006-01-02", "2020-08-23", time.UTC)
	s.Require().NoError(err)
	finish, err := time.ParseInLocation("2006-01-02", "2023-09-01", time.UTC)
	s.Require().NoError(err)
	fin := pgtype.Date{Time: finish, Valid: true}
	req := model.UpdateExperience{
		CompanyID:       &companyID,
		Company:         "Some another company name",
		Start:           pgtype.Date{Time: start, Valid: true},
		Finish:          &fin,
		Specializations: []int64{55, 66, 44},
	}

	_, err = s.client.UpdateExperience(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *ExperienceTestSuite) TestDeleteExperience() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteExperience(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetExperiences(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Experiences, 1)
	for _, e := range list.Experiences {
		s.NotEqual(int64(1), e.ID)
	}
}
