package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

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

type DTOConverterI interface {
	PostConvertToDTO(data *PostComplexData) (*PostDTO, error)
	CommentsConvertToDTO(data []*CommentComplexData) []*CommentDTO
	VotesConvertToDTO(data []*Vote) []*VoteDTO
	PostsConvertToDTO(data []*PostComplexData) ([]*PostDTO, error)
}

type TimeGetterI interface {
	GetCreated() string
}

type TimeGetter struct{}

func (timer *TimeGetter) GetCreated() string {
	return time.Now().Format(time.RFC3339)
}

type UUIDGetterI interface {
	GetUUID() string
}

type UUIDGetter struct{}

func (uuidGetter *UUIDGetter) GetUUID() string {
	return uuid.NewString()
}

type PostsHandler struct {
	PostsRepo      PostRepoI
	DTOConverter   DTOConverterI
	DictionaryRepo DictionaryRepoI
	CommentRepo    CommentRepoI
	TimeGetter     TimeGetterI
	UUIDGetter     UUIDGetterI
	Logger         *log.Logger
}

var ScoreDefault uint32 = 1

func NewPostsHandler(db *sql.DB) *PostsHandler {
	commentRepo := NewCommentRepo(db)
	return &PostsHandler{
		PostsRepo: NewPostsRepo(db),
		DTOConverter: &DTOConverter{
			CommentRepo: commentRepo,
			VoteRepo:    NewVoteRepo(db),
		},
		DictionaryRepo: NewDictionaryRepo(db),
		CommentRepo:    commentRepo,
		TimeGetter:     &TimeGetter{},
		UUIDGetter:     &UUIDGetter{},
		Logger:         nil,
	}
}

