package dao

import (
	"auth/model"
	"errors"
)

func GetUserByEmail(email string) (*model.User, error) {

	user := new(model.User)

	if res := DB.Preload("Roles").Where(&model.User{Email: email}).Take(&user); res.RowsAffected <= 0 {
		err := errors.New("Email not found")
		return nil, err
	}

	return user, nil
}
