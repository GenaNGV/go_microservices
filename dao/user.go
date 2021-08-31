package dao

import (
	"auth/enviroment"
	"auth/model"
	"errors"
)

func GetUserByEmail(email string) (*model.User, error) {

	user := new(model.User)

	if res := enviroment.DB.Preload("Roles").Where("email", email).Take(&user); res.RowsAffected <= 0 {
		err := errors.New("user not found")
		return nil, err
	}

	return user, nil
}
