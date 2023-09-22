package model

import (
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/Hvaekar/med-account/pkg/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Account struct {
	ID         int64        `json:"id"`
	CreatedAt  time.Time    `json:"created_at"`
	UpdatedAt  time.Time    `json:"updated_at"`
	DeletedAt  *time.Time   `json:"deleted_at,omitempty"`
	Login      string       `json:"login"`
	Password   string       `json:"-"`
	FirstName  *string      `json:"first_name,omitempty"`
	FatherName *string      `json:"father_name,omitempty"`
	LastName   *string      `json:"last_name,omitempty"`
	Sex        *string      `json:"sex,omitempty"`
	Photo      *string      `json:"photo,omitempty"`
	Birthday   *pgtype.Date `json:"birthday,omitempty"`
	Language   *string      `json:"language,omitempty"`
	Country    *string      `json:"country,omitempty"`

	ListEmails
	ListPhones
	ListAddresses
	ListLanguages

	Profiles ListProfiles `json:"profiles"`
	//Companies      []Company           `json:"companies"`
}

func (a *Account) ToResponse() IResponse {
	a.Password = ""
	a.ListEmails.ToResponse()
	a.ListPhones.ToResponse()
	a.ListAddresses.ToResponse()
	a.ListLanguages.ToResponse()
	return a
}

type ListAccountsRequest struct {
	OrderBy string `json:"order_by" form:"order_by" url:"order_by" binding:"omitempty,min=1"`
	Limit   uint64 `json:"limit" form:"limit" url:"limit" binding:"omitempty,gt=0"`
	Page    uint64 `json:"page" form:"page" url:"page" binding:"omitempty,gt=0"`
}

func (l *ListAccountsRequest) Prepare() {
	if l.OrderBy == "" {
		l.OrderBy = defaultOrderBy
	}
	if l.Limit == 0 {
		l.Limit = defaultLimit
	}
	if l.Page == 0 {
		l.Page = defaultPage
	}
}

func (l *ListAccountsRequest) Offset() uint64 {
	return l.Limit * (l.Page - 1)
}

type ListAccounts struct {
	Accounts []*Account `json:"accounts"`
}

type UpdateAccount struct {
	Login      string       `json:"login" binding:"required,min=5,max=100"`
	FirstName  *string      `json:"first_name,omitempty" binding:"omitempty,max=100"`
	FatherName *string      `json:"father_name,omitempty" binding:"omitempty,max=100"`
	LastName   *string      `json:"last_name,omitempty" binding:"omitempty,max=100"`
	Sex        *string      `json:"sex,omitempty" binding:"omitempty,oneof=man woman"`
	Birthday   *pgtype.Date `json:"birthday,omitempty"`
	Language   *string      `json:"language,omitempty"`
	Country    *string      `json:"country,omitempty" binding:"omitempty,iso3166_1_alpha2"`
}

type UpdatePassword struct {
	OldPassword string `json:"old_password" binding:"required,min=8,max=100"`
	NewPassword string `json:"new_password" binding:"required,min=8,max=100"`
}

type UpdatePhoto struct {
	Path *string `json:"path"`
}

type UpdateAccountFields map[string]interface{}

func (f UpdateAccountFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"updated_at":  struct{}{},
		"deleted_at":  struct{}{},
		"login":       struct{}{},
		"password":    struct{}{},
		"first_name":  struct{}{},
		"father_name": struct{}{},
		"last_name":   struct{}{},
		"sex":         struct{}{},
		"photo":       struct{}{},
		"birthday":    struct{}{},
		"language":    struct{}{},
		"country":     struct{}{},
	}
}

func (f UpdateAccountFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "password":
			f[k] = utils.HashPassword(v.(string))
		case "first_name", "father_name", "last_name", "sex", "language", "country", "photo":
			if v == nil {
				f[k] = storage.NullString(nil)
				continue
			}

			val := v.(*string)

			f[k] = storage.NullString(val)
		case "deleted_at":
			if v == nil {
				f[k] = storage.NullTime(nil)
				continue
			}

			val := v.(time.Time)

			f[k] = storage.NullTime(&val)
		case "birthday":
			if v == nil {
				f[k] = storage.NullDatePGX(nil)
				continue
			}

			t, err := time.Parse("2006-01-02", v.(string))
			if err != nil {
				f[k] = storage.NullDatePGX(nil)
				continue
			}

			val := pgtype.Date{
				Time:  t,
				Valid: true,
			}

			f[k] = storage.NullDatePGX(&val)
		default:
			f[k] = v
		}
	}
}
