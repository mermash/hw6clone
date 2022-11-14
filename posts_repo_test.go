package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestPostsGetAllSuccessed(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{
			"post_id", "title", "type",
			"description", "score", "user_id",
			"category_id", "post_created",
			"user_user_id", "login",
			"category_name",
		})

	expect := []*PostComplexData{
		{
			Post: Post{
				ID:          "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
				Title:       "test fashion",
				Type:        "text",
				Description: "test fashion",
				Score:       1,
				UserID:      "522cd619-841f-43d5-866d-f880e5f48d18",
				CategoryID:  1,
				Created:     "2022-11-09T19:51:42Z",
			},
			User: User{
				ID:    "522cd619-841f-43d5-866d-f880e5f48d18",
				Login: "mer",
			},
			Category: Category{
				Name: "fashion",
			}},
	}

	for _, post := range expect {
		rows = rows.AddRow(post.Post.ID, post.Post.Title,
			post.Post.Type, post.Post.Description, post.Post.Score,
			post.Post.UserID, post.Post.CategoryID, post.Post.Created, post.User.ID, post.User.Login,
			post.Category.Name)
	}

	mock.
		ExpectQuery(`
		SELECT 
		post.id AS post_id, title, type, description, score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		ORDER BY post.created DESC`).
		WillReturnRows(rows)

	postsRepo := NewPostsRepo(db)
	posts, err := postsRepo.GetAll()

	if nil != err {
		t.Errorf("unexpected error: %s", err)
	}

	if !reflect.DeepEqual(posts, expect) {
		t.Errorf("results not match. want %#v; have: %#v", expect, posts)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were an unfulfilled expectations: %s", err)
	}
}

func TestPostsGetAllQueryError(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when open stub connetcion", err)
	}
	defer db.Close()

	mock.ExpectQuery(
		`SELECT
		post.id AS post_id, title, type, description, score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		ORDER BY post.created DESC`).
		WillReturnError(fmt.Errorf("db_error"))

	postsRepo := NewPostsRepo(db)
	_, err = postsRepo.GetAll()

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were an unfulfilled expectations")
	}
}

func TestPostsGetAllScanError(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error %s was not expected when open stub connetcion", err)
	}
	defer db.Close()

	rows := sqlmock.
		NewRows([]string{
			"post_id", "title", "type",
		}).AddRow("dc1e2f25-76a5-4aac-9212-96e2121c16f1", "test fashion", "text")

	mock.ExpectQuery(
		`SELECT
		post.id AS post_id, title, type, description, score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		ORDER BY post.created DESC`).
		WillReturnRows(rows)

	postsRepo := NewPostsRepo(db)
	_, err = postsRepo.GetAll()

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were an unfulfilled expectations: %s", err)
	}
}

