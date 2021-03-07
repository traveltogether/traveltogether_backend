package general

import (
	"regexp"
	"strings"
)

var mailRegex = regexp.MustCompile(
	"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func IsEmailValid(mail string) (bool, string) {
	mail = strings.TrimSpace(mail)
	if len(mail) < 3 && len(mail) > 254 {
		return false, ""
	}

	return mailRegex.MatchString(mail), mail
}

func IsNameValid(name string) (bool, string) {
	name = strings.TrimSpace(name)
	if len(name) == 0 {
		return false, ""
	}

	if strings.Contains(name, "@") {
		return false, ""
	}

	return true, name
}
