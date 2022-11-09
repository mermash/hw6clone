package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserRepo struct {
	DB *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		DB: db,
	}
}

func (repo *UserRepo) GetById(id string) (*User, error) {
	fmt.Println("Get user by id")
	user := &User{}
	err := repo.DB.
		QueryRow("SELECT id, login, password FROM user WHERE id = ?", id).
		Scan(&user.ID, &user.Login, &user.Password)
	if nil != err {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepo) GetByLogin(login string) (*User, error) {
	fmt.Println("Get user by login")
	user := &User{}
	err := repo.DB.
		QueryRow("SELECT id, login, password FROM user WHERE login = ?", login).
		Scan(&user.ID, &user.Login, &user.Password)
	if nil != err {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepo) Create(user *User) (*string, error) {
	fmt.Println("Create new user")
	uuid := uuid.NewString()
	created := time.Now().Format(time.RFC3339)
	fmt.Println("created: ", created)
	_, err := repo.DB.Exec(
		"INSERT INTO user (id, login, password, created) VALUES(?, ?, ?, ?)",
		uuid,
		user.Login,
		user.Password,
		created,
	)

	if nil != err {
		return nil, err
	}

	fmt.Println("new id", uuid)
	return &uuid, nil
}
