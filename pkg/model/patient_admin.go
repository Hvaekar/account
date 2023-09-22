package model

type PatientAdmin struct {
	ID             int64   `json:"id"`
	FirstName      *string `json:"first_name,omitempty"`
	FatherName     *string `json:"father_name,omitempty"`
	LastName       *string `json:"last_name,omitempty"`
	Photo          *string `json:"photo,omitempty"`
	PermissionEdit bool    `json:"permission_edit"`
	Verified       bool    `json:"verified"`
}

type ListAdmins struct {
	Admins []*PatientAdmin `json:"admins,omitempty"`
}

type UpdateAdmin struct {
	PermissionEdit *bool `json:"permission_edit" binding:"required"`
}

type PatientAdminMessage struct {
	AdminID        int64 `json:"account_id"`
	PatientID      int64 `json:"patient_id"`
	PermissionEdit bool  `json:"permission_edit"`
}
