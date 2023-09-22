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

type EmailTestSuite struct {
	TestSuite
}

func TestEmailSuite(t *testing.T) {
	suite.Run(t, new(EmailTestSuite))
}

func (s *EmailTestSuite) TestAddEmail() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	open := true
	req := model.AddEmail{
		Type:  "personal",
		Email: "some@example.com",
		Open:  &open,
	}

	email, err := s.client.AddEmail(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.NotEmpty(email.ID)
	s.Equal(req.Type, email.Type)
	s.Equal(req.Email, email.Email)
	s.Equal(false, email.Verified)
	s.Equal(*req.Open, email.Open)
}

func (s *EmailTestSuite) TestGetEmails() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetEmails(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Emails, 3)
}

func (s *EmailTestSuite) TestUpdateEmail() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	open := true
	req := model.UpdateEmail{
		Type: "other",
		Open: &open,
	}

	email, err := s.client.UpdateEmail(s.ctx, s.token.Access, 1, &req)
	s.Require().NoError(err)

	s.Equal(int64(1), email.ID)
	s.Equal(req.Type, email.Type)
	s.Equal(*req.Open, email.Open)
}

func (s *EmailTestSuite) TestUpdateEmailNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	open := true
	req := model.UpdateEmail{
		Type: "other",
		Open: &open,
	}

	_, err := s.client.UpdateEmail(s.ctx, s.token.Access, 100, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *EmailTestSuite) TestVerifyEmailCode() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	email, cookie, err := s.client.VerifyEmailCode(s.ctx, s.token.Access, 2)
	s.Require().NoError(err)

	s.Equal(int64(2), email.ID)

	s.NotNil(cookie)
	s.Equal(s.cfg.Verify.VerifyCodeCookieName, cookie.Name)
}

func (s *EmailTestSuite) TestVerifyEmail() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	code := 111111
	req := model.VerifyEmail{
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

	email, err := s.client.VerifyEmail(s.ctx, s.token.Access, 2, &req, &cookie)
	s.Require().NoError(err)

	s.Equal(int64(2), email.ID)
	s.Equal(true, email.Verified)
}

func (s *EmailTestSuite) TestVerifyEmailBadCode() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	code := 111111
	req := model.VerifyEmail{
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

	_, err := s.client.VerifyEmail(s.ctx, s.token.Access, 2, &req, &cookie)
	s.Require().Error(err)
	s.Equal(
		"unexpected status code: 400, error: crypto/bcrypt: hashedPassword is not the hash of the given password",
		err.Error())
}

func (s *EmailTestSuite) TestDeleteEmail() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	err := s.client.DeleteEmail(s.ctx, s.token.Access, 1)
	s.Require().NoError(err)

	list, err := s.client.GetEmails(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Len(list.Emails, 2)
	for _, e := range list.Emails {
		s.NotEqual(int64(1), e.ID)
	}
}
