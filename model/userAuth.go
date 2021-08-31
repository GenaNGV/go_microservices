package model

import (
	"time"
)

type UserAuth struct {
	UserDetail User
	Token      string    `json:"token"`
	Expired    time.Time `json:"expired"`
}

//func (u *UserAuth) Put() {
//	enviroment.RDB.Set(u.Token, u, time.Duration(u.Expired.Unix()))
//}
//
//func (u *UserAuth) Get(token string)  {
//	*u = enviroment.RDB.Get(token)
//}
