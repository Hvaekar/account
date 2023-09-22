package model

import "github.com/Hvaekar/med-account/pkg/storage"

type Address struct {
	ID        int64  `json:"id"`
	AccountID int64  `json:"-"`
	Type      string `json:"type"`
	CityID    int64  `json:"city_id"`
	Address   string `json:"address"`
	Open      bool   `json:"open"`
}

func (a *Address) ToResponse() IResponse {
	a.AccountID = 0
	return a
}

type AddressJoin struct {
	ID        *int64  `json:"id"`
	AccountID *int64  `json:"account_id"`
	Type      *string `json:"type"`
	CityID    *int64  `json:"city_id"`
	Address   *string `json:"address"`
	Open      *bool   `json:"open"`
}

func (a AddressJoin) ConvertToAddress() Address {
	return Address{
		ID:        *a.ID,
		AccountID: *a.AccountID,
		Type:      *a.Type,
		CityID:    *a.CityID,
		Address:   *a.Address,
		Open:      *a.Open,
	}
}

type AddAddress struct {
	Type    string `json:"type" binding:"required,oneof=personal work other"`
	CityID  int64  `json:"city_id" binding:"required,gt=0"`
	Address string `json:"address" binding:"required,max=255"`
	Open    *bool  `json:"open" binding:"required"`
}

type UpdateAddress struct {
	Type string `json:"type" binding:"required,oneof=personal work other"`
	Open *bool  `json:"open" binding:"required"`
}

type ListAddresses struct {
	Addresses []*Address `json:"addresses"`
}

func (l *ListAddresses) ToResponse() IResponse {
	for _, v := range l.Addresses {
		v.ToResponse()
	}
	return l
}

type UpdateAddressFields map[string]interface{}

func (f UpdateAddressFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		//"updated_at":  struct{}{},
		"type": struct{}{},
		"open": struct{}{},
	}
}

func (f UpdateAddressFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "open":
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
