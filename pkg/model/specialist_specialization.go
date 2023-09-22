package model

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Specialization struct {
	SpecializationID int64       `json:"specialization_id"`
	Start            pgtype.Date `json:"start"`
}

func (s *Specialization) ToResponse() IResponse {
	return s
}

type SpecializationJoin struct {
	SpecializationID *int64       `json:"specialization_id"`
	Start            *pgtype.Date `json:"start"`
}

func (s SpecializationJoin) ConvertToSpecialization() Specialization {
	return Specialization{
		SpecializationID: *s.SpecializationID,
		Start:            *s.Start,
	}
}

type AddSpecialization struct {
	SpecializationID int64       `json:"specialization_id" binding:"required,gt=0"`
	Start            pgtype.Date `json:"start" binding:"required"`
}

type AddSpecializations struct {
	Specializations []AddSpecialization `json:"specializations" binding:"required"`
}

type UpdateSpecialization struct {
	Start pgtype.Date `json:"start" binding:"required"`
}

type ListSpecializations struct {
	Specializations []*Specialization `json:"specializations"`
}

func (l *ListSpecializations) ToResponse() IResponse {
	for _, v := range l.Specializations {
		v.ToResponse()
	}
	return l
}

type UpdateSpecializationFields map[string]interface{}

func (f UpdateSpecializationFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"start": struct{}{},
	}
}

func (f UpdateSpecializationFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		default:
			f[k] = v
		}
	}
}
