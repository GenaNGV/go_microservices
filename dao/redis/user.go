package redis

import (
	"auth/enviroment"
	"auth/model"
	"encoding/json"
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

func Put(u *model.UserAuth, duration time.Duration) error {
	entry, err := json.Marshal(u)

	if err == nil {
		enviroment.Env.RDB.Set(u.Token, entry, duration)
	}

	return err
}

func Get(token string) (*model.UserAuth, error) {

	val, err := enviroment.Env.RDB.Get(token).Result()
	if err != nil || val == "" {
		return nil, ErrUserNotFound
	} else {
		userAuth := &model.UserAuth{}
		err := json.Unmarshal([]byte(val), &userAuth)
		if err != nil {
			return nil, err
		}
		return userAuth, nil
	}
}
