package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AdminTestSuite struct {
	TestSuite
}

func TestAdminSuite(t *testing.T) {
	suite.Run(t, new(AdminTestSuite))
}

func (s *AdminTestSuite) TestAddAdmin() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	pEdit := true
	req := model.AddAdmin{
		AdminID:        3,
		PermissionEdit: &pEdit,
	}

	admin, err := s.client.AddAdmin(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.Equal(int64(3), admin.ID)
	s.Equal(*req.PermissionEdit, admin.PermissionEdit)
	s.Nil(admin.FirstName)
	s.Nil(admin.FatherName)
	s.Nil(admin.LastName)
	s.Nil(admin.Photo)
	s.False(admin.Verified)
}

func (s *AdminTestSuite) TestAddAdminBadReq() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	pEdit := true
	req := model.AddAdmin{
		AdminID:        1,
		PermissionEdit: &pEdit,
	}

	_, err := s.client.AddAdmin(s.ctx, s.token.Access, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 400, error: invalid input parameter", err.Error())
}

func (s *AdminTestSuite) TestGetAdmins() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetAdmins(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Admins, 1)

	// add one and check
	pEdit := true
	req := model.AddAdmin{
		AdminID:        3,
		PermissionEdit: &pEdit,
	}

	_, err = s.client.AddAdmin(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	list, err = s.client.GetAdmins(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Admins, 2)
}

func (s *AdminTestSuite) TestUpdateAdmin() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	pEdit := true
	req := model.UpdateAdmin{
		PermissionEdit: &pEdit,
	}

	admin, err := s.client.UpdateAdmin(s.ctx, s.token.Access, 2, &req)
	s.Require().NoError(err)

	s.Equal(int64(2), admin.ID)
	s.Equal(*req.PermissionEdit, admin.PermissionEdit)
}

func (s *AdminTestSuite) TestUpdateAdminNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	pEdit := true
	req := model.UpdateAdmin{
		PermissionEdit: &pEdit,
	}

	_, err := s.client.UpdateAdmin(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *AdminTestSuite) TestDeleteAdmin() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteAdmin(s.ctx, s.token.Access, 2)
	s.Require().NoError(err)

	list, err := s.client.GetAdmins(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Admins, 0)
}
