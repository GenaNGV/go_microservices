package redis

import (
	"auth/enviroment"
	"auth/model"
	"errors"
	"time"
)

func Put(u *model.UserAuth) {
	enviroment.Env.RDB.Set(u.Token, u, time.Duration(u.Expired.Unix()))
}

func Get(token string) (*model.UserAuth, error) {
	err := errors.New("not implemented")
	return nil, err
}
