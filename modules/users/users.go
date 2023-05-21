package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Username string `db:"username" json:"username"`
	RoleId   int `db:"role_id" json:"role_id"`
}

type UserRegisterReq struct {
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"-" db:"password"`
}

func (obj *UserRegisterReq) BcryptHashing() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(obj.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("bcrypt password failed: %v", err)
	}
	obj.Password = string(hash)
	return nil
}

func (obj *UserRegisterReq) BcryptPCompare(hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(obj.Password))
}

func (obj *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
}

type UserToken struct{
	AccessToken string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
	Id string `json:"id" db:"id"`
}


type UserPassport struct {
	User *User `json:"user"`
	Token *UserToken `json:"token"`
}