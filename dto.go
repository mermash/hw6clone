package main

import "fmt"

type AuthorDTO struct {
	UserName string `json:"username"`
	ID       string `json:"id"`
}

type VoteDTO struct {
	User string `json:"user"`
	Vote uint32 `json:"vote"`
}

type CommentDTO struct {
	Author  *AuthorDTO `json:"author"`
	Body    string     `json:"body"`
	Created string     `json:"created,datetime"`
	ID      string     `json:"id"`
}

type PostDTO struct {
	ID               string        `json:"id"`
	Author           *AuthorDTO    `json:"author"`
	Category         string        `json:"category"`
	Comments         []*CommentDTO `json:"comments"`
	Created          string        `json:"created,datetime"`
	Score            uint32        `json:"score"`
	Text             string        `json:"text"`
	Title            string        `json:"title"`
	Type             string        `json:"type"`
	UpVotePercentage uint          `json:"upvotepercentage"`
	Votes            []*VoteDTO    `json:"votes"`
	Views            uint32        `json:"views"`
}

type ErrorDTO struct {
	ID          string `json:"id"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type PostRequestDTO struct {
	Category string `json:"category"`
	Type     string `json:"type"`
	Title    string `json:"title"`
	Text     string `json:"text"`
}

type CommentRequestDTO struct {
	Comment string `json:"comment"`
}

type LoginDTO struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

type DTOConverter struct {
	CommentRepo CommentRepoI
	VoteRepo    VoteRepoI
}

func (converter *DTOConverter) PostConvertToDTO(data *PostComplexData) (*PostDTO, error) {
	postDTO := &PostDTO{
		ID: data.Post.ID,
		Author: &AuthorDTO{
			UserName: data.User.Login,
			ID:       data.User.ID,
		},
		Category:         data.Category.Name,
		Comments:         []*CommentDTO{},
		Created:          data.Post.Created,
		Score:            data.Post.Score,
		Text:             data.Post.Description,
		Title:            data.Post.Title,
		Type:             data.Post.Type,
		UpVotePercentage: 0,
		Votes:            []*VoteDTO{},
		Views:            0,
	}

	postIds := make([]string, 0, 1)
	postIds = append(postIds, data.Post.ID)
	comments, err := converter.CommentRepo.GetCommentsByPostIds(postIds)
	if nil != err {
		return nil, err
	}

	postDTO.Comments = converter.CommentsConvertToDTO(comments[data.Post.ID])

	votes, err := converter.VoteRepo.GetVotesByPostIds(postIds)
	if nil != err {
		return nil, err
	}
	postDTO.Votes = converter.VotesConvertToDTO(votes[data.Post.ID])

	return postDTO, nil
}

func (converter *DTOConverter) CommentsConvertToDTO(data []*CommentComplexData) []*CommentDTO {
	commentsDTO := []*CommentDTO{}
	for _, comment := range data {
		commentDTO := &CommentDTO{
			Author: &AuthorDTO{
				UserName: comment.User.Login,
				ID:       comment.User.ID,
			},
			Body:    comment.Comment.Body,
			Created: comment.Comment.Created,
			ID:      comment.Comment.ID,
		}
		commentsDTO = append(commentsDTO, commentDTO)
	}
	return commentsDTO
}

func (converter *DTOConverter) VotesConvertToDTO(data []*Vote) []*VoteDTO {
	votesDTO := []*VoteDTO{}
	for _, vote := range data {
		voteDTO := &VoteDTO{
			User: vote.UserID,
			Vote: vote.Vote,
		}
		votesDTO = append(votesDTO, voteDTO)
	}
	return votesDTO
}

func (converter *DTOConverter) PostsConvertToDTO(data []*PostComplexData) ([]*PostDTO, error) {
	postsDTO := []*PostDTO{}
	postIds := make([]string, 0, 10)
	for _, post := range data {
		postIds = append(postIds, post.Post.ID)
		postDTO := &PostDTO{
			ID: post.Post.ID,
			Author: &AuthorDTO{
				UserName: post.User.Login,
				ID:       post.User.ID,
			},
			Category:         post.Category.Name,
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
		postsDTO = append(postsDTO, postDTO)
	}
	if len(postIds) > 0 {
		comments, err := converter.CommentRepo.GetCommentsByPostIds(postIds)
		if nil != err {
			fmt.Println("get comments: ", err)
			return nil, err
		}
		votes, err := converter.VoteRepo.GetVotesByPostIds(postIds)
		if nil != err {
			fmt.Println("get votes: ", err)
			return nil, err
		}
		for _, post := range postsDTO {
			post.Comments = converter.CommentsConvertToDTO(comments[post.ID])
			post.Votes = converter.VotesConvertToDTO(votes[post.ID])
		}
	}

	return postsDTO, nil
}
