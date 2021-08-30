package service

import (
	"auth/dao"
	"auth/model"
	"errors"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

const INVALID_USERNAME_PASSWORD string = "Invalid email or password"

func Authenticate(email string, password string) (*model.User, error) {

	user, error := dao.GetUserByEmail(email)

	if error != nil {
		log.Error(error, " ", email)
		err := errors.New(INVALID_USERNAME_PASSWORD)
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.Error("Invalid password, email ", email)
		err = errors.New(INVALID_USERNAME_PASSWORD)
		return nil, err
	}

	if user.Deleted != nil {
		log.Error("User has been deleted, email ", email)
		err := errors.New(INVALID_USERNAME_PASSWORD)
		return nil, err
	}

	return user, nil
}
