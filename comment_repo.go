package main

import (
	"database/sql"
	"fmt"
	"strings"
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

func (repo *CommentRepo) GetCommentsByPostIds(postIds []string) (map[string][]*CommentComplexData, error) {

	lenPostId := len(postIds)
	placeHolders := make([]string, 0, lenPostId)
	args := make([]interface{}, 0, lenPostId)
	for _, id := range postIds {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}
	query :=
		`SELECT 
	comment.id AS comment_id, post_id, body, 
	comment.created AS comment_created,
	user.id AS user_id, user.login
	FROM comment 
	LEFT JOIN user ON user.id = comment.user_id
	WHERE post_id IN (` + strings.Join(placeHolders, ",") + `)`
	fmt.Println("get comments postIDs", postIds)
	fmt.Println("get comments sql query: ", query)
	rows, err := repo.DB.Query(query, args...)
	if nil != err {
		fmt.Println("get comments query:", err)
		return nil, err
	}
	defer rows.Close()
	comments := map[string][]*CommentComplexData{}
	for rows.Next() {
		data := &CommentComplexData{}
		err := rows.Scan(&data.Comment.ID, &data.Comment.PostId,
			&data.Comment.Body, &data.Comment.Created, &data.User.ID, &data.User.Login)
		if nil != err {
			fmt.Println("get comments scan:", err)
			return nil, err
		}
		fmt.Println("get comment for post id", data, data.Comment.PostId)
		if _, ok := comments[data.Comment.PostId]; !ok {
			comments[data.Comment.PostId] = []*CommentComplexData{}
		}
		comments[data.Comment.PostId] = append(comments[data.Comment.PostId], data)
	}
	fmt.Println("comments: ", comments)
	return comments, nil
}
