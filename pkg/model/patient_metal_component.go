package model

import "github.com/Hvaekar/med-account/pkg/storage"

type MetalComponent struct {
	ID          int64   `json:"id"`
	PatientID   int64   `json:"-"`
	Metal       *string `json:"metal,omitempty"`
	OrganID     int64   `json:"organ_id,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (c *MetalComponent) ToResponse() IResponse {
	c.PatientID = 0
	return c
}

type AddMetalComponent struct {
	Metal       *string `json:"metal,omitempty" binding:"omitempty,max=100"`
	OrganID     int64   `json:"organ_id,omitempty" binding:"required,gt=0"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=255"`
}

type MetalComponentJoin struct {
	ID          *int64  `json:"id"`
	PatientID   *int64  `json:"patient_id"`
	Metal       *string `json:"metal,omitempty"`
	OrganID     *int64  `json:"organ_id,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (m MetalComponentJoin) ConvertToMetalComponent() MetalComponent {
	return MetalComponent{
		ID:          *m.ID,
		PatientID:   *m.PatientID,
		Metal:       m.Metal,
		OrganID:     *m.OrganID,
		Description: m.Description,
	}
}

type UpdateMetalComponent struct {
	Metal       *string `json:"metal,omitempty" binding:"omitempty,max=100"`
	OrganID     int64   `json:"organ_id,omitempty" binding:"required,gt=0"`
	Description *string `json:"description,omitempty" binding:"omitempty,max=255"`
}

type ListMetalComponents struct {
	MetalComponents []*MetalComponent `json:"metal_components"`
}

func (c *ListMetalComponents) ToResponse() IResponse {
	for _, v := range c.MetalComponents {
		v.ToResponse()
	}
	return c
}

type UpdateMetalComponentFields map[string]interface{}

func (f UpdateMetalComponentFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		//"updated_at":              struct{}{},
		"metal":       struct{}{},
		"organ_id":    struct{}{},
		"description": struct{}{},
	}
}

func (f UpdateMetalComponentFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "metal", "description":
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
