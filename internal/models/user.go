package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	BaseFields
	Username     string     `gorm:"unique" json:"username"`
	Nickname     string     `json:"nickname"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone" `
	Password     string     `json:"-"`
	Gender       string     `json:"gender"`
	Age          int        `json:"age"`
	Status       int        `gorm:"default:1" json:"status"`
	DepartmentID int        `gorm:"column:department_id;default:0" json:"department_id"`
	Roles        []RBACRole `gorm:"many2many:rbac_user_roles;" json:"roles"`
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
