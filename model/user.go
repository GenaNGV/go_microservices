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

	Roles []*Role `gorm:"many2many:user_role"`
}

func (user *User) TableName() string {
	return "user"
}

type Role struct {
	Id   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

func (role *Role) TableName() string {
	return "role"
}

type UserRole struct {
	UserId uint `gorm:"primaryKey;autoIncrement:false"`
	RoleId uint `gorm:"primaryKey;autoIncrement:false"`
}

func (userRole *UserRole) TableName() string {
	return "user_role"
}
