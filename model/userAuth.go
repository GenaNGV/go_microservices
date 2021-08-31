package model

type UserAuth struct {
	UserDetail User
	Token      string `json:"token"`
	Expired    int64  `json:"expired"`
}
