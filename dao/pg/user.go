package pg

import (
	"auth/enviroment"
	"auth/model"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

func GetUserByEmail(email string) (*model.User, error) {

	user := &model.User{}

	if res := enviroment.Env.DB.Preload("Roles").Where("email", email).Take(&user); res.RowsAffected <= 0 {
		return nil, ErrUserNotFound
	}

	return user, nil
}
