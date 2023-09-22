package model

import "github.com/Hvaekar/med-account/pkg/storage"

type Phone struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"-"`
	Type      string `json:"type"`
	Code      string `json:"code"`
	Phone     string `json:"value"`
	Verified  bool   `json:"verified"`
	Open      bool   `json:"open"`
}

func (p *Phone) ToResponse() IResponse {
	p.AccountID = 0
	return p
}

type PhoneJoin struct {
	ID        *int64  `json:"id"`
	AccountID *int64  `json:"account_id"`
	Type      *string `json:"type"`
	Code      *string `json:"code"`
	Phone     *string `json:"value"`
	Verified  *bool   `json:"verified"`
	Open      *bool   `json:"open"`
}

func (p PhoneJoin) ConvertToPhone() Phone {
	return Phone{
		ID:        *p.ID,
		AccountID: *p.AccountID,
		Type:      *p.Type,
		Code:      *p.Code,
		Phone:     *p.Phone,
		Verified:  *p.Verified,
		Open:      *p.Open,
	}
}

type AddPhone struct {
	Type  string `json:"type" binding:"required,oneof=personal work other"`
	Code  string `json:"code" binding:"required,max=20"`
	Phone string `json:"phone" binding:"required,max=20"`
	Open  *bool  `json:"open" binding:"required"`
}

type UpdatePhone struct {
	Type string `json:"type" binding:"required,oneof=personal work other"`
	Open *bool  `json:"open" binding:"required"`
}

type VerifyPhone struct {
	Code int `json:"code" binding:"required"`
}

type ListPhones struct {
	Phones []*Phone `json:"phones"`
}

func (l *ListPhones) ToResponse() IResponse {
	for _, v := range l.Phones {
		v.ToResponse()
	}
	return l
}

type UpdatePhoneFields map[string]interface{}

func (f UpdatePhoneFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		//"updated_at":  struct{}{},
		"type":     struct{}{},
		"open":     struct{}{},
		"verified": struct{}{},
	}
}

func (f UpdatePhoneFields) Prepare() {
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
