package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock"
)

func TestPostsGetAllEmpty(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{"id", "title", "type", "description", "score", "user_id", "category_id", "created"})

	mock.
		ExpectQuery("SELECT id, title, type, description, score, user_id, category_id, created").
		WillReturnRows(rows)

	postsRepo := NewPostsRepo(db)
	posts, err := postsRepo.GetAll()

	if len(posts) > 0 {
		t.Errorf("posts must be empty. get: %s", posts)
	}

	if nil != err {
		t.Errorf("unexpected error: %s", err)
	}
}

func TestPostsGetAllIsNotEmpty(t *testing.T) {

	postsRepo := NewPostsRepo()
	posts, err := postsRepo.GetAll()

	if len(posts) == 0 || posts == nil || err != nil {
		t.Fatalf("Posts must not be empty. Get value %#v and err %#v", posts, err)
	}
}

func TestPostsGetAllError(t *testing.T) {

	postsRepo := NewPostsRepo()
	posts, err := postsRepo.GetAll()

	if err != nil {
		t.Fatalf("Must get error. Get value %#v and err %#v", posts, err)
	}
}

func TestPostsGetByIdError(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	id := "123"

	mock.
		ExpectQuery("SELECT id, title, type, description, score, user_id, category_id, created WHERE").
		WithArgs(id).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = postsRepo.GetById(&id)

	if err = mock.ExpectationsWereMet(); nil != err {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if nil == err {
		t.Errorf("expected error, got nil")
	}
}

func TestPostsGetByIdIsNotEmpty(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	id := "123"

	rows := sqlmock.
		NewRows([]string{"id", "title", "type", "description", "score", "user_id", "category_id", "created"})

	expect := []*Post{
		{"123", "test", "music", "description", 1, "123", 2, "2022-10-19 23:00:00"},
	}

	for _, post := range expect {
		rows = rows.AddRow(post.ID, post.Title, post.Type, post.Description, post.Score, post.UserID, post.CategoryID, post.Created)
	}

	mock.
		ExpectQuery("SELECT id, title, type, description, score, user_id, category_id, created").
		WithArgs(id).
		WillReturnRows(rows)

	post, err := postsRepo.GetById(&id)

	if nil != err {
		t.Errorf("unexpected error: %s", err)
	}

	if err = mock.ExpectationsWereMet(); nil != err {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if !reflect.DeepEqual(post, expect[0]) {
		t.Errorf("results not match. want %#v; have: %#v", expect[0], post)
		return
	}
}

func TestPostsGetByIdEmpty(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	id := "123"
	post, err := postsRepo.GetById(&id)

	if post != nil && err == nil {
		t.Fatalf("Posts must be empty. Get value %#v and err %#v", post, err)
	}
}

func TestPostsAddSuccessed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	post := &Post{
		Description: "test",
		Created:     "2022-10-18",
		UserID:      "1",
	}
	postId, err := postsRepo.Add(post)

	if postId == nil || err != nil {
		t.Fatalf("Must be to add successed. Get value %#v and err %#v", postId, err)
	}
}

func TestPostsAddError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	post := &Post{
		ID:          "123",
		Title:       "test",
		Type:        "text",
		Description: "test",
		Score:       1,
		UserID:      "123",
		CategoryID:  1,
		Created:     "2022-10-19",
	}
	postId, err := postsRepo.Add(post)

	if postId != nil || err == nil {
		t.Fatalf("Must be to get error. Get value %#v and err %#v", postId, err)
	}
}

func TestPostsUpdateSuccessed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	post := &Post{
		ID:          "123",
		Title:       "test",
		Type:        "text",
		Description: "test",
		Score:       1,
		UserID:      "123",
		CategoryID:  1,
		Created:     "2022-10-19",
	}
	isUpdated, err := postsRepo.Update(post)

	if isUpdated == false || err != nil {
		t.Fatalf("Must be to update successed. Get value %#v and err %#v", isUpdated, err)
	}
}

func TestPostsUpdateError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	post := &Post{
		ID:          "123",
		Title:       "test",
		Type:        "text",
		Description: "test",
		Score:       1,
		UserID:      "123",
		CategoryID:  1,
		Created:     "2022-10-19",
	}
	isUpdated, err := postsRepo.Update(post)

	if true == isUpdated && nil != err {
		t.Fatalf("Must be to update with error. Get value %#v and err %#v", isUpdated, err)
	}
}

func TestPostsDeleteSuccessed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	id := "123"
	isDeleted, err := postsRepo.Delete(id)

	if false == isDeleted && nil != err {
		t.Fatalf("Must be to delete successed. Get value %#v and err %#v", isDeleted, err)
	}
}

func TestPostsDeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	id := "123"
	isDeleted, err := postsRepo.Delete(id)

	if true == isDeleted && nil == err {
		t.Fatalf("Nust be to delete with error. Get value %#v and err %#v", isDeleted, err)
	}
}
