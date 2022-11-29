package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

var (
	ErrNoAuth = errors.New("No session found")
)

type SessionsDBManagerJWT struct {
	DB *sql.DB
}

type UserJWtClaims struct {
	UserName string `json:"username"`
	ID       string `json:"id"`
	SessID   string `json:"sess_id"`
}

type SessionJWTClaims struct {
	User UserJWtClaims `json:"user"`
	jwt.StandardClaims
}

func NewSessionDBManagerJWT(db *sql.DB) *SessionsDBManagerJWT {
	return &SessionsDBManagerJWT{
		DB: db,
	}
}

func (sm *SessionsDBManagerJWT) Check(r *http.Request) (*Session, error) {

	var err error
	authHeader := r.Header.Get("Authorization")
	_, tokenString, _ := strings.Cut(authHeader, "Bearer ")
	if tokenString == "" {
		err = fmt.Errorf("no token found: %s", authHeader)
		fmt.Println(err)
		return nil, err
	}

	var secretKey = []byte(os.Getenv("SECRET_KEY"))

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return secretKey, nil
	}
	payload := &SessionJWTClaims{}
	_, err = jwt.ParseWithClaims(tokenString, payload, hashSecretGetter)

	if nil != err || payload.Valid() != nil {
		fmt.Println(authHeader)
		fmt.Println(tokenString)
		err = fmt.Errorf("bad token: %s %s", err.Error(), tokenString)
		fmt.Println(err)
		return nil, err
	}

	sess := &Session{}
	fmt.Printf("check session %#v\n", payload)
	row := sm.DB.QueryRow("SELECT id, user_id FROM sessions WHERE id = ?", payload.User.SessID)

	err = row.Scan(&sess.ID, &sess.UserID)

	fmt.Printf("check session result %#v\n", sess)

	if err == sql.ErrNoRows {
		fmt.Println("Check session no rows")
		return nil, ErrNoAuth
	} else if err != nil {
		fmt.Println("Check session err: ", err)
		return nil, ErrNoAuth
	}

	return sess, nil
}

func (sm *SessionsDBManagerJWT) Create(w http.ResponseWriter, user *User) (*Session, error) {
	sessID := RandStringRunes(32)
	_, err := sm.DB.Exec("INSERT INTO sessions (user_id, id) VALUES(?, ?)", user.ID, sessID)
	if err != nil {
		return nil, err
	}

	return &Session{
		UserID: user.ID,
		ID:     sessID,
	}, nil
}

func (sm *SessionsDBManagerJWT) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
	sess, err := SessionFromContext(r.Context())
	if err == nil {
		_, err = sm.DB.Exec("DELETE FROM sessions WHERE id = ?", sess.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sm *SessionsDBManagerJWT) DestroyAll(w http.ResponseWriter, user *User) error {
	result, err := sm.DB.Exec("DELETE FROM sessions WHERE user_id = ?", user.ID)
	if err != nil {
		return err
	}

	affected, _ := result.RowsAffected()
	log.Println("destroyed sessions", affected, "for user", user.ID)

	return nil
}
