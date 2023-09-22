package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PublicationLinkTestSuite struct {
	TestSuite
}

func TestPublicationLinkSuite(t *testing.T) {
	suite.Run(t, new(PublicationLinkTestSuite))
}

func (s *PublicationLinkTestSuite) TestAddPublicationLink() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.AddPublicationLink{
		Title: "Some title",
		Link:  "https://example.com/some",
	}

	pl, err := s.client.AddPublicationLink(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.NotEmpty(pl.ID)
	s.Equal(req.Title, pl.Title)
	s.Equal(req.Link, pl.Link)
}

func (s *PublicationLinkTestSuite) TestGetPublicationLinks() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetPublicationLinks(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.PublicationLinks, 2)
}

func (s *PublicationLinkTestSuite) TestUpdatePublicationLink() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.UpdatePublicationLink{
		Title: "Some another title",
		Link:  "https://example.com/some-anotehr",
	}

	pl, err := s.client.UpdatePublicationLink(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(req.Title, pl.Title)
	s.Equal(req.Link, pl.Link)
}

func (s *PublicationLinkTestSuite) TestUpdatePublicationLinkNotUnique() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.UpdatePublicationLink{
		Title: "Some title",
		Link:  "https://example.com/link2",
	}

	_, err := s.client.UpdatePublicationLink(s.ctx, s.token.Access, 1, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 500, error: unique constraint fail", err.Error())
}

func (s *PublicationLinkTestSuite) TestUpdatePublicationLinkNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.UpdatePublicationLink{
		Title: "Some title",
		Link:  "https://example.com/some",
	}

	_, err := s.client.UpdatePublicationLink(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *PublicationLinkTestSuite) TestDeletePublicationLink() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeletePublicationLink(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetPublicationLinks(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.PublicationLinks, 1)
	for _, e := range list.PublicationLinks {
		s.NotEqual(int64(1), e.ID)
	}
}
