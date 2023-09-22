package account

import (
	"github.com/Hvaekar/med-account/pkg/model"
	"github.com/Hvaekar/med-account/pkg/storage/postgres"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/gin-gonic/gin"
)

func (s *PostgresStorage) Register(c *gin.Context, req *model.RegisterRequest) (*model.Account, error) {
	psql := s.SetFormat().RunWith(s.DB)

	q := psql.Insert(accountsTableName).
		Columns(
			"login",
			"password",
		).
		Values(
			req.Login,
			utils.HashPassword(req.Password),
		).
		Suffix("RETURNING \"id\"")

	var id int64
	if err := q.QueryRowContext(c).Scan(&id); err != nil {
		return nil, postgres.ConvertError(err)
	}

	_ = psql.Insert(patientProfilesTableName).Columns("account_id").Values(id).QueryRowContext(c)

	return s.GetAccountByID(c, id)
}

func (s *PostgresStorage) Login(c *gin.Context, req *model.LoginRequest) (*model.Account, error) {
	a, err := s.GetAccountByLogin(c, req.Login)
	if err != nil {
		return nil, err
	}

	if err := utils.ValidatePassword(a.Password, req.Password); err != nil {
		return nil, err
	}

	a.Profiles.Patients, err = s.GetPatientProfiles(c, a.ID)
	if err != nil {
		return nil, err
	}

	return a, nil
}
