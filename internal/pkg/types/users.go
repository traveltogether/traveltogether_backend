package types

import "reflect"

var (
	AuthInformationType   = reflect.TypeOf(&AuthenticationInformation{})
	UserType              = reflect.TypeOf(&User{})
	PwHashInformationType = reflect.TypeOf(&PasswordHashInformation{})
)

type User struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	MailAddress string `json:"mail" db:"mail"`
}

type AuthenticationInformation struct {
	Id         int    `json:"id" db:"id"`
	Name       string `json:"name" db:"name"`
	SessionKey string `json:"session_key" db:"session_key"`
}

type PasswordHashInformation struct {
	PasswordHash string `db:"password"`
}
