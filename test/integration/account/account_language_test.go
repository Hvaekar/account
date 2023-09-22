package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type LanguageTestSuite struct {
	TestSuite
}

func TestLanguageSuite(t *testing.T) {
	suite.Run(t, new(LanguageTestSuite))
}

func (s *LanguageTestSuite) TestAddLanguage() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.AddLanguage{
		Language: "aa",
		Level:    "a1",
	}

	language, err := s.client.AddLanguage(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.Equal(req.Language, language.Language)
	s.Equal(req.Level, language.Level)
}

func (s *LanguageTestSuite) TestGetLanguages() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetLanguages(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Languages, 2)
}

func (s *LanguageTestSuite) TestUpdateLanguage() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.UpdateLanguage{
		Level: "a1",
	}

	language, err := s.client.UpdateLanguage(s.ctx, s.token.Access, "en", &req)
	s.Require().NoError(err)

	s.Equal(req.Level, language.Level)
}

func (s *LanguageTestSuite) TestUpdateLanguageNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.UpdateLanguage{
		Level: "a1",
	}

	_, err := s.client.UpdateLanguage(s.ctx, s.token.Access, "aa", &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *LanguageTestSuite) TestDeleteLanguage() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteLanguage(s.ctx, s.token.Access, "en")
	s.Require().NoError(err)

	list, err := s.client.GetLanguages(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Languages, 1)
	for _, l := range list.Languages {
		s.NotEqual("en", l.Language)
	}
}
