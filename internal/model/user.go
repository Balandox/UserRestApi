package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
	"my-rest-api/internal/enum"
)

type User struct {
	Id       int       `json:"id"`
	UserName string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     enum.Role `json:"role" `
}

func (u *User) Validate() error {

	if !u.Role.IsValid() {
		return enum.RoleIsNotValid
	}

	return validation.ValidateStruct(u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(5, 100)),
		validation.Field(&u.UserName, validation.Required, validation.Length(3, 50)),
	)
}

func (u *User) BeforeCreate() error {

	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		u.Password = enc
	}

	return nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
