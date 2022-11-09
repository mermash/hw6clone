package main

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
