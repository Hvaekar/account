package model

import (
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Experience struct {
	ID              int64        `json:"id"`
	ProfileID       int64        `json:"-"`
	CompanyID       *int64       `json:"company_id,omitempty"`
	Company         string       `json:"company"`
	Start           pgtype.Date  `json:"start"`
	Finish          *pgtype.Date `json:"finish,omitempty"`
	Specializations []int64      `json:"specializations"`
}

func (e *Experience) ToResponse() IResponse {
	e.ProfileID = 0
	return e
}

type ExperienceJoin struct {
	ID        *int64       `json:"id"`
	ProfileID *int64       `json:"profile_id"`
	CompanyID *int64       `json:"company_id"`
	Company   *string      `json:"company"`
	Start     *pgtype.Date `json:"start"`
	Finish    *pgtype.Date `json:"finish"`
}

func (e ExperienceJoin) ConvertToExperience() Experience {
	return Experience{
		ID:        *e.ID,
		ProfileID: *e.ProfileID,
		CompanyID: e.CompanyID,
		Company:   *e.Company,
		Start:     *e.Start,
		Finish:    e.Finish,
	}
}

type AddExperience struct {
	CompanyID       *int64       `json:"company_id" binding:"omitempty,gt=0"`
	Company         string       `json:"company" binding:"required,max=255"`
	Start           pgtype.Date  `json:"start" binding:"required"`
	Finish          *pgtype.Date `json:"finish"`
	Specializations []int64      `json:"specializations" binding:"required"`
}

type UpdateExperience AddExperience

type ListExperiences struct {
	Experiences []*Experience `json:"experiences"`
}

func (l *ListExperiences) ToResponse() IResponse {
	for _, v := range l.Experiences {
		v.ToResponse()
	}
	return l
}

type UpdateExperienceFields map[string]interface{}

func (f UpdateExperienceFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"company_id": struct{}{},
		"company":    struct{}{},
		"start":      struct{}{},
		"finish":     struct{}{},
	}
}

func (f UpdateExperienceFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "company_id":
			if v == nil {
				f[k] = storage.NullInt64(nil)
				continue
			}

			val := int64(v.(float64))

			f[k] = storage.NullInt64(&val)
		case "finish":
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
