package model

import (
	"github.com/Hvaekar/med-account/pkg/storage"
)

type Email struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"-"`
	Type      string `json:"type"`
	Email     string `json:"value"`
	Verified  bool   `json:"verified"`
	Open      bool   `json:"open"`
}

func (e *Email) ToResponse() IResponse {
	e.AccountID = 0
	return e
}

type EmailJoin struct {
	ID        *int64  `json:"id"`
	AccountID *int64  `json:"account_id"`
	Type      *string `json:"type"`
	Email     *string `json:"value"`
	Verified  *bool   `json:"verified"`
	Open      *bool   `json:"open"`
}

func (e EmailJoin) ConvertToEmail() Email {
	return Email{
		ID:        *e.ID,
		AccountID: *e.AccountID,
		Type:      *e.Type,
		Email:     *e.Email,
		Verified:  *e.Verified,
		Open:      *e.Open,
	}
}

type AddEmail struct {
	Type  string `json:"type" binding:"required,oneof=personal work other"`
	Email string `json:"email" binding:"required,email"`
	Open  *bool  `json:"open" binding:"required"`
}

type UpdateEmail struct {
	Type string `json:"type" binding:"required,oneof=personal work other"`
	Open *bool  `json:"open" binding:"required"`
}

type VerifyEmail struct {
	Code int `json:"code" binding:"required"`
}

type ListEmails struct {
	Emails []*Email `json:"emails"`
}

func (l *ListEmails) ToResponse() IResponse {
	for _, v := range l.Emails {
		v.ToResponse()
	}
	return l
}

type UpdateEmailFields map[string]interface{}

func (f UpdateEmailFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		//"updated_at":  struct{}{},
		"type":     struct{}{},
		"open":     struct{}{},
		"verified": struct{}{},
	}
}

func (f UpdateEmailFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "open", "verified":
			if v == nil {
				f[k] = storage.NullBool(nil)
				continue
			}

			val := v.(bool)

			f[k] = storage.NullBool(&val)
		default:
			f[k] = v
		}
	}
}
