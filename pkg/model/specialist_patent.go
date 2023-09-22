package model

import "github.com/Hvaekar/med-account/pkg/storage"

type Patent struct {
	ID        int64   `json:"id"`
	ProfileID int64   `json:"-"`
	Number    string  `json:"number"`
	Name      string  `json:"name"`
	Link      *string `json:"link,omitempty"`
}

func (p *Patent) ToResponse() IResponse {
	p.ProfileID = 0
	return p
}

type PatentJoin struct {
	ID        *int64  `json:"id"`
	ProfileID *int64  `json:"profile_id"`
	Number    *string `json:"number"`
	Name      *string `json:"name"`
	Link      *string `json:"link"`
}

func (e PatentJoin) ConvertToPatent() Patent {
	return Patent{
		ID:        *e.ID,
		ProfileID: *e.ProfileID,
		Number:    *e.Number,
		Name:      *e.Name,
		Link:      e.Link,
	}
}

type AddPatent struct {
	Number string  `json:"number" binding:"required,max=100"`
	Name   string  `json:"name" binding:"required,max=255"`
	Link   *string `json:"link" binding:"omitempty,http_url,max=255"`
}

type UpdatePatent AddPatent

type ListPatents struct {
	Patents []*Patent `json:"patents"`
}

func (l *ListPatents) ToResponse() IResponse {
	for _, v := range l.Patents {
		v.ToResponse()
	}
	return l
}

type UpdatePatentFields map[string]interface{}

func (f UpdatePatentFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"number": struct{}{},
		"name":   struct{}{},
		"link":   struct{}{},
	}
}

func (f UpdatePatentFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "link":
			if v == nil {
				f[k] = storage.NullString(nil)
				continue
			}

			val := v.(string)

			f[k] = storage.NullString(&val)
		default:
			f[k] = v
		}
	}
}
