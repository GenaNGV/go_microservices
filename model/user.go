package model

import (
	"time"
)

type User struct {
	id         uint   `gorm:"primaryKey"`
	Email      string `json:"email"`
	Password   string
	created    time.Time `json:"created"`
	deleted    time.Time `json:"deleted"`
	deleted_by uint      `json:"deletedBy"`
	last_login time.Time `json:"lastLogin"`
	first_name string    `json:"firstName"`
	last_name  string    `json:"lastName"`
}
