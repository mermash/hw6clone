// package main

// import (
// 	"database/sql"
// 	"errors"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"time"
// )

// var (
// 	ErrNoAuth   = errors.New("No session found")
// 	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
// )

// type SessionsDBManager struct {
// 	DB *sql.DB
// }

// func RandStringRunes(n int) string {
// 	b := make([]rune, n)
// 	for i := range b {
// 		b[i] = letterRunes[rand.Intn(len(letterRunes))]
// 	}
// 	return string(b)
// }

// func NewSessionDBManager(db *sql.DB) *SessionsDBManager {
// 	return &SessionsDBManager{
// 		DB: db,
// 	}
// }

// func (sm *SessionsDBManager) Check(r *http.Request) (*Session, error) {
// 	fmt.Println("Check session db")
// 	sessionCookie, err := r.Cookie("session_id")

// 	if err == http.ErrNoCookie {
// 		fmt.Println("check session no rows")
// 		return nil, ErrNoAuth
// 	}

// 	fmt.Println("session: ", sessionCookie.Value)
// 	sess := &Session{}
// 	row := sm.DB.QueryRow("SELECT user_id FROM sessions WHERE id = ?", sessionCookie.Value)

// 	err = row.Scan(sess.UserID)

// 	if err == sql.ErrNoRows {
// 		fmt.Println("Check session no rows")
// 		return nil, ErrNoAuth
// 	} else if err != nil {
// 		fmt.Println("Check session err: ", err)
// 		return nil, ErrNoAuth
// 	}

// 	sess.ID = sessionCookie.Value
// 	return sess, nil
// }

// func (sm *SessionsDBManager) Create(w http.ResponseWriter, user *User) (*string, error) {
// 	sessID := RandStringRunes(32)
// 	_, err := sm.DB.Exec("INSERT INTO sessions (user_id, id) VALUES(?, ?)", user.ID, sessID)
// 	if err != nil {
// 		return nil, err
// 	}

// 	cookie := &http.Cookie{
// 		Name:     "session_id",
// 		Value:    sessID,
// 		Expires:  time.Now().Add(90 * 24 * time.Hour),
// 		Path:     "/",
// 		HttpOnly: true,
// 		SameSite: http.SameSiteStrictMode,
// 		Secure:   true,
// 	}
// 	http.SetCookie(w, cookie)
// 	return &sessID, nil
// }

// func (sm *SessionsDBManager) DestroyCurrent(w http.ResponseWriter, r *http.Request) error {
// 	sess, err := SessionFromContext(r.Context())
// 	if err == nil {
// 		_, err = sm.DB.Exec("DELETE FROM sessions WHERE id = ?", sess.ID)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	cookie := &http.Cookie{
// 		Name:     "session_id",
// 		Expires:  time.Now().AddDate(0, 0, -1),
// 		Path:     "/",
// 		HttpOnly: true,
// 		SameSite: http.SameSiteStrictMode,
// 		Secure:   true,
// 	}
// 	http.SetCookie(w, cookie)
// 	return nil
// }

// func (sm *SessionsDBManager) DestroyAll(w http.ResponseWriter, user *User) error {
// 	result, err := sm.DB.Exec("DELETE FROM sessions WHERE user_id = ?", user.ID)
// 	if err != nil {
// 		return err
// 	}

// 	affected, _ := result.RowsAffected()
// 	log.Println("destroyed sessions", affected, "for user", user.ID)

// 	return nil
// }
