package models

import (
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ModelFields `s2m:"-"`
	Username    string `json:"username" gorm:"type:char(32)"`
	Email       string `json:"email" gorm:"unique"`
	Mobile      string `json:"mobile" gorm:"unique"`
	Password    string `json:"password,omitempty"`
	Active      bool   `json:"active"`
}

func (u *User) Validate() error {
	if !ValidateEmail(u.Email) {
		return fmt.Errorf("invalid email")
	}
	if !ValidateMobile(u.Mobile) {
		return fmt.Errorf("invalid mobile")
	}
	if !ValidatePassword(u.Password) {
		return fmt.Errorf("password must be at least 6 characters")
	}
	return nil
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)
	return nil
}

func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		u.Password = string(hash)
	}
	return nil
}

func (u *User) CheckPassword(pwd string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd)); err != nil {
		return err
	}
	return nil
}

func (u *User) MarshalJSON() ([]byte, error) {
	type user User
	userCopy := user(*u)
	userCopy.Password = ""
	return json.Marshal(userCopy)
}
