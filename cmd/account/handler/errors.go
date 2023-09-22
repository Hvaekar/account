package handler

import "errors"

var (
	ErrNoPermissions = errors.New("you have no permissions here")
	ErrMissingParam  = errors.New("missing incoming parameter")
	ErrInvalidParam  = errors.New("invalid input parameter")
	ErrInvalidField  = errors.New("invalid input field")
	ErrNoToken       = errors.New("request does not contain a token")
	ErrLimitExceeded = errors.New("limit is exceeded")
)
