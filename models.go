package main

type Post struct {
	ID          string
	Title       string
	Type        string
	Description string
	Score       uint32
	UserID      string
	CategoryID  uint
	Created     string
}

type User struct {
	ID       string
	Login    string
	Password string
}

type Comment struct {
	ID      string
	Body    string
	PostId  string
	UserId  string
	Created string
}

type Vote struct {
	PostID string
	UserID string
	Vote   uint32
}

type Category struct {
	ID   uint32
	Name string
}

type PostComplexData struct {
	Post
	User
	Category
}

type CommentComplexData struct {
	Comment
	User
}
