package main

import (
	"database/sql"
	"fmt"
)

type PostsRepo struct {
	DB *sql.DB
}

func NewPostsRepo(db *sql.DB) *PostsRepo {
	postsRepo := &PostsRepo{
		DB: db,
	}
	fmt.Println("Create new postsRepo", postsRepo)
	return postsRepo
}

func (repo *PostsRepo) GetAll() ([]*PostComplexData, error) {
	fmt.Println("Repo post: get all posts")

	rows, err := repo.DB.
		Query(
			`
			SELECT 
			post.id AS post_id, title, type, description, score, user_id, category_id, post.created AS post_created,
			user.id AS user_id, user.login,
			category.name AS category_name
			FROM post
			LEFT JOIN user ON user.id = post.user_id
			LEFT JOIN category ON category.id = post.category_id
			ORDER BY post.created DESC
			`)
	if nil != err {
		fmt.Println("get all: ", err)
		return nil, err
	}
	defer rows.Close()

	posts := make([]*PostComplexData, 0, 10)
	for rows.Next() {
		data := &PostComplexData{}
		err := rows.Scan(&data.Post.ID, &data.Post.Title,
			&data.Post.Type, &data.Post.Description,
			&data.Post.Score, &data.Post.UserID,
			&data.Post.CategoryID, &data.Post.Created,
			&data.User.ID, &data.User.Login,
			&data.Category.Name)
		if nil != err {
			fmt.Println("scan: ", err)
			return nil, err
		}
		posts = append(posts, data)
	}

	return posts, nil
}

func (repo *PostsRepo) GetById(id string) (*PostComplexData, error) {
	fmt.Println("Repo post: get by id post")

	row := repo.DB.QueryRow(`
	SELECT 
	post.id AS post_id, title, type, description, 
	score, user_id, category_id, post.created AS post_created,
	user.id AS user_user_id, user.login,
	category.name AS category_name
	FROM post 
	LEFT JOIN user ON user.id = post.user_id
	LEFT JOIN category ON category.id = post.category_id
	WHERE post.id = ?`, id)

	data := &PostComplexData{}
	err := row.Scan(&data.Post.ID, &data.Post.Title, &data.Post.Type,
		&data.Post.Description, &data.Post.Score, &data.Post.UserID,
		&data.Post.CategoryID, &data.Post.Created,
		&data.User.ID, &data.User.Login,
		&data.Category.Name)
	if nil != err {
		return nil, err
	}

	return data, nil
}

func (repo *PostsRepo) GetByCategoryName(categoryName string) ([]*PostComplexData, error) {
	fmt.Println("Repo post: get posts by categoryName")
	rows, err := repo.DB.Query(
		`
		SELECT
		post.id AS post_id, title, type, description, 
		score, user_id, category_id, post.created AS post_created,
		user.id AS user_id, user.login,
		category.name AS category_name
		FROM post 
		LEFT JOIN user ON user.id = post.user_id
		LEFT JOIN category ON category.id = post.category_id
		WHERE category.name = ?`,
		categoryName)
	if nil != err {
		fmt.Println("get all: ", err)
		return nil, err
	}
	defer rows.Close()

	posts := make([]*PostComplexData, 0, 10)
	for rows.Next() {
		data := &PostComplexData{}
		err := rows.Scan(&data.Post.ID, &data.Post.Title,
			&data.Post.Type, &data.Post.Description,
			&data.Post.Score, &data.Post.UserID,
			&data.Post.CategoryID, &data.Post.Created,
			&data.User.ID, &data.User.Login,
			&data.Category.Name)
		if nil != err {
			fmt.Println("scan: ", err)
			return nil, err
		}
		posts = append(posts, data)
	}

	return posts, nil
}

func (repo *PostsRepo) GetByUserLogin(userLogin string) ([]*PostComplexData, error) {
	fmt.Println("Repo post: get posts by user login")

	rows, err := repo.DB.Query(`
	SELECT 
	post.id AS post_id, title, type, description, 
	score, user_id, category_id, post.created AS post_created,
	user.id AS user_id, user.login,
	category.name AS category_name
	FROM post 
	LEFT JOIN user ON user.id = post.user_id
	LEFT JOIN category ON category.id = post.category_id
	WHERE user.login = ?`,
		userLogin)
	if nil != err {
		fmt.Println("get all: ", err)
		return nil, err
	}
	defer rows.Close()

	posts := make([]*PostComplexData, 0, 10)
	for rows.Next() {
		data := &PostComplexData{}
		err := rows.Scan(&data.Post.ID, &data.Post.Title,
			&data.Post.Type, &data.Post.Description,
			&data.Post.Score, &data.Post.UserID,
			&data.Post.CategoryID, &data.Post.Created,
			&data.User.ID, &data.User.Login,
			&data.Category.Name)
		if nil != err {
			fmt.Println("scan: ", err)
			return nil, err
		}
		posts = append(posts, data)
	}
	return posts, nil
}

func (repo *PostsRepo) Add(post *Post) (*string, error) {
	fmt.Println("Repo post: add post")

	result, err := repo.DB.Exec(`INSERT INTO post 
	(id, title, type, description, score, user_id, category_id, created) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		post.ID, post.Title, post.Type, post.Description, post.Score, post.UserID, post.CategoryID, post.Created)

	if err != nil {
		return nil, err
	}
	rows, err := result.RowsAffected()
	fmt.Printf("affected rows : %d\n", rows)
	if err != nil {
		return nil, err
	}
	return &post.ID, nil
}

func (repo *PostsRepo) Delete(id string) (bool, error) {
	fmt.Println("Repo post: delete post")

	result, err := repo.DB.Exec(`DELETE FROM post WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected != 1 {
		return false, fmt.Errorf("wrong deleted rows: %d", affected)
	}
	return true, nil
}

func (repo *PostsRepo) UpVote(id string) (bool, error) {
	fmt.Println("Repo post: upvote")

	result, err := repo.DB.Exec(`UPDATE post SET score = score + 1 WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected != 1 {
		return false, fmt.Errorf("wrong update rows: %d for id post: %s", affected, id)
	}
	return true, nil
}

func (repo *PostsRepo) DownVote(id string) (bool, error) {
	fmt.Println("Repo post: downvote")
	result, err := repo.DB.Exec(`UPDATE post SET score = IF(score = 0, score, score - 1) WHERE id = ?`, id)
	if err != nil {
		return false, err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return false, err
	}
	return true, nil
}
