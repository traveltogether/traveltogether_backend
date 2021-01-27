package users

import (
	"github.com/andskur/argon2-hashing"
	"github.com/google/uuid"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
)

func hashPassword(password string) (string, error) {
	hash, err := argon2.GenerateFromPassword([]byte(password), argon2.DefaultParams)

	return string(hash), err
}

func compareHashAndPassword(hash string, password string) error {
	return argon2.CompareHashAndPassword([]byte(hash), []byte(password))
}

func GetUserByAuthenticationKey(key string) (*types.User, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.UserType,
		"SELECT id, username, mail, first_name, profile_image, disabilities FROM users WHERE session_key = $1", key)

	if err != nil {
		return nil, err
	}

	users := slice.([]*types.User)
	if len(users) == 0 {
		return nil, UserNotFound
	}

	return users[0], err
}

func Login(nameOrMail string, password string) (*types.AuthenticationInformation, error) {
	passwordHash, err := getUserPasswordHash(nameOrMail)
	if err != nil {
		return nil, err
	}

	if err = compareHashAndPassword(passwordHash, password); err != nil {
		return nil, IncorrectPassword
	}

	sessionKey, err := getNewSessionKey()
	if err != nil {
		return nil, err
	}

	slice, err := database.QueryAsync(database.DefaultTimeout, types.AuthInformationType,
		"UPDATE users SET session_key = $1 WHERE mail = $2 OR username = $3 RETURNING id, username, session_key",
		sessionKey, nameOrMail, nameOrMail)
	if err != nil {
		return nil, err
	}

	authInfos := slice.([]*types.AuthenticationInformation)
	if len(authInfos) == 0 {
		return nil, UserNotFound
	}

	return authInfos[0], err
}

func Register(name string, mailAddress string, password string, firstname *string, profileImage *string,
	disabilities *string) (*types.AuthenticationInformation, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
		"SELECT id FROM users WHERE mail = $1 OR username = $2", mailAddress, name)
	if err != nil {
		return nil, err
	}

	if len(slice.([]*types.IdInformation)) != 0 {
		return nil, UserAlreadyExists
	}

	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	sessionKey, err := getNewSessionKey()
	if err != nil {
		return nil, err
	}

	slice, err = database.QueryAsync(database.DefaultTimeout, types.AuthInformationType,
		"INSERT INTO users(username, mail, first_name, password, session_key, profile_image, disabilities) "+
			"VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id, username, session_key",
		name, mailAddress, firstname, passwordHash, sessionKey, profileImage, disabilities)
	if err != nil {
		return nil, err
	}

	authInfos := slice.([]*types.AuthenticationInformation)
	if len(authInfos) == 0 {
		return nil, UserNotFound
	}

	return authInfos[0], err
}

func getNewSessionKey() (string, error) {
	var key string
	var err error

	for {
		key = uuid.New().String()

		slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
			"SELECT id FROM users WHERE session_key = $1", key)
		if err != nil {
			break
		}

		if len(slice.([]*types.IdInformation)) == 0 {
			break
		}
	}

	return key, err
}

func getUserPasswordHash(nameOrMail string) (string, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.PwHashInformationType,
		"SELECT password FROM users WHERE mail = $1 OR username = $2", nameOrMail, nameOrMail)
	if err != nil {
		return "", err
	}

	passwords := slice.([]*types.PasswordHashInformation)
	if len(passwords) == 0 {
		return "", UserNotFound
	}

	return passwords[0].PasswordHash, nil
}
