package auth_utils

import "errors"

var (
	ErrInvalidFirstName  = errors.New("please provide a valid first name")
	ErrInvalidLastName   = errors.New("please provide a valid last name")
	ErrInvalidEmail      = errors.New("please provide a valid email address")
	ErrInvalidPhone      = errors.New("please provide a valid phone number")
	ErrInvalidUsername   = errors.New("please provide a valid username")
	ErrInvalidPassword   = errors.New("please provide a valid password")
	ErrInvalidMiddleName = errors.New("please provide a valid middle name")

	ErrUserNotFound      = errors.New("user could not be found")
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUsernameTaken     = errors.New("this username is taken")
	ErrEmailTaken        = errors.New("this email is taken")

	ErrUserBanned          = errors.New("user is banned")
	ErrUserNotVerified     = errors.New("user is not verified")
	ErrUnauthorizedRequest = errors.New("unauthorized request")

	ErrNoUsersFound = errors.New("users could not be found")
)
