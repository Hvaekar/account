package model

type RegisterRequest struct {
	Login    string `json:"login" binding:"required,alphanum,min=5,max=100"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type LoginRequest struct {
	Login    string `json:"login" binding:"required,alphanum,min=5,max=100"`
	Password string `json:"password" binding:"required,min=8,max=100"`
}

type LoginResponse struct {
	Data  any   `json:"data"`
	Token Token `json:"token"`
}

//func (r *LoginResponse) ToResponse() IResponse {
//	r.Account.ToResponse()
//	return r
//}