func TestPostsGetByIdQueryError(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	id := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectQuery(
			`
		SELECT 
	post.id AS post_id, title, type, description, 
	score, user_id, category_id, post.created AS post_created,
	user.id AS user_user_id, user.login,
	category.name AS category_name
	FROM post 
	LEFT JOIN user ON user.id = post.user_id
	LEFT JOIN category ON category.id = post.category_id WHERE post.id =`).
		WithArgs(id).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = postsRepo.GetById(id)

	if err := mock.ExpectationsWereMet(); nil != err {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if nil == err {
		t.Errorf("expected error, got nil")
	}
}

func TestPostsGetByIdSuccessed(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	id := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	rows := sqlmock.
		NewRows([]string{
			"post_id", "title", "type",
			"description", "score", "user_id",
			"category_id", "post_created",
			"user_user_id", "login",
			"category_name"})

	expect := []*PostComplexData{
		{
			Post: Post{
				ID:          "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
				Title:       "test fashion",
				Type:        "text",
				Description: "test fashion",
				Score:       1,
				UserID:      "522cd619-841f-43d5-866d-f880e5f48d18",
				CategoryID:  1,
				Created:     "2022-11-09T19:51:42Z",
			},
			User: User{
				ID:    "522cd619-841f-43d5-866d-f880e5f48d18",
				Login: "mer",
			},
			Category: Category{
				Name: "fashion",
			}},
	}

	for _, post := range expect {
		rows = rows.AddRow(post.Post.ID, post.Post.Title,
			post.Post.Type, post.Post.Description, post.Post.Score,
			post.Post.UserID, post.Post.CategoryID, post.Post.Created, post.User.ID, post.User.Login,
			post.Category.Name)
	}

	mock.
		ExpectQuery(`
		SELECT 
	post.id AS post_id, title, type, description, 
	score, user_id, category_id, post.created AS post_created,
	user.id AS user_user_id, user.login,
	category.name AS category_name
	FROM post 
	LEFT JOIN user ON user.id = post.user_id
	LEFT JOIN category ON category.id = post.category_id WHERE post.id =`).
		WithArgs(id).
		WillReturnRows(rows)

	post, err := postsRepo.GetById(id)

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

func TestPostsGetByIdScanError(t *testing.T) {

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	id := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	rows := sqlmock.
		NewRows([]string{"post_id", "title", "type"}).
		AddRow("dc1e2f25-76a5-4aac-9212-96e2121c16f1", "test fashion", "text")

	mock.
		ExpectQuery(`
		SELECT 
	post.id AS post_id, title, type, description, 
	score, user_id, category_id, post.created AS post_created,
	user.id AS user_user_id, user.login,
	category.name AS category_name
	FROM post 
	LEFT JOIN user ON user.id = post.user_id
	LEFT JOIN category ON category.id = post.category_id WHERE post.id =`).
		WithArgs(id).
		WillReturnRows(rows)

	_, err = postsRepo.GetById(id)

	if err := mock.ExpectationsWereMet(); nil != err {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}

	if err == nil {
		t.Errorf("expected error, got nil")
		return
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
		ID:          "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
		Title:       "test",
		Type:        "text",
		Description: "test",
		Score:       1,
		Created:     "2022-11-09T19:51:42Z",
		UserID:      "522cd619-841f-43d5-866d-f880e5f48d18",
		CategoryID:  1,
	}

	mock.
		ExpectExec(`INSERT INTO post`).
		WithArgs(post.ID, post.Title, post.Type, post.Description, post.Score, post.UserID, post.CategoryID, post.Created).
		WillReturnResult(sqlmock.NewResult(0, 1))

	postId, err := postsRepo.Add(post)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if postId != &post.ID {
		t.Errorf("expected lastId: %s", post.ID)
	}
}

func TestPostsAddQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	post := &Post{
		ID:          "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
		Title:       "test",
		Type:        "text",
		Description: "test",
		Score:       1,
		Created:     "2022-11-09T19:51:42Z",
		UserID:      "522cd619-841f-43d5-866d-f880e5f48d18",
		CategoryID:  1,
	}

	mock.
		ExpectExec(`INSERT INTO post`).
		WithArgs(post.ID, post.Title, post.Type, post.Description, post.Score, post.UserID, post.CategoryID, post.Created).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = postsRepo.Add(post)

	if err == nil {
		t.Errorf("expected error got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsAddResultError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	post := &Post{
		ID:          "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
		Title:       "test",
		Type:        "text",
		Description: "test",
		Score:       1,
		Created:     "2022-11-09T19:51:42Z",
		UserID:      "522cd619-841f-43d5-866d-f880e5f48d18",
		CategoryID:  1,
	}

	mock.
		ExpectExec(`INSERT INTO post`).
		WithArgs(post.ID, post.Title, post.Type, post.Description, post.Score, post.UserID, post.CategoryID, post.Created).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = postsRepo.Add(post)

	if err == nil {
		t.Errorf("expected error got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsDeleteSuccessed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec(`DELETE FROM post`).
		WithArgs(postId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := postsRepo.Delete(postId)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	if !result {
		t.Errorf("expected true")
	}
}

func TestPostsDeleteQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	//query error
	mock.
		ExpectExec(`DELETE FROM post`).
		WithArgs(postId).
		WillReturnError(fmt.Errorf("bad_query"))

	_, err = postsRepo.Delete(postId)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsDeleteResultError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec(`DELETE FROM post`).
		WithArgs(postId).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))

	_, err = postsRepo.Delete(postId)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsDeleteRowsAffectedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec(`DELETE FROM post`).
		WithArgs(postId).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err = postsRepo.Delete(postId)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsUpVoteSuccessed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connetcion", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec("UPDATE post SET").
		WithArgs(postId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := postsRepo.UpVote(postId)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if !result {
		t.Errorf("expected true")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsUpVoteQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec(`UPDATE post SET`).
		WithArgs(postId).
		WillReturnError(fmt.Errorf("bad_query"))

	_, err = postsRepo.UpVote(postId)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsUpVoteResultError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec(`UPDATE post SET`).
		WithArgs(postId).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad result")))

	_, err = postsRepo.UpVote(postId)

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsUpVoteRowsAffectedError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec(`UPDATE post SET`).
		WithArgs(postId).
		WillReturnResult(sqlmock.NewResult(0, 0))

	_, err = postsRepo.UpVote(postId)

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsDownVoteSuccessed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec(`UPDATE post SET`).
		WithArgs(postId).
		WillReturnResult(sqlmock.NewResult(0, 1))

	result, err := postsRepo.DownVote(postId)

	if !result {
		t.Errorf("expected true")
		return
	}

	if err != nil {
		t.Errorf("expected result, got error %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsDownVoteQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec(`UPDATE post SET`).
		WithArgs(postId).
		WillReturnError(fmt.Errorf("bad_query"))

	_, err = postsRepo.DownVote(postId)

	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsDownVoteResultError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"

	mock.
		ExpectExec(`UPDATE post SET`).
		WithArgs(postId).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad result")))

	_, err = postsRepo.DownVote(postId)

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestPostsGetByCategoryNameSuccessed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	categoryName := "fashion"

	rows := sqlmock.
		NewRows([]string{
			"post_id", "title", "type",
			"description", "score", "user_id",
			"category_id", "post_created",
			"user_user_id", "login",
			"category_name"})

	expect := []*PostComplexData{
		{
			Post: Post{
				ID:          "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
				Title:       "test fashion",
				Type:        "text",
				Description: "test fashion",
				Score:       1,
				UserID:      "522cd619-841f-43d5-866d-f880e5f48d18",
				CategoryID:  1,
				Created:     "2022-11-09T19:51:42Z",
			},
			User: User{
				ID:    "522cd619-841f-43d5-866d-f880e5f48d18",
				Login: "mer",
			},
			Category: Category{
				Name: "fashion",
			}},
	}

	for _, post := range expect {
		rows = rows.AddRow(post.Post.ID, post.Post.Title,
			post.Post.Type, post.Post.Description, post.Post.Score,
			post.Post.UserID, post.Post.CategoryID, post.Post.Created, post.User.ID, post.User.Login,
			post.Category.Name)
	}

	mock.
		ExpectQuery(`		SELECT
		post.id AS post_id, title, type, description, 
		score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post 
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		WHERE category.name = `).
		WithArgs(categoryName).
		WillReturnRows(rows)

	posts, err := postsRepo.GetByCategoryName(categoryName)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if !reflect.DeepEqual(posts, expect) {
		t.Errorf("results are not matched; want: %#v; have: %#v", expect, posts)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were an unfulfilled expectations %s", err)
	}
}

func TestPostsGetByCategoryNameScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	categoryName := "fashion"

	rows := sqlmock.
		NewRows([]string{
			"post_id", "title", "type"}).
		AddRow("dc1e2f25-76a5-4aac-9212-96e2121c16f1", "test fashion", "text")

	mock.
		ExpectQuery(`		SELECT
		post.id AS post_id, title, type, description, 
		score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post 
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		WHERE category.name = `).
		WithArgs(categoryName).
		WillReturnRows(rows)

	_, err = postsRepo.GetByCategoryName(categoryName)

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were an unfulfilled expectations %s", err)
	}
}

func TestPostsGetByCategoryNameQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	categoryName := "fashion"

	mock.
		ExpectQuery(`		SELECT
		post.id AS post_id, title, type, description, 
		score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post 
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		WHERE category.name = `).
		WithArgs(categoryName).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = postsRepo.GetByCategoryName(categoryName)

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were an unfulfilled expectations %s", err)
	}
}

func TestPostsGetByLoginSuccessed(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	login := "test"

	rows := sqlmock.
		NewRows([]string{
			"post_id", "title", "type",
			"description", "score", "user_id",
			"category_id", "post_created",
			"user_user_id", "login",
			"category_name"})

	expect := []*PostComplexData{
		{
			Post: Post{
				ID:          "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
				Title:       "test fashion",
				Type:        "text",
				Description: "test fashion",
				Score:       1,
				UserID:      "522cd619-841f-43d5-866d-f880e5f48d18",
				CategoryID:  1,
				Created:     "2022-11-09T19:51:42Z",
			},
			User: User{
				ID:    "522cd619-841f-43d5-866d-f880e5f48d18",
				Login: "mer",
			},
			Category: Category{
				Name: "fashion",
			}},
	}

	for _, post := range expect {
		rows = rows.AddRow(post.Post.ID, post.Post.Title,
			post.Post.Type, post.Post.Description, post.Post.Score,
			post.Post.UserID, post.Post.CategoryID, post.Post.Created, post.User.ID, post.User.Login,
			post.Category.Name)
	}

	mock.
		ExpectQuery(`	SELECT 
		post.id AS post_id, title, type, description, 
		score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post 
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		WHERE user.login = ?`).
		WithArgs(login).
		WillReturnRows(rows)

	posts, err := postsRepo.GetByUserLogin(login)

	if err != nil {
		t.Errorf("unexpected error: %s", err)
		return
	}

	if !reflect.DeepEqual(posts, expect) {
		t.Errorf("result is not matched; want: %#v; have: %#v", posts, expect)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("there were unfulfilled expectations")
		return
	}
}

func TestPostsGetByLoginScanError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	login := "test"

	rows := sqlmock.
		NewRows([]string{
			"post_id", "title", "type"}).
		AddRow("dc1e2f25-76a5-4aac-9212-96e2121c16f1", "test fashion", "text")

	mock.
		ExpectQuery(`	SELECT 
		post.id AS post_id, title, type, description, 
		score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post 
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		WHERE user.login = ?`).
		WithArgs(login).
		WillReturnRows(rows)

	_, err = postsRepo.GetByUserLogin(login)

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("there were unfulfilled expectations")
		return
	}
}

func TestPostsGetByLoginQueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	postsRepo := NewPostsRepo(db)
	login := "test"

	mock.
		ExpectQuery(`	SELECT 
		post.id AS post_id, title, type, description, 
		score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post 
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		WHERE user.login = ?`).
		WithArgs(login).
		WillReturnError(fmt.Errorf("db_error"))

	_, err = postsRepo.GetByUserLogin(login)

	if err == nil {
		t.Error("expected error, got nil")
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("there were unfulfilled expectations")
		return
	}
}
