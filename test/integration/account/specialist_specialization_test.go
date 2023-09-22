package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type SpecializationTestSuite struct {
	TestSuite
}

func TestSpecializationSuite(t *testing.T) {
	suite.Run(t, new(SpecializationTestSuite))
}

func (s *SpecializationTestSuite) TestAddSpecialization() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	start, err := time.ParseInLocation("2006-01-02", "2010-04-23", time.UTC)
	s.Require().NoError(err)
	req := model.AddSpecialization{
		SpecializationID: 100,
		Start:            pgtype.Date{Time: start, Valid: true},
	}

	specialization, err := s.client.AddSpecialization(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.Equal(req.SpecializationID, specialization.SpecializationID)
	s.Equal(req.Start.Time, specialization.Start.Time)
}

func (s *SpecializationTestSuite) TestGetSpecializations() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetSpecializations(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Specializations, 3)
}

func (s *SpecializationTestSuite) TestUpdateSpecialization() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	start, err := time.ParseInLocation("2006-01-02", "2010-04-23", time.UTC)
	s.Require().NoError(err)
	req := model.UpdateSpecialization{
		Start: pgtype.Date{Time: start, Valid: true},
	}

	specialization, err := s.client.UpdateSpecialization(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(int64(1), specialization.SpecializationID)
	s.Equal(req.Start.Time, specialization.Start.Time)
}

func (s *SpecializationTestSuite) TestUpdateSpecializationNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	start, err := time.ParseInLocation("2006-01-02", "2010-04-23", time.UTC)
	s.Require().NoError(err)
	req := model.UpdateSpecialization{
		Start: pgtype.Date{Time: start, Valid: true},
	}

	_, err = s.client.UpdateSpecialization(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *SpecializationTestSuite) TestDeleteSpecialization() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteSpecialization(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetSpecializations(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Specializations, 2)
	for _, v := range list.Specializations {
		s.NotEqual(int64(1), v.SpecializationID)
	}
}
