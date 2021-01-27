package types

type LoginData struct {
	UsernameOrMail string `json:"usernameOrMail"`
	Password       string `json:"password"`
}

type RegistrationData struct {
	Username             string  `json:"username"`
	Mail                 string  `json:"mail"`
	Password             string  `json:"password"`
	FirstName            *string `json:"first_name,omitempty"`
	Disabilities         *string `json:"disabilities,omitempty"`
	ProfileImageAsBase64 *string `json:"profile_image,omitempty"`
}
