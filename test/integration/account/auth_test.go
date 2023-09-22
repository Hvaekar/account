package account

import (
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/stretchr/testify/suite"
	"net/http"
	"testing"
	"time"
)

type AuthTestSuite struct {
	TestSuite
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (s *AuthTestSuite) TestRegister() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.RegisterRequest{
		Login:    "tester",
		Password: "test_password",
	}

	token, cookie, err := s.client.Register(s.ctx, &req)
	s.Require().NoError(err)

	s.NotNil(cookie)
	s.NotEmpty(token.Access)
}

func (s *AuthTestSuite) TestLogin() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.LoginRequest{
		Login:    "account1",
		Password: "account1",
	}

	token, cookie, err := s.client.Login(s.ctx, &req)
	s.Require().NoError(err)

	s.NotNil(cookie)
	s.NotEmpty(token.Access)
}

func (s *AuthTestSuite) TestLoginBadPassword() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.LoginRequest{
		Login:    "account1",
		Password: "account4",
	}

	_, _, err := s.client.Login(s.ctx, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 400, error: crypto/bcrypt: hashedPassword is not the hash of the given password", err.Error())
}

func (s *AuthTestSuite) TestLoginNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.LoginRequest{
		Login:    "account4",
		Password: "account4",
	}

	_, _, err := s.client.Login(s.ctx, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *AuthTestSuite) TestRefreshToken() {
	payload := model.TokenPayload{
		AccountID:    1,
		PatientID:    1,
		SpecialistID: 1,
	}

	refreshToken, err := jwt.GenerateJWT(s.cfg.JWT.RefreshTokenExpiresAt, payload, s.cfg.JWT.RefreshTokenSecretKey)
	s.Require().NoError(err)

	cookie := http.Cookie{
		Name:     s.cfg.JWT.RefreshTokenCookieName,
		Value:    *refreshToken,
		Expires:  time.Now().Add(time.Hour),
		Path:     "/",
		Domain:   "localhost",
		Secure:   false,
		HttpOnly: true,
	}

	token, err := s.client.RefreshToken(s.ctx, &cookie)
	s.Require().NoError(err)

	s.NotEmpty(token.Access)
}

func (s *AuthTestSuite) TestLogout() {
	cookie, err := s.client.Logout(s.ctx)
	s.Require().NoError(err)

	s.NotNil(cookie)
}
