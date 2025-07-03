package models

import (
	"golang.org/x/crypto/bcrypt"
)

// SetPassword хеширует и сохраняет пароль пользователя
func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPasswordStr := string(hashedPassword)
	u.Password = &hashedPasswordStr
	return nil
}

// CheckPassword проверяет соответствие пароля хешу
func (u *User) CheckPassword(password string) bool {
	if u.Password == nil {
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(*u.Password), []byte(password))
	return err == nil
}

// HasPassword проверяет установлен ли пароль у пользователя
func (u *User) HasPassword() bool {
	return u.Password != nil && *u.Password != ""
}
