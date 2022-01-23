package models

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not found")
	ErrAlreadyExist = errors.New("already exist")
	ErrConflict     = errors.New("conflict")
)

type FieldError struct {
	Field string
	Code  string
}

type InternalError string

func (e InternalError) Error() string {
	return string(e)
}

type BadRequest struct {
	Msg    string `json:"-"`
	Errors []FieldError
}

func (e BadRequest) Error() string {
	return e.Msg
}