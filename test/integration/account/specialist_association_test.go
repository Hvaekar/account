package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AssociationTestSuite struct {
	TestSuite
}

func TestAssociationSuite(t *testing.T) {
	suite.Run(t, new(AssociationTestSuite))
}

func (s *AssociationTestSuite) TestAddAssociation() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var associationID int64 = 100
	jobTitle := "Some job title"
	req := model.AddAssociation{
		AssociationID: &associationID,
		Name:          "Some association name",
		JobTitle:      &jobTitle,
	}

	association, err := s.client.AddAssociation(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.NotEmpty(association.ID)
	s.Equal(*req.AssociationID, *association.AssociationID)
	s.Equal(req.Name, association.Name)
	s.Equal(*req.JobTitle, *association.JobTitle)
}

func (s *AssociationTestSuite) TestGetAssociations() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetAssociations(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Associations, 3)
}

func (s *AssociationTestSuite) TestUpdateAssociation() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var associationID int64 = 100
	jobTitle := "Some another job title"
	req := model.UpdateAssociation{
		AssociationID: &associationID,
		Name:          "Some another association name",
		JobTitle:      &jobTitle,
	}

	association, err := s.client.UpdateAssociation(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(*req.AssociationID, *association.AssociationID)
	s.Equal(req.Name, association.Name)
	s.Equal(*req.JobTitle, *association.JobTitle)
}

func (s *AssociationTestSuite) TestUpdateAssociationNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	var associationID int64 = 100
	jobTitle := "Some another job title"
	req := model.UpdateAssociation{
		AssociationID: &associationID,
		Name:          "Some another association name",
		JobTitle:      &jobTitle,
	}

	_, err := s.client.UpdateAssociation(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *AssociationTestSuite) TestDeleteAssociation() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteAssociation(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetAssociations(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Associations, 2)
	for _, v := range list.Associations {
		s.NotEqual(int64(1), v.ID)
	}
}
