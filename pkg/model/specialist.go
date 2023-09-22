package model

import (
	"github.com/Hvaekar/med-account/pkg/storage"
	"time"
)

type Specialist struct {
	ID              int64     `json:"id"`
	UpdatedAt       time.Time `json:"updated_at"`
	FirstName       *string   `json:"first_name,omitempty"`
	FatherName      *string   `json:"father_name,omitempty"`
	LastName        *string   `json:"last_name,omitempty"`
	Sex             *string   `json:"sex,omitempty"`
	Photo           *string   `json:"photo,omitempty"`
	Phone           *string   `json:"phone,omitempty"`
	Email           *string   `json:"email,omitempty"`
	About           *string   `json:"about,omitempty"`
	MedicalCategory *string   `json:"medical_category"`
	CuresDiseases   []int64   `json:"cures_diseases"`
	Services        []int64   `json:"services"`
	TreatsAdults    bool      `json:"treats_adults"`
	TreatsChildren  bool      `json:"treats_children"`
	ListSpecializations
	ListEducations
	ListExperiences
	ListAssociations
	ListPatents
	ListPublicationLinks
	//EducationalCourses []EducationalCourse `json:"educational_courses"`
	//Teaching     []Teaching        `json:"teaching"`
	//Speeches     []Speech          `json:"speeches"`
}

func (s *Specialist) ToResponse() IResponse {
	s.ListSpecializations.ToResponse()
	s.ListEducations.ToResponse()
	s.ListExperiences.ToResponse()
	s.ListAssociations.ToResponse()
	s.ListPatents.ToResponse()
	s.ListPublicationLinks.ToResponse()
	return s
}

type ListSpecialistsRequest struct {
	OrderBy string  `json:"order_by" form:"order_by" url:"order_by" binding:"omitempty,min=1"`
	Limit   uint64  `json:"limit" form:"limit" url:"limit" binding:"omitempty,gt=0"`
	Page    uint64  `json:"page" form:"page" url:"page" binding:"omitempty,gt=0"`
	IDList  []int64 `json:"id_list" form:"page" url:"page" binding:"omitempty,gt=0"`
}

func (l *ListSpecialistsRequest) Prepare() {
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

func (l *ListSpecialistsRequest) Offset() uint64 {
	return l.Limit * (l.Page - 1)
}

type ListSpecialists struct {
	Specialists []*Specialist `json:"specialists"`
}

type AddSpecialistProfile struct {
	PhoneID         *int64  `json:"phone_id" binding:"omitempty,gt=0"`
	EmailID         *int64  `json:"email_id" binding:"omitempty,gt=0"`
	About           *string `json:"about"`
	MedicalCategory *string `json:"medical_category" binding:"omitempty,oneof=0 1 2 3"`
	CuresDiseases   []int64 `json:"cures_diseases" binding:"required"`
	Services        []int64 `json:"services" binding:"required"`
	TreatsAdults    *bool   `json:"treats_adults" binding:"required"`
	TreatsChildren  *bool   `json:"treats_children" binding:"required"`
}

type UpdateSpecialistProfile AddSpecialistProfile

type UpdateSpecialistProfileFields map[string]interface{}

func (f UpdateSpecialistProfileFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"updated_at":       struct{}{},
		"phone_id":         struct{}{},
		"email_id":         struct{}{},
		"about":            struct{}{},
		"medical_category": struct{}{},
		"cures_diseases":   struct{}{},
		"services":         struct{}{},
		"treats_adults":    struct{}{},
		"treats_children":  struct{}{},
	}
}

func (f UpdateSpecialistProfileFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "phone_id", "email_id":
			if v == nil {
				f[k] = storage.NullInt64(nil)
				continue
			}

			val := int64(v.(float64))

			f[k] = storage.NullInt64(&val)
		case "about", "medical_category":
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
