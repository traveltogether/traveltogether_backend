package types

import "reflect"

var (
	AuthInformationType   = reflect.TypeOf(&AuthenticationInformation{})
	UserType              = reflect.TypeOf(&User{})
	PwHashInformationType = reflect.TypeOf(&PasswordHashInformation{})
	UsernameType          = reflect.TypeOf(&UsernameInformation{})
	MailInformationType   = reflect.TypeOf(&MailInformation{})
)

type User struct {
	Id                   int     `json:"id" db:"id"`
	Username             string  `json:"username" db:"username"`
	MailAddress          string  `json:"mail" db:"mail"`
	FirstName            *string `json:"first_name,omitempty" db:"first_name"`
	ProfileImageAsBase64 *string `json:"profile_image,omitempty" db:"profile_image"`
	Disabilities         *string `json:"disabilities,omitempty" db:"disabilities"`
}

type AuthenticationInformation struct {
	Id         int    `json:"id" db:"id"`
	Username   string `json:"username" db:"username"`
	SessionKey string `json:"session_key" db:"session_key"`
}

type PasswordHashInformation struct {
	PasswordHash string `db:"password"`
}

type UsernameInformation struct {
	Username string `db:"username"`
}

type MailInformation struct {
	Mail string `db:"mail"`
}
