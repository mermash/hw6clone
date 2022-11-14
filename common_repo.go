package main

type PostRepoI interface {
	GetAll() ([]*PostComplexData, error)
	GetById(id string) (*PostComplexData, error)
	GetByCategoryName(categoryName string) ([]*PostComplexData, error)
	GetByUserLogin(userLogin string) ([]*PostComplexData, error)
	Add(post *Post) (*string, error)
	Delete(id string) (bool, error)
	UpVote(id string) (bool, error)
	DownVote(id string) (bool, error)
}

type CommentRepoI interface {
	Add(comment *Comment) (*string, error)
	Delete(id string) (bool, error)
	GetCommentsByPostIds(postIds []string) (map[string][]*CommentComplexData, error)
}

type VoteRepoI interface {
	GetVotesByPostIds(postIds []string) (map[string][]*Vote, error)
}

type DictionaryRepoI interface {
	GetCategoryByName(name string) (*Category, error)
}
