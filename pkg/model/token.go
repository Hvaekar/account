package model

type Token struct {
	Access  string `json:"access_token"`
	Refresh string `json:"refresh_token,omitempty"`
}

type TokenPayload struct {
	AccountID    int64 `json:"account_id"`
	PatientID    int64 `json:"patient_id"`
	SpecialistID int64 `json:"specialist_id"`
}

type IdentifyRequest struct {
	Account    bool `json:"account"`
	Patient    bool `json:"patient"`
	Specialist bool `json:"specialist"`
}
