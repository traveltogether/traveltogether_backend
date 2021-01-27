package users

import "errors"

var (
	UserNotFound      = errors.New("user not found")
	UserAlreadyExists = errors.New("user already exists")
	IncorrectPassword = errors.New("password is incorrect")
	MailAlreadyInUse  = errors.New("mail is already in use")
)
