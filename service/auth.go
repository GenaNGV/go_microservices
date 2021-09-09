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
const TokenTTL = time.Hour * 1

var ErrInvalidUsernamePassword = errors.New("invalid email or password")
var ErrUnauthorized = errors.New("unauthorized")
var ErrSystem = errors.New("oops, we will fix it quickly")

func Login(email string, password string) (*model.UserAuth, error) {

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
	expired := time.Now().Add(TokenTTL)
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.Id)),
		ExpiresAt: expired.Unix(),
	})

	token, _ := claims.SignedString([]byte(SecretKey))

	auth := model.UserAuth{User: *user, Token: token, Expired: expired}

	if err := redis.Put(&auth, time.Duration(TokenTTL)); err != nil {
		log.WithFields(log.Fields{"email": email}).Error("System Exception")
		return nil, ErrSystem
	}

	return &auth, nil
}

func TokenDetail(token string) (*model.UserAuth, error) {

	if token == "" {
		return nil, ErrUnauthorized
	}

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

	setAuthTokenTTL(user)
	_ = redis.Put(user, time.Duration(TokenTTL))

	return user, nil
}

func setAuthTokenTTL(u *model.UserAuth) {
	u.Expired = time.Now().Add(TokenTTL)
}
