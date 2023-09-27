package account

import (
	"github.com/Hvaekar/med-account/pkg/jwt"
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/test/integration/account/fixtures"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/suite"
	"math"
	"testing"
	"time"
)

type AccountTestSuite struct {
	TestSuite
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountTestSuite))
}

func (s *AccountTestSuite) TestAccountGetMe() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	account, err := s.client.GetMe(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.Equal(int64(1), account.ID)
}

func (s *AccountTestSuite) TestDeleteAccount() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	cookie, err := s.client.DeleteAccount(s.ctx, s.token.Access)
	s.Require().NoError(err)

	s.NotNil(cookie)

	payload := model.TokenPayload{
		AccountID: 2,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	list, err := s.client.GetAccounts(s.ctx, *token, &model.ListAccountsRequest{Limit: math.MaxInt64})
	s.Require().NoError(err)
	for _, a := range list.Accounts {
		if a.ID == 1 {
			s.NotEqual(nil, a.DeletedAt)
		}
	}
}

func (s *AccountTestSuite) TestUpdateAccountMain() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	name := "Another"
	fatherName := "Father"
	lastName := "LastName"
	sex := "woman"
	birthday, err := time.ParseInLocation("2006-01-02", "1980-12-01", time.UTC)
	s.Require().NoError(err)
	language := "af"
	country := "GB"
	req := model.UpdateAccount{
		Login:      "account11",
		FirstName:  &name,
		FatherName: &fatherName,
		LastName:   &lastName,
		Sex:        &sex,
		Birthday:   &pgtype.Date{Time: birthday, Valid: true},
		Language:   &language,
		Country:    &country,
	}

	account, err := s.client.UpdateAccountMain(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.Equal(req.Login, account.Login)
	s.Equal(req.FirstName, account.FirstName)
	s.Equal(req.FatherName, account.FatherName)
	s.Equal(req.LastName, account.LastName)
	s.Equal(req.Sex, account.Sex)
	s.NotEqual(req.Birthday.Time, account.Birthday.Time)
	s.Equal(req.Language, account.Language)
	s.Equal(req.Country, account.Country)
}

func (s *AccountTestSuite) TestUpdateAccount2Main() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 2,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	name := "Another"
	fatherName := "Father"
	lastName := "LastName"
	sex := "woman"
	birthday, err := time.ParseInLocation("2006-01-02", "1980-12-01", time.UTC)
	s.Require().NoError(err)
	language := "af"
	country := "GB"
	req := model.UpdateAccount{
		Login:      "account11",
		FirstName:  &name,
		FatherName: &fatherName,
		LastName:   &lastName,
		Sex:        &sex,
		Birthday:   &pgtype.Date{Time: birthday, Valid: true},
		Language:   &language,
		Country:    &country,
	}

	account, err := s.client.UpdateAccountMain(s.ctx, *token, &req)
	s.Require().NoError(err)

	s.Equal(req.Login, account.Login)
	s.Equal(req.FirstName, account.FirstName)
	s.Equal(req.FatherName, account.FatherName)
	s.Equal(req.LastName, account.LastName)
	s.Equal(req.Sex, account.Sex)
	s.Equal(req.Birthday.Time, account.Birthday.Time)
	s.Equal(req.Language, account.Language)
	s.Equal(req.Country, account.Country)
}

func (s *AccountTestSuite) TestUpdateAccountMainNotFound() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	payload := model.TokenPayload{
		AccountID: 100,
	}

	token, err := jwt.GenerateJWT(s.cfg.JWT.AccessTokenExpiresAt, payload, s.cfg.JWT.AccessTokenSecretKey)
	s.Require().NoError(err)

	name := "Another"
	fatherName := "Father"
	lastName := "LastName"
	sex := "woman"
	birthday, err := time.ParseInLocation("2006-01-02", "1980-12-01", time.UTC)
	s.Require().NoError(err)
	language := "af"
	country := "GB"
	req := model.UpdateAccount{
		Login:      "account11",
		FirstName:  &name,
		FatherName: &fatherName,
		LastName:   &lastName,
		Sex:        &sex,
		Birthday:   &pgtype.Date{Time: birthday, Valid: true},
		Language:   &language,
		Country:    &country,
	}

	_, err = s.client.UpdateAccountMain(s.ctx, *token, &req)
	s.Require().Error(err)
	s.Equal("unexpected status code: 404, error: not found", err.Error())
}

func (s *AccountTestSuite) TestUpdateAccountPassword() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	req := model.UpdatePassword{
		OldPassword: "account1",
		NewPassword: "some_new_password",
	}

	err := s.client.UpdatePassword(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	reqL := model.LoginRequest{
		Login:    "account1",
		Password: req.NewPassword,
	}

	_, _, err = s.client.Login(s.ctx, &reqL)
	s.Require().NoError(err)
}

func (s *AccountTestSuite) TestUpdateAccountPhoto() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	// path as name
	path := "example2.jpg"
	req := model.UpdatePhoto{
		Name: &path,
	}

	account, err := s.client.UpdatePhoto(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.Equal(*req.Name, *account.Photo)

	// path as nil
	req = model.UpdatePhoto{
		Name: nil,
	}

	account, err = s.client.UpdatePhoto(s.ctx, s.token.Access, &req)
	s.Require().NoError(err)

	s.Nil(account.Photo)
}

func (s *AccountTestSuite) TestGetAccounts() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	list, err := s.client.GetAccounts(s.ctx, s.token.Access, &model.ListAccountsRequest{Limit: math.MaxInt64})
	s.Require().NoError(err)
	s.Len(list.Accounts, 3)
	for i, a := range list.Accounts {
		s.Equal(int64(i+1), a.ID)
		s.NotEqual(nil, a.DeletedAt)
	}
}

func (s *AccountTestSuite) TestGetAccount() {
	s.Require().NoError(s.db.TruncateTables(s.ctx, truncateTables...))
	s.Require().NoError(fixtures.PopulateDB(s.ctx, s.db.GetDB()))

	account, err := s.client.GetAccount(s.ctx, s.token.Access, 2)
	s.Require().NoError(err)
	s.Require().Equal("account2", account.Login)
}
