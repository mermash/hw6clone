package main

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type UserUtils struct{}

func (u *UserUtils) GenerateJWT(user *User, sessID string) (string, error) {
	var signingKey = []byte(os.Getenv("SECRET_KEY"))
	data := &SessionJWTClaims{
		User: UserJWtClaims{
			UserName: user.Login,
			ID:       user.ID,
			SessID:   sessID,
		},
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(90 * 24 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, data).SignedString(signingKey)

	if nil != err {
		fmt.Printf("Error during generate token: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}

func (u *UserUtils) GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (u *UserUtils) CheckPasswordHash(passwordReceived string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordReceived))
	isEqual := nil == err
	if !isEqual {
		fmt.Printf("check password hash: %s", err.Error())
		return false
	}
	return isEqual
}
