package main

import (
	"database/sql"
	"fmt"
)

type CommentRepo struct {
	DB *sql.DB
}

func NewCommentRepo(db *sql.DB) *CommentRepo {
	return &CommentRepo{
		DB: db,
	}
}

func (repo *CommentRepo) Add(comment *Comment) (*string, error) {
	fmt.Println("Comment repo: add comment")
	result, err := repo.DB.Exec(`INSERT INTO comment
	(id, post_id, user_id, body, created) 
	VALUES (?, ?, ?, ?, ?)`,
		comment.ID, comment.PostId, comment.UserId, comment.Body, comment.Created)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected != 1 {
		return nil, fmt.Errorf("wrong affected rows: %d for comment id: %s", affected, comment.ID)
	}
	return &comment.ID, nil
}

func (repo *CommentRepo) Delete(id string) (bool, error) {
	fmt.Println("Comment repo: delete comment")
	result, err := repo.DB.Exec(`DELETE FROM comment WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected != 1 {
		return false, fmt.Errorf("wrong affected rows: %d for comment id %s", affected, id)
	}
	return true, nil
}
