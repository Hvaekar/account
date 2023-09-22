package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PatentTestSuite struct {
	TestSuite
}

func TestPatentSuite(t *testing.T) {
	suite.Run(t, new(PatentTestSuite))
}

func (s *PatentTestSuite) TestAddPatent() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	link := "https://example.com/some"
	req := model.AddPatent{
		Number: "UA5775738992",
		Name:   "Some patent name",
		Link:   &link,
	}

	patent, err := s.client.AddPatent(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.NotEmpty(patent.ID)
	s.Equal(req.Number, patent.Number)
	s.Equal(req.Name, patent.Name)
	s.Equal(*req.Link, *patent.Link)
}

func (s *PatentTestSuite) TestGetPatents() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetPatents(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Patents, 2)
}

func (s *PatentTestSuite) TestUpdatePatent() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	link := "https://example.com/some"
	req := model.UpdatePatent{
		Number: "UA5775738992",
		Name:   "Some another patent name",
		Link:   &link,
	}

	patent, err := s.client.UpdatePatent(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(req.Number, patent.Number)
	s.Equal(req.Name, patent.Name)
	s.Equal(*req.Link, *patent.Link)
}

func (s *PatentTestSuite) TestUpdatePatentNotUnique() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	link := "https://example.com/some"
	req := model.UpdatePatent{
		Number: "UA123456789",
		Name:   "Some another patent name",
		Link:   &link,
	}

	_, err := s.client.UpdatePatent(s.ctx, s.token.Access, 1, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 500, error: unique constraint fail", err.Error())
}

func (s *PatentTestSuite) TestUpdatePatentNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	link := "https://example.com/some"
	req := model.UpdatePatent{
		Number: "UA5775738992",
		Name:   "Some another patent name",
		Link:   &link,
	}

	_, err := s.client.UpdatePatent(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *PatentTestSuite) TestDeletePatent() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeletePatent(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetPatents(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Patents, 1)
	for _, e := range list.Patents {
		s.NotEqual(int64(1), e.ID)
	}
}
