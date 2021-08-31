package model

import (
	"time"
)

type User struct {
	Id        uint       `json:"id" gorm:"primaryKey"`
	Email     string     `json:"email"`
	Password  string     `json:"-"`
	Created   time.Time  `json:"created"`
	Deleted   *time.Time `json:"deleted,omitempty" gorm:"index"`
	DeletedBy *uint      `json:"deletedBy,omitempty"`
	LastLogin *time.Time `json:"lastLogin,omitempty"`
	FirstName string     `json:"firstName"`
	LastName  string     `json:"lastName"`

	Roles []*Role `json:"roles" gorm:"many2many:user_role"`
}

type Role struct {
	Id   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type UserRole struct {
	UserId uint `gorm:"primaryKey;autoIncrement:false"`
	RoleId uint `gorm:"primaryKey;autoIncrement:false"`
}

func (user *User) TableName() string {
	return "user"
}

func (role *Role) TableName() string {
	return "role"
}

func (userRole *UserRole) TableName() string {
	return "user_role"
}
