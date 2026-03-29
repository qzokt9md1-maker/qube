package service

import "errors"

var (
	ErrUnauthorizedAction = errors.New("you are not authorized to perform this action")
	ErrNotFound           = errors.New("resource not found")
	ErrAlreadyExists      = errors.New("resource already exists")
	ErrSelfAction         = errors.New("cannot perform this action on yourself")
	ErrBlocked            = errors.New("user is blocked")
)
