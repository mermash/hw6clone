package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(passwordReceived string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordReceived))
	isEqual := nil == err
	if !isEqual {
		fmt.Errorf("check password hash: %s", err.Error())
	}
	return isEqual
}
