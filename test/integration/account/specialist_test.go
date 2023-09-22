package account

import (
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SpecialistTestSuite struct {
	TestSuite
}

func TestSpecialistSuite(t *testing.T) {
	suite.Run(t, new(SpecialistTestSuite))
}

func (s *SpecialistTestSuite) TestGetSpecialistProfile() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	specialist, err := s.client.GetSpecialistProfile(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Equal(int64(1), specialist.ID)
}

func (s *SpecialistTestSuite) TestUpdateSpecialistProfileMain() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var phoneID int64 = 2
	var emailID int64 = 2
	about := "Some another about text"
	medicalCategory := "1"
	treatsAdults := false
	treatsChildren := true
	req := model.UpdateSpecialistProfile{
		PhoneID:         &phoneID,
		EmailID:         &emailID,
		About:           &about,
		MedicalCategory: &medicalCategory,
		CuresDiseases:   []int64{10, 20, 30, 40},
		Services:        []int64{10, 20, 30, 40},
		TreatsAdults:    &treatsAdults,
		TreatsChildren:  &treatsChildren,
	}

	specialist, err := s.client.UpdateSpecialist(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.Equal(*req.About, *specialist.About)
	s.Equal(*req.MedicalCategory, *specialist.MedicalCategory)
	s.Len(req.CuresDiseases, 4)
	for _, v := range req.CuresDiseases {
		s.Contains(specialist.CuresDiseases, v)
	}
	s.Len(req.Services, 4)
	for _, v := range req.Services {
		s.Contains(specialist.Services, v)
	}
	s.Equal(*req.TreatsAdults, specialist.TreatsAdults)
	s.Equal(*req.TreatsChildren, specialist.TreatsChildren)
}

func (s *SpecialistTestSuite) TestSpecialistProfileMainNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID:    3,
		SpecialistID: 3,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	var phoneID int64 = 2
	var emailID int64 = 2
	about := "Some another about text"
	medicalCategory := "1"
	treatsAdults := false
	treatsChildren := true
	req := model.UpdateSpecialistProfile{
		PhoneID:         &phoneID,
		EmailID:         &emailID,
		About:           &about,
		MedicalCategory: &medicalCategory,
		CuresDiseases:   []int64{10, 20, 30, 40},
		Services:        []int64{10, 20, 30, 40},
		TreatsAdults:    &treatsAdults,
		TreatsChildren:  &treatsChildren,
	}

	_, err = s.client.UpdateSpecialist(s.ctx, *token, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *SpecialistTestSuite) TestGetSpecialists() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.ListSpecialistsRequest{}

	list, err := s.client.GetSpecialists(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)
	s.Len(list.Specialists, 2)
	for i, p := range list.Specialists {
		s.Equal(int64(i+1), p.ID)
	}

	req.IDList = []int64{2}

	list, err = s.client.GetSpecialists(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)
	s.Len(list.Specialists, 1)
	for i, p := range list.Specialists {
		s.Equal(int64(i+2), p.ID)
	}
}

func (s *SpecialistTestSuite) TestGetSpecialist() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	specialist, err := s.client.GetSpecialist(s.ctx, s.token.Access, 2)
	s.Require().NoError(err)

	s.Equal(int64(2), specialist.ID)
}
