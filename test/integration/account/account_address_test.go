package account

import (
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AddressTestSuite struct {
	TestSuite
}

func TestAddressSuite(t *testing.T) {
	suite.Run(t, new(AddressTestSuite))
}

func (s *AddressTestSuite) TestAddAddress() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	open := true
	req := model.AddAddress{
		Type:    "personal",
		CityID:  1,
		Address: "Some address string, 1",
		Open:    &open,
	}

	address, err := s.client.AddAddress(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.NotEmpty(address.ID)
	s.Equal(req.Type, address.Type)
	s.Equal(req.CityID, address.CityID)
	s.Equal(req.Address, address.Address)
	s.Equal(*req.Open, address.Open)
}

func (s *AddressTestSuite) TestAddAddressNotFoundAccount() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 100,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	open := true
	req := model.AddAddress{
		Type:    "personal",
		CityID:  1,
		Address: "Some address string, 1",
		Open:    &open,
	}

	_, err = s.client.AddAddress(s.ctx, *token, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *AddressTestSuite) TestGetAddresses() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetAddresses(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Addresses, 2)
}

func (s *AddressTestSuite) TestUpdateAddress() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	open := true
	req := model.UpdateAddress{
		Type: "work",
		Open: &open,
	}

	address, err := s.client.UpdateAddress(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(int64(1), address.ID)
	s.Equal(req.Type, address.Type)
	s.Equal(*req.Open, address.Open)
}

func (s *AddressTestSuite) TestUpdateAddressNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	open := true
	req := model.UpdateAddress{
		Type: "work",
		Open: &open,
	}

	_, err := s.client.UpdateAddress(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *AddressTestSuite) TestDeleteAddress() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteAddress(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetAddresses(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Addresses, 1)
	for _, a := range list.Addresses {
		s.NotEqual(int64(1), a.ID)
	}
}
