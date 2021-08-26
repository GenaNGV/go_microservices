package service

import (
	"auth/dao"
	"auth/model"
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(email string, password string) (*model.User, error) {

	user := new(model.User)

	if res := dao.DB.Where(&model.User{Email: email}).First(&user); res.RowsAffected <= 0 {
		err := errors.New("Invalid Credentials")
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		err = errors.New("Invalid Credentials")
		return nil, err
	}

	return user, nil
}
