package account

import (
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ProfileTestSuite struct {
	TestSuite
}

func TestProfileSuite(t *testing.T) {
	suite.Run(t, new(ProfileTestSuite))
}

func (s *ProfileTestSuite) TestGetProfiles() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetProfiles(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Equal(int64(1), list.PatientProfileID)
	s.Equal(int64(1), list.SpecialistProfileID)
	s.Len(list.Patients, 2)
	for _, p := range list.Patients {
		s.NotEqual(int64(1), p.ID)
	}
}

func (s *ProfileTestSuite) TestVerifyPatientProfile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	patient, err := s.client.VerifyPatientProfile(s.ctx, s.token.Access, 2)
	s.Require().NoError(err)

	s.Equal(int64(2), patient.ID)

	list, err := s.client.GetProfiles(s.ctx, s.token.Access)
	s.Require().NoError(err)
	for _, p := range list.Patients {
		if p.ID == 2 {
			s.True(p.Verified)
		}
	}
}

func (s *ProfileTestSuite) TestVerifyPatientProfileBadAccount() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 100,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	_, err = s.client.VerifyPatientProfile(s.ctx, *token, 1)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *ProfileTestSuite) TestSelectPatientProfile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	newToken, cookie, err := s.client.SelectPatientProfile(s.ctx, s.token.Access, 2)
	s.Require().NoError(err)

	s.NotNil(cookie)
	s.NotEmpty(newToken.Access)

	patient, err := s.client.GetPatientProfile(s.ctx, newToken.Access)
	s.Require().NoError(err)

	s.Equal(int64(2), patient.ID)
}

func (s *ProfileTestSuite) TestSelectPatientProfileBadPatient() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	_, _, err := s.client.SelectPatientProfile(s.ctx, s.token.Access, 100)
	s.Require().Error(err)
	s.Equal("unexpected status code: 400, error: invalid input parameter", err.Error())
}

func (s *ProfileTestSuite) TestDeletePatientProfile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeletePatientProfile(s.ctx, s.token.Access, 2)
	s.Require().NoError(err)

	list, err := s.client.GetProfiles(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Patients, 1)
	for _, p := range list.Patients {
		s.NotEqual(int64(2), p.ID)
	}
}

func (s *ProfileTestSuite) TestAddSpecialistProfile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 3,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	about := "Some about text"
	medicalCategory := "1"
	treatsAdults := true
	treatsChildren := false
	req := model.AddSpecialistProfile{
		About:           &about,
		MedicalCategory: &medicalCategory,
		CuresDiseases:   []int64{1, 2, 3},
		Services:        []int64{1, 2, 3},
		TreatsAdults:    &treatsAdults,
		TreatsChildren:  &treatsChildren,
	}

	newToken, cookie, err := s.client.AddSpecialistProfile(s.ctx, *token, &req)
	s.Require().NoError(err)

	s.NotNil(cookie)
	s.NotEmpty(newToken.Access)

	specialist, err := s.client.GetSpecialistProfile(s.ctx, newToken.Access)
	s.Require().NoError(err)

	s.Equal(*req.About, *specialist.About)
	s.Equal(*req.MedicalCategory, *specialist.MedicalCategory)
	s.Equal(req.CuresDiseases, specialist.CuresDiseases)
	s.Equal(req.Services, specialist.Services)
	s.Equal(*req.TreatsAdults, specialist.TreatsAdults)
	s.Equal(*req.TreatsChildren, specialist.TreatsChildren)
}

func (s *ProfileTestSuite) TestAddSpecialistProfileBadReq() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 3,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	var phoneID int64 = 1
	var emailID int64 = 1
	about := "Some about text"
	medicalCategory := "1"
	treatsAdults := true
	treatsChildren := false
	req := model.AddSpecialistProfile{
		PhoneID:         &phoneID,
		EmailID:         &emailID,
		About:           &about,
		MedicalCategory: &medicalCategory,
		CuresDiseases:   []int64{1, 2, 3},
		Services:        []int64{1, 2, 3},
		TreatsAdults:    &treatsAdults,
		TreatsChildren:  &treatsChildren,
	}

	_, _, err = s.client.AddSpecialistProfile(s.ctx, *token, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 400, error: invalid input field", err.Error())
}
