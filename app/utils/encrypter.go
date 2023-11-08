package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func EncryptPassword(password string) string {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return ""
	}
	return string(res)
}

func DecryptPassword(password string, passwordToCompare string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(passwordToCompare))
	return err == nil
}
