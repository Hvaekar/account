package account

import (
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type MetalComponentTestSuite struct {
	TestSuite
}

func TestMetalComponentSuite(t *testing.T) {
	suite.Run(t, new(MetalComponentTestSuite))
}

func (s *MetalComponentTestSuite) TestAddMetalComponent() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	metal := "some metal"
	description := "some another description"
	req := model.AddMetalComponent{
		Metal:       &metal,
		OrganID:     100,
		Description: &description,
	}

	mc, err := s.client.AddMetalComponent(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.NotEmpty(mc.ID)
	s.Equal(*req.Metal, *mc.Metal)
	s.Equal(req.OrganID, req.OrganID)
	s.Equal(*req.Description, *mc.Description)
}

func (s *MetalComponentTestSuite) TestAddMetalComponentNoPermission() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 1,
		PatientID: 3,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	metal := "some metal"
	description := "some another description"
	req := model.AddMetalComponent{
		Metal:       &metal,
		OrganID:     100,
		Description: &description,
	}

	_, err = s.client.AddMetalComponent(s.ctx, *token, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 400, error: you have no permissions here", err.Error())
}

func (s *MetalComponentTestSuite) TestGetMetalComponents() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetMetalComponents(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.MetalComponents, 2)
}

func (s *MetalComponentTestSuite) TestUpdateMetalComponent() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	metal := "some metal"
	description := "some another description"
	req := model.UpdateMetalComponent{
		Metal:       &metal,
		OrganID:     100,
		Description: &description,
	}

	mc, err := s.client.UpdateMetalComponent(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(*req.Metal, *mc.Metal)
	s.Equal(req.OrganID, mc.OrganID)
	s.Equal(*req.Description, *mc.Description)
}

func (s *MetalComponentTestSuite) TestUpdateMetalComponentNoPermission() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 1,
		PatientID: 3,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	metal := "some metal"
	description := "some another description"
	req := model.UpdateMetalComponent{
		Metal:       &metal,
		OrganID:     100,
		Description: &description,
	}

	_, err = s.client.UpdateMetalComponent(s.ctx, *token, 4, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 400, error: you have no permissions here", err.Error())
}

func (s *MetalComponentTestSuite) TestUpdateMetalComponentNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	metal := "some metal"
	description := "some another description"
	req := model.UpdateMetalComponent{
		Metal:       &metal,
		OrganID:     100,
		Description: &description,
	}

	_, err := s.client.UpdateMetalComponent(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *MetalComponentTestSuite) TestDeleteMetalComponent() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteMetalComponent(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetMetalComponents(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.MetalComponents, 1)
	for _, e := range list.MetalComponents {
		s.NotEqual(int64(1), e.ID)
	}
}

func (s *MetalComponentTestSuite) TestDeleteMetalComponentNoPermission() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 1,
		PatientID: 3,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	err = s.client.DeleteMetalComponent(s.ctx, *token, 4)
	s.Require().Error(err)
	s.Equal("unexpected status code: 400, error: you have no permissions here", err.Error())

	list, err := s.client.GetMetalComponents(s.ctx, *token)
	s.Require().NoError(err)

	s.Len(list.MetalComponents, 1)
}
