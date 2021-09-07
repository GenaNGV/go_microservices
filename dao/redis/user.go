package redis

import (
	"auth/enviroment"
	"auth/model"
	"encoding/json"
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

func Put(u *model.UserAuth) {
	enviroment.Env.RDB.Set(u.Token, u, time.Duration(u.Expired.Unix()))
}

func Get(token string) (*model.UserAuth, error) {

	val := enviroment.Env.RDB.Get(token)
	if val == nil {
		return nil, ErrUserNotFound
	} else {
		userAuth := &model.UserAuth{}
		err := json.Unmarshal([]byte(val.Val()), &userAuth)
		if err != nil {
			return nil, err
		}
		return userAuth, nil
	}
}
