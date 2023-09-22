package model

type ListProfiles struct {
	PatientProfileID    int64             `json:"patient_profile_id"`
	SpecialistProfileID int64             `json:"specialist_profile_id"`
	Patients            []*AccountPatient `json:"patients,omitempty"`
}

type AccountPatient struct {
	ID             int64   `json:"id"`
	FirstName      *string `json:"first_name,omitempty"`
	FatherName     *string `json:"father_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	Photo          *string `json:"photo,omitempty"`
	PermissionEdit bool    `json:"permission_edit"`
	Verified       bool    `json:"verified"`
}

type UpdateAdminFields map[string]interface{}

func (f UpdateAdminFields) DBColumns() map[string]interface{} {
	return map[string]interface{}{
		//"updated_at":  struct{}{},
		"permission_edit": struct{}{},
	}
}

func (f UpdateAdminFields) Prepare() {
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
