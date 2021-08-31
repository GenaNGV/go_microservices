package model

import (
	"time"
)

type UserAuth struct {
	User
	Token   string    `json:"token"`
	Expired time.Time `json:"expired"`
}
