package main

import (
	"database/sql"
	"fmt"
	"strings"
)

type VoteRepo struct {
	DB *sql.DB
}

func NewVoteRepo(db *sql.DB) VoteRepoI {
	return &VoteRepo{
		DB: db,
	}
}

func (repo *VoteRepo) GetVotesByPostIds(postIds []string) (map[string][]*Vote, error) {
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
	votes := map[string][]*Vote{}
	for rows.Next() {
		vote := &Vote{}
		err := rows.Scan(&vote.PostID, &vote.UserID, &vote.Vote)
		if nil != err {
			fmt.Println("get votes scan: ", err)
			return nil, err
		}
		if _, ok := votes[vote.PostID]; !ok {
			votes[vote.PostID] = []*Vote{}
		}
		votes[vote.PostID] = append(votes[vote.PostID], vote)
	}
	return votes, err
}
