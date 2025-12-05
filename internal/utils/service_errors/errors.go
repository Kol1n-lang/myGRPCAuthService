package service_errors

import "errors"

var (
	InternalServerError     = errors.New("internal server error")
	UserAlreadyExistsError  = errors.New("user already exists")
	UserNotFoundError       = errors.New("user not found")
	InvalidCredentialsError = errors.New("invalid credentials")
)
