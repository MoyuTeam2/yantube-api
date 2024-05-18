package models

import "api/pkg/utils"

type User struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	Username   string `json:"username" gorm:"unique"`
	Password   string `json:"password"`
	StreamCode string `json:"stream_code"`
}

func (User) TableName() string {
	return "user"
}

func (u *User) ValidatePassword(password string) (bool, error) {
	return utils.ValidatePbkdf2Cipher(u.Password, password)
}

func NewUserWithPassword(username, password string) (*User, error) {
	salt, err := utils.RandomBytesWithDefaultLen()
	if err != nil {
		return nil, err
	}
	return &User{
		Username:   username,
		Password:   utils.DefaultPbkdf2.FormatedString([]byte(password), salt),
		StreamCode: utils.GenerateRandomString(16),
	}, nil
}
