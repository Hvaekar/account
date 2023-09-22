package model

type IUpdateFields interface {
	DBColumns() map[string]interface{}
	Prepare()
}

type IResponse interface {
	ToResponse() IResponse
}
