package model

type Status struct {
	Success bool `json:"success"`
}

type IDMessage struct {
	ID int64 `json:"id"`
}

type KeyMessage struct {
	Key any `json:"key"`
}

type Verify struct {
	Value string `json:"value"`
	Code  string `json:"code"`
}
