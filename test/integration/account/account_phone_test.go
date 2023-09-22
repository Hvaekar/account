package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type PhoneTestSuite struct {
	TestSuite
}

func TestPhoneSuite(t *testing.T) {
	suite.Run(t, new(PhoneTestSuite))
}

func (s *PhoneTestSuite) TestAddPhone() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	open := true
	req := model.AddPhone{
		Type:  "personal",
		Code:  "38",
		Phone: "0123456789",
		Open:  &open,
	}

	phone, err := s.client.AddPhone(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.NotEmpty(phone.ID)
	s.Equal(req.Type, phone.Type)
	s.Equal(req.Code, phone.Code)
	s.Equal(req.Phone, phone.Phone)
	s.Equal(false, phone.Verified)
	s.Equal(*req.Open, phone.Open)
}

func (s *PhoneTestSuite) TestGetPhones() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetPhones(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Phones, 3)
}

func (s *PhoneTestSuite) TestUpdatePhone() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	open := true
	req := model.UpdatePhone{
		Type: "other",
		Open: &open,
	}

	phone, err := s.client.UpdatePhone(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(int64(1), phone.ID)
	s.Equal(req.Type, phone.Type)
	s.Equal(*req.Open, phone.Open)
}

func (s *PhoneTestSuite) TestUpdatePhoneNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	open := true
	req := model.UpdatePhone{
		Type: "other",
		Open: &open,
	}

	_, err := s.client.UpdatePhone(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *PhoneTestSuite) TestVerifyPhoneCode() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	phone, cookie, err := s.client.VerifyPhoneCode(s.ctx, s.token.Access, 2)
	s.Require().NoError(err)

	s.Equal(int64(2), phone.ID)

	s.NotNil(cookie)
	s.Equal(s.cfg.Verify.VerifyCodeCookieName, cookie.Name)
}

func (s *PhoneTestSuite) TestVerifyPhone() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	code := 111111
	req := model.VerifyPhone{
		Code: code,
	}

	cookie := http.Cookie{
		Name:     s.cfg.Verify.VerifyCodeCookieName,
		Value:    utils.HashPassword(strconv.Itoa(code)),
		Expires:  time.Now().Add(s.cfg.Verify.VerifyCodeExpiresAt),
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}

	phone, err := s.client.VerifyPhone(s.ctx, s.token.Access, 2, &req, &cookie)
	s.Require().NoError(err)

	s.Equal(int64(2), phone.ID)
	s.Equal(true, phone.Verified)
}

func (s *PhoneTestSuite) TestVerifyPhoneBadCode() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	code := 111111
	req := model.VerifyPhone{
		Code: 123456,
	}

	cookie := http.Cookie{
		Name:     s.cfg.Verify.VerifyCodeCookieName,
		Value:    utils.HashPassword(strconv.Itoa(code)),
		Expires:  time.Now().Add(s.cfg.Verify.VerifyCodeExpiresAt),
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}

	_, err := s.client.VerifyPhone(s.ctx, s.token.Access, 2, &req, &cookie)
	s.Require().Error(err)
	s.Equal(
		"unexpected status code: 400, error: crypto/bcrypt: hashedPassword is not the hash of the given password",
		err.Error())
}

func (s *PhoneTestSuite) TestDeletePhone() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeletePhone(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetPhones(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Phones, 2)
	for _, e := range list.Phones {
		s.NotEqual(int64(1), e.ID)
	}
}
