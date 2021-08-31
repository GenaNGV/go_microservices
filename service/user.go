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

const SecretKey = "iTechArtGoLab"

var ErrInvalidUsernamePassword = errors.New("invalid email or password")
var ErrUnauthorized = errors.New("unauthorized")

func Authenticate(email string, password string) (*model.UserAuth, error) {

	user, err := pg.GetUserByEmail(email)

	if err != nil {
		log.Error(err, " ", email)
		return nil, ErrInvalidUsernamePassword
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		log.WithFields(log.Fields{"email": email}).Error("Invalid password")
		return nil, ErrInvalidUsernamePassword
	}

	if user.Deleted != nil {
		log.WithFields(log.Fields{"email": email}).Error("User has been deleted")
		return nil, ErrInvalidUsernamePassword
	}

	// 2 hours
	expired := time.Now().Add(time.Hour * 2)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: expired.Unix(),
	})

	token, _ := claims.SignedString([]byte(SecretKey))

	auth := model.UserAuth{User: *user, Token: token, Expired: expired}

	redis.Put(&auth)

	return &auth, nil
}

func Check(token string) (*model.UserAuth, error) {
	user, err := redis.Get(token)

	if err != nil {
		log.Error(err)
		return nil, ErrUnauthorized
	}

	if user == nil {
		log.Error(err)
		return nil, ErrUnauthorized
	}

	if user.Deleted != nil {
		log.WithFields(log.Fields{"email": user.Email}).Error("User has been deleted")
		return nil, ErrUnauthorized
	}

	if user.Expired.Before(time.Now()) {
		log.WithFields(log.Fields{"email": user.Email}).Error("Session has been expired")
		return nil, ErrUnauthorized
	}

	user.Expired = time.Now().Add(time.Hour * 2)
	redis.Put(user)

	return user, nil
}