func (h *PostsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["POST_ID"]
	fmt.Printf("param: %#v", params)
	data, err := h.PostsRepo.GetById(id)
	if nil != err {
		fmt.Println("can't get post by id", err)
		jsonError(w, http.StatusInternalServerError, "can't get post by id")
		return
	}

	postDTO, err := h.DTOConverter.PostConvertToDTO(data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert post to dto")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postDTO)

}

func (h *PostsHandler) GetByCategoryName(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	categoryName := params["CATEGORY_NAME"]
	data, err := h.PostsRepo.GetByCategoryName(categoryName)

	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't get posts by category")
		return
	}

	postsDTO, err := h.DTOConverter.PostsConvertToDTO(data)
	if err != nil {
		fmt.Println("can't convert posts to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postsDTO)
}

func (h *PostsHandler) List(w http.ResponseWriter, r *http.Request) {
	data, err := h.PostsRepo.GetAll()

	if nil != err {
		jsonError(w, http.StatusInternalServerError, "DB err")
		return
	}

	postsDTO, err := h.DTOConverter.PostsConvertToDTO(data)
	if err != nil {
		fmt.Println("can't convert posts to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postsDTO)
}

func (h *PostsHandler) Add(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "can't receive session")
		return
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if nil != err {
		jsonError(w, http.StatusInternalServerError, "read reqeust err")
		return
	}

	requestData := &PostRequestDTO{}
	err = json.Unmarshal(body, requestData)
	if nil != err {
		jsonError(w, http.StatusBadRequest, "can't unpack payload")
		return
	}

	category, err := h.DictionaryRepo.GetCategoryByName(requestData.Category)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "can't get category")
		return
	}

	newPost := &Post{
		ID:          h.UUIDGetter.GetUUID(),
		Title:       requestData.Title,
		Type:        requestData.Type,
		Description: requestData.Text,
		Score:       ScoreDefault,
		UserID:      sess.UserID,
		CategoryID:  uint(category.ID),
		Created:     h.TimeGetter.GetCreated(),
	}

	lastID, err := h.PostsRepo.Add(newPost)

	if nil != err {
		fmt.Println("can't add post", err)
		jsonError(w, http.StatusInternalServerError, "can't add post")
		return
	}

	data, err := h.PostsRepo.GetById(*lastID)
	if nil != err {
		fmt.Println("can't get by id the added post", err)
		jsonError(w, http.StatusInternalServerError, "can't get by id the added post")
		return
	}

	postDTO, err := h.DTOConverter.PostConvertToDTO(data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
		return
	}

	jsonResponse(w, postDTO)
}

func (h *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["POST_ID"]

	isDeleted, err := h.PostsRepo.Delete(id)

	if nil != err || !isDeleted {
		jsonError(w, http.StatusInternalServerError, "can't delete post, err")
		return
	}

	fmt.Println("Delete post", id)

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"message": "success"}`))
}

func (h *PostsHandler) UpVote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["POST_ID"]

	isUpVoted, err := h.PostsRepo.UpVote(postId)

	if nil != err || !isUpVoted {
		jsonError(w, http.StatusInternalServerError, "can't up vote")
		return
	}

	data, err := h.PostsRepo.GetById(postId)

	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't get upvoted post")
		return
	}

	postUpdatedDTO, err := h.DTOConverter.PostConvertToDTO(data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postUpdatedDTO)
}

func (h *PostsHandler) DownVote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["POST_ID"]

	_, err := h.PostsRepo.DownVote(postId)

	if nil != err {
		fmt.Println("can't down vote", err)
		jsonError(w, http.StatusInternalServerError, "can't down vote")
		return
	}

	data, err := h.PostsRepo.GetById(postId)

	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't get updated post")
		return
	}

	postUpdatedDTO, err := h.DTOConverter.PostConvertToDTO(data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postUpdatedDTO)
}

func (h *PostsHandler) UnVote(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["POST_ID"]

	_, err := h.PostsRepo.DownVote(postId)

	if nil != err {
		fmt.Println("can't down vote", err)
		jsonError(w, http.StatusInternalServerError, "can't down vote")
		return
	}

	data, err := h.PostsRepo.GetById(postId)

	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't get updated post")
		return
	}

	postUpdatedDTO, err := h.DTOConverter.PostConvertToDTO(data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postUpdatedDTO)
}

func (h *PostsHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["POST_ID"]
	sess, err := SessionFromContext(r.Context())
	if err != nil {
		fmt.Println("err: ", err)
		jsonError(w, http.StatusInternalServerError, "can't get session from context")
		return
	}
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if nil != err {
		jsonError(w, http.StatusInternalServerError, "read request err")
		return
	}
	commentRequest := &CommentRequestDTO{}
	err = json.Unmarshal(body, commentRequest)
	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't unpack payload")
		return
	}
	newComment := &Comment{
		ID:      h.UUIDGetter.GetUUID(),
		Body:    commentRequest.Comment,
		PostId:  postId,
		UserId:  sess.UserID,
		Created: h.TimeGetter.GetCreated(),
	}
	_, err = h.CommentRepo.Add(newComment)
	if nil != err {
		fmt.Println("can't add comment", err)
		jsonError(w, http.StatusInternalServerError, "can't add comment")
		return
	}
	data, err := h.PostsRepo.GetById(postId)
	if nil != err {
		fmt.Println("can't get updated post", err)
		jsonError(w, http.StatusInternalServerError, "can't get by id updated post")
		return
	}
	postUpdatedDTO, err := h.DTOConverter.PostConvertToDTO(data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postUpdatedDTO)
}

func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	postId := params["POST_ID"]
	commentId := params["COMMENT_ID"]
	isDeleted, err := h.CommentRepo.Delete(commentId)
	if nil != err || !isDeleted {
		jsonError(w, http.StatusInternalServerError, "can't delete comment, err")
		return
	}
	fmt.Println("Delete comment")
	data, err := h.PostsRepo.GetById(postId)
	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't get updated post")
		return
	}
	postUpdatedDTO, err := h.DTOConverter.PostConvertToDTO(data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postUpdatedDTO)
}
