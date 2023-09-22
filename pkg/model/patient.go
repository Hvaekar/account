package model

import (
	"github.com/Hvaekar/med-account/pkg/storage"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	PatientAdminPermissionEdit = "edit"
)

type Patient struct {
	ID         int64        `json:"id"`
	AccountID  int64        `json:"-"`
	FirstName  *string      `json:"first_name,omitempty"`
	FatherName *string      `json:"father_name,omitempty"`
	LastName   *string      `json:"last_name,omitempty"`
	Sex        *string      `json:"sex,omitempty"`
	Photo      *string      `json:"photo,omitempty"`
	Birthday   *pgtype.Date `json:"birthday,omitempty"`
	Phone      *string      `json:"phone,omitempty" `
	Email      *string      `json:"email,omitempty"`
	Body
	Blood
	Vision
	Disability
	LifeStyle
	ListMetalComponents
	ListAdmins
}

func (p *Patient) ToResponse() IResponse {
	p.ListMetalComponents.ToResponse()
	for _, v := range p.Disability.Files {
		v.ToResponse()
	}
	return p
}

type Body struct {
	Height   *float64 `json:"height,omitempty"`
	Weight   *float64 `json:"weight,omitempty"`
	BodyType *string  `json:"body_type,omitempty"`
}

type Blood struct {
	BloodType *string `json:"blood_type,omitempty"`
	Rh        *bool   `json:"rh,omitempty"`
}

type Vision struct {
	LeftEye  *float64 `json:"left_eye,omitempty"`
	RightEye *float64 `json:"right_eye,omitempty"`
}

type Disability struct {
	Group       *string `json:"disability_group,omitempty"`
	Reason      *string `json:"disability_reason,omitempty"`
	DocumentNum *string `json:"disability_document_num,omitempty"`
	Files       []*File `json:"disability_files,omitempty"`
}

type LifeStyle struct {
	Activity  *string `json:"activity,omitempty"`
	Nutrition *string `json:"nutrition,omitempty"`
	Work      *string `json:"work,omitempty"`
}

type ListPatientsRequest struct {
	OrderBy string  `json:"order_by" form:"order_by" url:"order_by" binding:"omitempty,min=1"`
	Limit   uint64  `json:"limit" form:"limit" url:"limit" binding:"omitempty,gt=0"`
	Page    uint64  `json:"page" form:"page" url:"page" binding:"omitempty,gt=0"`
	IDList  []int64 `json:"id_list" form:"page" url:"page" binding:"omitempty,gt=0"`
}

func (l *ListPatientsRequest) Prepare() {
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

func (l *ListPatientsRequest) Offset() uint64 {
	return l.Limit * (l.Page - 1)
}

type ListPatients struct {
	Patients []*Patient `json:"patients"`
}

type AccountPatientMessage struct {
	AdminID        int64 `json:"account_id"`
	PatientID      int64 `json:"patient_id"`
	PermissionEdit bool  `json:"permission_edit"`
}

type AddAdmin struct {
	AdminID        int64 `json:"admin_id" binding:"required,gt=0"`
	PermissionEdit *bool `json:"permission_edit"`
}

type UpdatePatientProfile struct {
	PhoneID               *int64   `json:"phone_id,omitempty" binding:"omitempty,gt=0"`
	EmailID               *int64   `json:"email_id,omitempty" binding:"omitempty,gt=0"`
	Height                *float64 `json:"height,omitempty" binding:"omitempty,gte=10,lte=300"`
	Weight                *float64 `json:"weight,omitempty" binding:"omitempty,gte=0,lte=600"`
	BodyType              *string  `json:"body_type,omitempty" binding:"omitempty,oneof=ectomorph mesomorph endomorph other"`
	BloodType             *string  `json:"blood_type,omitempty" binding:"omitempty,oneof=A B AB O"`
	Rh                    *bool    `json:"rh,omitempty"`
	LeftEye               *float64 `json:"left_eye,omitempty" binding:"omitempty,gte=-40,lte=10"`
	RightEye              *float64 `json:"right_eye,omitempty" binding:"omitempty,gte=-40,lte=10"`
	DisabilityGroup       *string  `json:"disability_group,omitempty" binding:"omitempty,max=1"`
	DisabilityReason      *string  `json:"disability_reason,omitempty" binding:"omitempty,max=255"`
	DisabilityDocumentNum *string  `json:"disability_document_num,omitempty" binding:"omitempty,max=100"`
	DisabilityFiles       []*File  `json:"disability_files,omitempty"`
	Activity              *string  `json:"activity,omitempty" binding:"omitempty,max=255"`
	Nutrition             *string  `json:"nutrition,omitempty" binding:"omitempty,max=255"`
	Work                  *string  `json:"work,omitempty" binding:"omitempty,max=255"`
}

type UpdatePatientProfileFields map[string]interface{}

func (f UpdatePatientProfileFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"updated_at":              struct{}{},
		"phone_id":                struct{}{},
		"email_id":                struct{}{},
		"height":                  struct{}{},
		"weight":                  struct{}{},
		"body_type":               struct{}{},
		"rh":                      struct{}{},
		"left_eye":                struct{}{},
		"right_eye":               struct{}{},
		"disability_group":        struct{}{},
		"disability_reason":       struct{}{},
		"disability_document_num": struct{}{},
		"activity":                struct{}{},
		"nutrition":               struct{}{},
		"work":                    struct{}{},
	}
}

func (f UpdatePatientProfileFields) Prepare() {
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
		case "height", "weight", "left_eye", "right_eye":
			if v == nil {
				f[k] = storage.NullFloat64(nil)
				continue
			}

			val := v.(float64)

			f[k] = storage.NullFloat64(&val)
		case "body_type", "blood_type", "disability_group", "disability_reason", "disability_document_num", "activity", "nutrition", "work":
			if v == nil {
				f[k] = storage.NullString(nil)
				continue
			}

			val := v.(string)

			f[k] = storage.NullString(&val)
		case "rh":
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
