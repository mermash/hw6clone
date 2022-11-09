package main

import (
	"database/sql"
	"fmt"
	"strings"
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

func (repo *PostsRepo) getCommentsByPostIds(postIds []string) (map[string][]*CommentDTO, error) {

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
	comments := map[string][]*CommentDTO{}
	for rows.Next() {
		comment := &Comment{}
		user := &User{}
		err := rows.Scan(&comment.ID, &comment.PostId, &comment.Body, &comment.Created, &user.ID, &user.Login)
		if nil != err {
			fmt.Println("get comments scan:", err)
			return nil, err
		}
		fmt.Println("get comment for post id", comment, comment.PostId)
		if _, ok := comments[comment.PostId]; !ok {
			comments[comment.PostId] = []*CommentDTO{}
		}
		commentDTO := &CommentDTO{
			Author: &AuthorDTO{
				UserName: user.Login,
				ID:       user.ID,
			},
			Body:    comment.Body,
			Created: comment.Created,
			ID:      comment.ID,
		}
		comments[comment.PostId] = append(comments[comment.PostId], commentDTO)
	}
	fmt.Println("comments: ", comments)
	return comments, nil
}

func (repo *PostsRepo) getVotesByPostIds(postIds []string) (map[string][]*VoteDTO, error) {
	lenPostId := len(postIds)
	placeHolders := make([]string, 0, lenPostId)
	args := make([]interface{}, 0, lenPostId)
	for _, id := range postIds {
		placeHolders = append(placeHolders, "?")
		args = append(args, id)
	}
	query := `
	SELECT post_id, user_id, vote 
	FROM vote 
	WHERE post_id IN (` + strings.Join(placeHolders, ",") + `)`
	fmt.Println("get votes sql query: ", query)
	rows, err := repo.DB.Query(query, args...)
	if nil != err {
		fmt.Println("get votes query: ", err)
		return nil, err
	}
	defer rows.Close()
	votes := map[string][]*VoteDTO{}
	for rows.Next() {
		vote := &Vote{}
		err := rows.Scan(&vote.PostID, &vote.UserID, &vote.Vote)
		if nil != err {
			fmt.Println("get votes scan: ", err)
			return nil, err
		}
		if _, ok := votes[vote.PostID]; !ok {
			votes[vote.PostID] = []*VoteDTO{}
		}
		voteDTO := &VoteDTO{
			User: vote.UserID,
			Vote: vote.Vote,
		}
		votes[vote.PostID] = append(votes[vote.PostID], voteDTO)
	}
	return votes, err
}

func (repo *PostsRepo) PostsConvertToDTO(rows *sql.Rows) ([]*PostDTO, error) {
	posts := []*PostDTO{}
	postIds := make([]string, 0, 10)
	for rows.Next() {
		post := &Post{}
		user := &User{}
		category := &Category{}
		err := rows.Scan(&post.ID, &post.Title, &post.Type, &post.Description, &post.Score, &post.UserID, &post.CategoryID, &post.Created,
			&user.ID, &user.Login,
			&category.Name)
		if nil != err {
			fmt.Println("scan: ", err)
			return nil, err
		}
		postIds = append(postIds, post.ID)
		postDTO := &PostDTO{
			ID: post.ID,
			Author: &AuthorDTO{
				UserName: user.Login,
				ID:       user.ID,
			},
			Category:         category.Name,
			Comments:         []*CommentDTO{},
			Created:          post.Created,
			Score:            post.Score,
			Text:             post.Description,
			Title:            post.Title,
			Type:             post.Type,
			UpVotePercentage: 0,
			Votes:            []*VoteDTO{},
			Views:            0,
		}
		posts = append(posts, postDTO)
	}
	if len(postIds) > 0 {
		comments, err := repo.getCommentsByPostIds(postIds)
		if nil != err {
			fmt.Println("get comments: ", err)
			return nil, err
		}
		votes, err := repo.getVotesByPostIds(postIds)
		if nil != err {
			fmt.Println("get votes: ", err)
			return nil, err
		}
		for _, post := range posts {
			post.Comments = comments[post.ID]
			post.Votes = votes[post.ID]
		}
	}

	return posts, nil
}

func (repo *PostsRepo) PostConvertToDTO(row *sql.Row) (*PostDTO, error) {
	post := &Post{}
	user := &User{}
	category := &Category{}
	err := row.Scan(&post.ID, &post.Title, &post.Type, &post.Description, &post.Score, &post.UserID, &post.CategoryID, &post.Created,
		&user.ID, &user.Login,
		&category.Name)
	if nil != err {
		return nil, err
	}
	postDTO := &PostDTO{
		ID: post.ID,
		Author: &AuthorDTO{
			UserName: user.Login,
			ID:       user.ID,
		},
		Category:         category.Name,
		Comments:         []*CommentDTO{},
		Created:          post.Created,
		Score:            post.Score,
		Text:             post.Description,
		Title:            post.Title,
		Type:             post.Type,
		UpVotePercentage: 0,
		Votes:            []*VoteDTO{},
		Views:            0,
	}

	postIds := make([]string, 0, 1)
	postIds = append(postIds, post.ID)
	comments, err := repo.getCommentsByPostIds(postIds)
	if nil != err {
		return nil, err
	}
	votes, err := repo.getVotesByPostIds(postIds)
	if nil != err {
		return nil, err
	}
	if commentsByPostId, ok := comments[post.ID]; ok {
		postDTO.Comments = commentsByPostId
	}
	if votesByPostId, ok := votes[post.ID]; ok {
		postDTO.Votes = votesByPostId
	}

	return postDTO, nil
}

func (repo *PostsRepo) GetAll() ([]*PostDTO, error) {
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
	return repo.PostsConvertToDTO(rows)
}

func (repo *PostsRepo) GetById(id string) (*PostDTO, error) {
	fmt.Println("Repo post: get by id post")

	row := repo.DB.QueryRow(`
	SELECT 
	post.id AS post_id, title, type, description, 
	score, user_id, category_id, post.created AS post_created,
	user.id AS user_id, user.login,
	category.name AS category_name
	FROM post 
	LEFT JOIN user ON user.id = post.user_id
	LEFT JOIN category ON category.id = post.category_id
	WHERE post.id = ?`, id)

	return repo.PostConvertToDTO(row)
}

func (repo *PostsRepo) GetByCategoryName(categoryName string) ([]*PostDTO, error) {
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
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return repo.PostsConvertToDTO(rows)
}

func (repo *PostsRepo) GetByUserLogin(userLogin string) ([]*PostDTO, error) {
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
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return repo.PostsConvertToDTO(rows)
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
