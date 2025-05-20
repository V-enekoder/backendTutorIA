package user

import (
	"errors"
	"time"

	"github.com/V-enekoder/backendTutorIA/src/schema"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserService(userDTO UserCreateDTO) (uint, error) {
	if exists, err := UserExistsByFieldService("email", userDTO.Email, 0); err != nil {
		return 0, err
	} else if exists {
		return 0, HandleUniquenessError("email")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	user := schema.User{
		Name:     userDTO.Name,
		Email:    userDTO.Email,
		Password: string(hashedPassword),
	}

	id, err := CreateUserRepository(user)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func UserExistsByFieldService(field string, value interface{}, excludeId uint) (bool, error) {
	return UserExistsByFieldRepository(field, value, excludeId)
}

func HandleUniquenessError(type_ string) error {
	switch type_ {
	case "email":
		return errors.New("Correo ya registrado")
	default:
		return nil
	}
}

func GetUserByIdService(id uint) (UserResponseDTO, error) {
	user, err := GetUserByIdRepository(id)
	if err != nil {
		return UserResponseDTO{}, err
	}

	userResponse := UserResponseDTO{
		ID:                user.ID,
		Name:              user.Name,
		Email:             user.Email,
		NumberOfDocuments: uint(len(user.Documents)),
		NumberOfProjects:  uint(len(user.Projects)),
		CreatedAt:         user.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         user.UpdatedAt.Format(time.RFC3339),
	}

	return userResponse, nil
}

func UpdatePasswordUserService(id uint, password UserUpdatePasswordDTO) error {
	dbPassword, err := GetPasswordUserRepository(id)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password.OldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(password.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = UpdatePasswordUserRepository(id, string(hashedNewPassword))
	return err
}

func UpdateUserService(id uint, userDTO UserUpdateDTO) error {
	err := UpdateUserRepository(id, userDTO)
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserByIdService(id uint) error {
	err := DeleteUserbyIDRepository(id)
	if err != nil {
		return err
	}

	return nil
}
