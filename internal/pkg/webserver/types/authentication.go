package types

type LoginData struct {
	UsernameOrMail string `json:"nameOrMail"`
	Password       string `json:"password"`
}

type RegistrationData struct {
	Username string `json:"name"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
}
