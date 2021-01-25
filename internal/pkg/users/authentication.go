package users

import (
	"errors"
	"github.com/andskur/argon2-hashing"
	"github.com/google/uuid"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
)

var (
	UserNotFound      = errors.New("user not found")
	UserAlreadyExists = errors.New("user already exists")
	IncorrectPassword = errors.New("password is incorrect")
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
		"SELECT id, name, mail FROM users WHERE session_key = $1", key)

	if err != nil {
		return nil, err
	}

	users := slice.([]*types.User)
	if len(users) == 0 {
		return nil, UserNotFound
	}

	return users[0], err
}

func Login(mailAddress string, password string) (*types.AuthenticationInformation, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.PwHashInformationType,
		"SELECT password FROM users WHERE mail=$1", mailAddress)
	if err != nil {
		return nil, err
	}

	passwords := slice.([]*types.PasswordHashInformation)
	if len(passwords) == 0 {
		return nil, UserNotFound
	}

	if err = compareHashAndPassword(passwords[0].PasswordHash, password); err != nil {
		return nil, IncorrectPassword
	}

	sessionKey, err := getNewSessionKey()
	if err != nil {
		return nil, err
	}

	slice, err = database.QueryAsync(database.DefaultTimeout, types.AuthInformationType,
		"UPDATE users SET session_key=$1 WHERE mail=$2 RETURNING id, name, session_key;", sessionKey, mailAddress)
	if err != nil {
		return nil, err
	}

	authInfos := slice.([]*types.AuthenticationInformation)
	if len(authInfos) == 0 {
		return nil, UserNotFound
	}

	return authInfos[0], err
}

func Register(name string, mailAddress string, password string) (*types.AuthenticationInformation, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.IdInformationType,
		"SELECT id FROM users WHERE mail = $1 OR name = $2", mailAddress, name)
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
		"INSERT INTO users(name, mail, password, session_key) VALUES($1, $2, $3, $4) RETURNING id, name, session_key",
		name, mailAddress, passwordHash, sessionKey)
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
