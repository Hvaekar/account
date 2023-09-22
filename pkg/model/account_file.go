package model

import (
	"github.com/Hvaekar/med-account/pkg/storage"
	"time"
)

type File struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
	AccountID   int64     `json:"-"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
}

func (f *File) ToResponse() IResponse {
	f.CreatedAt = time.Time{}
	f.UpdatedAt = time.Time{}
	f.AccountID = 0
	return f
}

type FileJoin struct {
	ID          *int64     `json:"id"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	AccountID   *int64     `json:"account_id"`
	Name        *string    `json:"name"`
	Description *string    `json:"description"`
}

func (f FileJoin) ConvertToFile() File {
	return File{
		ID:          *f.ID,
		CreatedAt:   *f.CreatedAt,
		UpdatedAt:   *f.UpdatedAt,
		AccountID:   *f.AccountID,
		Name:        *f.Name,
		Description: f.Description,
	}
}

type AddFile struct {
	Path string `json:"path" binding:"required"`
}

type UpdateFile struct {
	Description *string `json:"description" binding:"omitempty,max=255"`
}

type ListFiles struct {
	Files []*File `json:"files"`
}

func (l *ListFiles) ToResponse() IResponse {
	for _, v := range l.Files {
		v.ToResponse()
	}
	return l
}

type UpdateFileFields map[string]interface{}

func (f UpdateFileFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		"updated_at":  struct{}{},
		"description": struct{}{},
	}
}

func (f UpdateFileFields) Prepare() {
	for k, v := range f {
		columns := f.DBColumns()
		if _, ok := columns[k]; !ok {
			delete(f, k)
			continue
		}

		switch k {
		case "description":
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

func MatchingUniqueFiles(input []*File) []*File {
	if input == nil || len(input) == 0 {
		return nil
	}

	u := make([]*File, 0, len(input))
	m := make(map[int64]bool)

	for _, val := range input {
		if _, ok := m[val.ID]; !ok {
			m[val.ID] = true
			u = append(u, val)
		}
	}

	return u
}
