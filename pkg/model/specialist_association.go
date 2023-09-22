package model

import (
	"github.com/Hvaekar/med-account/pkg/storage"
)

type Association struct {
	ID            int64   `json:"id"`
	ProfileID     int64   `json:"-"`
	AssociationID *int64  `json:"association_id,omitempty"`
	Name          string  `json:"name"`
	JobTitle      *string `json:"job_title,omitempty"`
}

func (a *Association) ToResponse() IResponse {
	a.ProfileID = 0
	return a
}

type AssociationJoin struct {
	ID            *int64  `json:"id"`
	ProfileID     *int64  `json:"profile_id"`
	AssociationID *int64  `json:"association_id"`
	Name          *string `json:"name"`
	JobTitle      *string `json:"job_title"`
}

func (e AssociationJoin) ConvertToAssociation() Association {
	return Association{
		ID:            *e.ID,
		ProfileID:     *e.ProfileID,
		AssociationID: e.AssociationID,
		Name:          *e.Name,
		JobTitle:      e.JobTitle,
	}
}

type AddAssociation struct {
	AssociationID *int64  `json:"association_id" binding:"omitempty,gt=0"`
	Name          string  `json:"name" binding:"required,max=255"`
	JobTitle      *string `json:"job_title" binding:"omitempty,max=255"`
}

type UpdateAssociation AddAssociation

type ListAssociations struct {
	Associations []*Association `json:"associations"`
}

func (l *ListAssociations) ToResponse() IResponse {
	for _, v := range l.Associations {
		v.ToResponse()
	}
	return l
}

type UpdateAssociationFields map[string]interface{}

func (f UpdateAssociationFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"association_id": struct{}{},
		"name":           struct{}{},
		"job_title":      struct{}{},
	}
}

func (f UpdateAssociationFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "association_id":
			if v == nil {
				f[k] = storage.NullInt64(nil)
				continue
			}

			val := int64(v.(float64))

			f[k] = storage.NullInt64(&val)
		case "job_title":
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
