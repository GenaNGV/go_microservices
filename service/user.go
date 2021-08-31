package service

import (
	"auth/dao"
	"auth/model"
	"errors"
	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

const INVALID_USERNAME_PASSWORD string = "Invalid email or password"
const SECRET_KEY = "iTechArtGoLab"

func Authenticate(email string, password string) (*model.UserAuth, error) {

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

	// 2 hours
	expired := time.Now().Add(time.Hour * 2).Unix()
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: expired,
	})

	token, _ := claims.SignedString([]byte(SECRET_KEY))

	auth := model.UserAuth{UserDetail: *user, Token: token, Expired: expired}

	return &auth, nil
}
