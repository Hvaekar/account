package model

import (
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/jackc/pgx/v5/pgtype"
)

type Education struct {
	ID            int64       `json:"id"`
	ProfileID     int64       `json:"-"`
	InstitutionID int64       `json:"institution_id"`
	FacultyID     *int64      `json:"faculty_id,omitempty"`
	DepartmentID  *int64      `json:"department_id,omitempty"`
	FormID        *int64      `json:"form_id,omitempty"`
	DegreeID      *int64      `json:"degree_id,omitempty"`
	Graduation    pgtype.Date `json:"graduation"`
	Verified      bool        `json:"verified"`
	Files         []*File     `json:"files"`
}

func (e *Education) ToResponse() IResponse {
	e.ProfileID = 0
	for _, v := range e.Files {
		v.ToResponse()
	}
	return e
}

type EducationJoin struct {
	ID            *int64       `json:"id"`
	ProfileID     *int64       `json:"profile_id"`
	InstitutionID *int64       `json:"institution_id"`
	FacultyID     *int64       `json:"faculty_id"`
	DepartmentID  *int64       `json:"department_id"`
	FormID        *int64       `json:"form_id"`
	DegreeID      *int64       `json:"degree_id"`
	Graduation    *pgtype.Date `json:"graduation"`
	Verified      *bool        `json:"verified"`
}

func (e EducationJoin) ConvertToEducation() Education {
	return Education{
		ID:            *e.ID,
		ProfileID:     *e.ProfileID,
		InstitutionID: *e.InstitutionID,
		FacultyID:     e.FacultyID,
		DepartmentID:  e.DepartmentID,
		FormID:        e.FormID,
		DegreeID:      e.DegreeID,
		Graduation:    *e.Graduation,
		Verified:      *e.Verified,
	}
}

type AddEducation struct {
	InstitutionID int64       `json:"institution_id" binding:"required,gt=0"`
	FacultyID     *int64      `json:"faculty_id" binding:"omitempty,gt=0"`
	DepartmentID  *int64      `json:"department_id" binding:"omitempty,gt=0"`
	FormID        *int64      `json:"form_id" binding:"omitempty,gt=0"`
	DegreeID      *int64      `json:"degree_id" binding:"omitempty,gt=0"`
	Graduation    pgtype.Date `json:"graduation"`
	Files         []*File     `json:"files"`
}

type AddEducations struct {
	Educations []AddEducation `json:"educations" binding:"required"`
}

type UpdateEducation struct {
	InstitutionID int64       `json:"institution_id" binding:"required,gt=0"`
	FacultyID     *int64      `json:"faculty_id" binding:"omitempty,gt=0"`
	DepartmentID  *int64      `json:"department_id" binding:"omitempty,gt=0"`
	FormID        *int64      `json:"form_id" binding:"omitempty,gt=0"`
	DegreeID      *int64      `json:"degree_id" binding:"omitempty,gt=0"`
	Graduation    pgtype.Date `json:"graduation"`
	Files         []*File     `json:"files"`
}

type ListEducations struct {
	Educations []*Education `json:"educations"`
}

func (l *ListEducations) ToResponse() IResponse {
	for _, v := range l.Educations {
		v.ToResponse()
	}
	return l
}

type UpdateEducationFields map[string]interface{}

func (f UpdateEducationFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"institution_id": struct{}{},
		"faculty_id":     struct{}{},
		"department_id":  struct{}{},
		"form_id":        struct{}{},
		"degree_id":      struct{}{},
		"graduation":     struct{}{},
	}
}

func (f UpdateEducationFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "faculty_id", "department_id", "form_id", "degree_id":
			if v == nil {
				f[k] = storage.NullInt64(nil)
				continue
			}

			val := int64(v.(float64))

			f[k] = storage.NullInt64(&val)
		default:
			f[k] = v
		}
	}
}
