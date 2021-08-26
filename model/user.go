package model

import (
	"time"
)

type User struct {
	Id        uint       `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Created   time.Time  `json:"created"`
	Deleted   *time.Time `json:"deleted,omitempty"`
	DeletedBy *uint      `json:"deletedBy,omitempty"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`
}

type Tabler interface {
	TableName() string
}

func (User) TableName() string {
	return "user"
}
