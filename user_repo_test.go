package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when open stub connetcion", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	userExpected := &User{
		ID:       "522cd619-841f-43d5-866d-f880e5f48d18",
		Login:    "mer",
		Password: "test",
	}

	// success
	rows := sqlmock.NewRows([]string{
		"id", "login", "password",
	}).AddRow(userExpected.ID, userExpected.Login, userExpected.Password)

	mock.ExpectQuery(`SELECT id, login, password FROM user WHERE id = `).
		WithArgs(userExpected.ID).
		WillReturnRows(rows)
	user, err := userRepo.GetById(userExpected.ID)
	if err != nil {
		t.Errorf("not expected error %s", err)
		return
	}
	if !reflect.DeepEqual(user, userExpected) {
		t.Errorf("results are not matched; want: %#v, have: %#v", userExpected, user)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there are unfulfilled expectations: %s", err)
		return
	}

	//query error
	mock.ExpectQuery(`SELECT id, login, password FROM user WHERE id = `).
		WithArgs(userExpected.ID).
		WillReturnError(fmt.Errorf("db error"))
	_, err = userRepo.GetById(userExpected.ID)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there are unfulfilled expectations: %s", err)
		return
	}

	//scan error
	rows = sqlmock.NewRows([]string{
		"id", "login",
	}).AddRow(userExpected.ID, userExpected.Login)
	mock.ExpectQuery(`SELECT id, login, password FROM user WHERE id = `).
		WithArgs(userExpected.ID).
		WillReturnRows(rows)

	_, err = userRepo.GetById(userExpected.ID)
	if err == nil {
		t.Errorf("scan error expected, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserGetByLogin(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when open stub connetcion", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	userExpected := &User{
		ID:       "522cd619-841f-43d5-866d-f880e5f48d18",
		Login:    "mer",
		Password: "test",
	}

	// success
	rows := sqlmock.NewRows([]string{
		"id", "login", "password",
	}).AddRow(userExpected.ID, userExpected.Login, userExpected.Password)

	mock.ExpectQuery(`SELECT id, login, password FROM user WHERE login = `).
		WithArgs(userExpected.Login).
		WillReturnRows(rows)
	user, err := userRepo.GetByLogin(userExpected.Login)
	if err != nil {
		t.Errorf("not expected error %s", err)
		return
	}
	if !reflect.DeepEqual(user, userExpected) {
		t.Errorf("results are not matched; want: %#v, have: %#v", userExpected, user)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there are unfulfilled expectations: %s", err)
		return
	}

	//query error
	mock.ExpectQuery(`SELECT id, login, password FROM user WHERE login = `).
		WithArgs(userExpected.Login).
		WillReturnError(fmt.Errorf("db error"))
	_, err = userRepo.GetByLogin(userExpected.Login)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there are unfulfilled expectations: %s", err)
		return
	}

	//scan error
	rows = sqlmock.NewRows([]string{
		"id", "login",
	}).AddRow(userExpected.ID, userExpected.Login)
	mock.ExpectQuery(`SELECT id, login, password FROM user WHERE login = `).
		WithArgs(userExpected.Login).
		WillReturnRows(rows)

	_, err = userRepo.GetByLogin(userExpected.Login)
	if err == nil {
		t.Errorf("scan error expected, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
}

func TestUserCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when open stub connetcion", err)
	}
	defer db.Close()

	userRepo := NewUserRepo(db)

	userExpected := &User{
		ID:       "522cd619-841f-43d5-866d-f880e5f48d18",
		Login:    "mer",
		Password: "test",
		Created:  "2022-11-09T19:51:42Z",
	}

	// success
	mock.ExpectExec(`INSERT INTO user`).
		WithArgs(userExpected.ID, userExpected.Login, userExpected.Password, userExpected.Created).
		WillReturnResult(sqlmock.NewResult(0, 1))
	lastID, err := userRepo.Create(userExpected)
	if err != nil {
		t.Errorf("not expected error %s", err)
		return
	}
	if userExpected.ID != *lastID {
		t.Errorf("results are not matched; want: %#v, have: %#v", userExpected.ID, lastID)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there are unfulfilled expectations: %s", err)
		return
	}

	// query error
	mock.ExpectExec(`INSERT INTO user`).
		WithArgs(userExpected.ID, userExpected.Login, userExpected.Password, userExpected.Created).
		WillReturnError(fmt.Errorf("db error"))
	_, err = userRepo.Create(userExpected)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there are unfulfilled expectations: %s", err)
		return
	}
}
