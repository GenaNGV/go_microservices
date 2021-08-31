package service

import (
	"auth/dao/pg"
	"auth/dao/redis"
	"auth/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

const INVALID_USERNAME_PASSWORD string = "Invalid email or password"
const UNAUTORIZED string = "Unauthorized"
const SECRET_KEY = "iTechArtGoLab"

func Authenticate(email string, password string) (*model.UserAuth, error) {

	user, error := pg.GetUserByEmail(email)

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

	// 2 hours
	expired := time.Now().Add(time.Hour * 2)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: expired.Unix(),
	})

	token, _ := claims.SignedString([]byte(SECRET_KEY))

	auth := model.UserAuth{User: *user, Token: token, Expired: expired}

	redis.Put(&auth)

	return &auth, nil
}

func Check(token string) (*model.UserAuth, error) {
	user, error := redis.Get(token)

	if error != nil {
		log.Error(error)
		err := errors.New(UNAUTORIZED)
		return nil, err
	}

	if user == nil {
		log.Error(error)
		err := errors.New(UNAUTORIZED)
		return nil, err
	}

	if user.Deleted != nil {
		log.Error("User has been deleted, email ", user.Email)
		err := errors.New(UNAUTORIZED)
		return nil, err
	}

	if user.Expired.Before(time.Now()) {
		log.Error("Session has been expired, email ", user.Email)
		err := errors.New(UNAUTORIZED)
		return nil, err
	}

	user.Expired = time.Now().Add(time.Hour * 2)
	redis.Put(user)

	return user, nil
}
