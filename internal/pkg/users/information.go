package users

import (
	"github.com/traveltogether/traveltogether_backend/internal/pkg/database"
	"github.com/traveltogether/traveltogether_backend/internal/pkg/types"
)

func ChangeUsername(user *types.User, newUsername string) error {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.UsernameType,
		"SELECT username FROM users WHERE username = $1", newUsername)
	if err != nil {
		return err
	}

	users := slice.([]*types.UsernameInformation)
	if len(users) != 0 {
		return UserAlreadyExists
	}

	err = database.PrepareAsync(database.DefaultTimeout, "UPDATE users SET username = $1 WHERE id = $2",
		newUsername, user.Id)
	if err != nil {
		return err
	}

	user.Username = newUsername
	return nil
}

func ChangeProfileImage(user *types.User, newProfileImage string) error {
	err := database.PrepareAsync(database.DefaultTimeout, "UPDATE users SET profile_image = $1 WHERE id = $2",
		newProfileImage, user.Id)
	if err != nil {
		return err
	}

	user.ProfileImageAsBase64 = &newProfileImage
	return nil
}

func ChangeDisabilities(user *types.User, newDisabilities string) error {
	err := database.PrepareAsync(database.DefaultTimeout, "UPDATE users SET disabilities = $1 WHERE id = $2",
		newDisabilities, user.Id)
	if err != nil {
		return err
	}

	user.Disabilities = &newDisabilities
	return nil
}

func ChangeFirstname(user *types.User, newFirstname string) error {
	err := database.PrepareAsync(database.DefaultTimeout, "UPDATE users SET first_name = $1 WHERE id = $2",
		newFirstname, user.Id)
	if err != nil {
		return err
	}

	user.FirstName = &newFirstname
	return nil
}

func ChangeMailAddress(user *types.User, newMailAddress string) error {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.MailInformationType,
		"SELECT mail FROM users WHERE mail = $1", newMailAddress)
	if err != nil {
		return err
	}

	users := slice.([]*types.MailInformation)
	if len(users) != 0 {
		return MailAlreadyInUse
	}

	err = database.PrepareAsync(database.DefaultTimeout, "UPDATE users SET mail = $1 WHERE id = $2",
		newMailAddress, user.Id)
	if err != nil {
		return err
	}

	user.MailAddress = &newMailAddress
	return nil
}

func ChangePassword(user *types.User, oldPassword string, newPassword string) error {
	passwordHash, err := getUserPasswordHash(user.Username)
	if err != nil {
		return err
	}

	if err = compareHashAndPassword(passwordHash, oldPassword); err != nil {
		return IncorrectPassword
	}

	newPasswordHash, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	return database.PrepareAsync(database.DefaultTimeout, "UPDATE users SET password = $1 WHERE id = $2",
		newPasswordHash, user.Id)
}

func GetUserById(id int) (*types.User, error) {
	slice, err := database.QueryAsync(database.DefaultTimeout, types.UserType,
		"SELECT id, username, mail, first_name, profile_image, disabilities FROM users WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	users := slice.([]*types.User)
	if len(users) == 0 {
		return nil, UserNotFound
	}

	return users[0], err
}
