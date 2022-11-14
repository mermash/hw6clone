package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type PostsHandler struct {
	Tmpl           *template.Template
	PostsRepo      PostRepoI
	CommentRepo    CommentRepoI
	DictionaryRepo DictionaryRepoI
	VoteRepo       VoteRepoI
	Logger         *log.Logger
}

var ScoreDefault uint32 = 1

func NewPostsHandler(db *sql.DB, templates *template.Template) *PostsHandler {
	return &PostsHandler{
		Tmpl:           templates,
		PostsRepo:      NewPostsRepo(db),
		CommentRepo:    NewCommentRepo(db),
		VoteRepo:       NewVoteRepo(db),
		DictionaryRepo: NewDictionaryRepo(db),
		Logger:         nil,
	}
}

func (h *PostsHandler) GetById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["POST_ID"]
	data, err := h.PostsRepo.GetById(id)
	if nil != err {
		fmt.Println("can't get post by id", err)
		jsonError(w, http.StatusInternalServerError, "can't get post by id")
	}

	postDTO, err := PostConvertToDTO(h.CommentRepo, h.VoteRepo, data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert post to dto")
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
	}

	postsDTO, err := PostsConvertToDTO(h.CommentRepo, h.VoteRepo, data)
	if err != nil {
		fmt.Println("can't convert posts to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
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

	postsDTO, err := PostsConvertToDTO(h.CommentRepo, h.VoteRepo, data)
	if err != nil {
		fmt.Println("can't convert posts to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
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

	body, err := ioutil.ReadAll(r.Body)
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
		ID:          uuid.NewString(),
		Title:       requestData.Title,
		Type:        requestData.Type,
		Description: requestData.Text,
		Score:       ScoreDefault,
		UserID:      sess.UserID,
		CategoryID:  uint(category.ID),
		Created:     time.Now().Format(time.RFC3339),
	}

	lastID, err := h.PostsRepo.Add(newPost)

	fmt.Println("add post id", lastID, *lastID)

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

	postDTO, err := PostConvertToDTO(h.CommentRepo, h.VoteRepo, data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
	}

	jsonResponse(w, postDTO)
}

func (h *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["POST_ID"]

	isDeleted, err := h.PostsRepo.Delete(id)

	if nil != err || !isDeleted {
		jsonError(w, http.StatusInternalServerError, "can't delete post, err")
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

	postUpdatedDTO, err := PostConvertToDTO(h.CommentRepo, h.VoteRepo, data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
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
	}

	data, err := h.PostsRepo.GetById(postId)

	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't get updated post")
	}

	postUpdatedDTO, err := PostConvertToDTO(h.CommentRepo, h.VoteRepo, data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
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
	}

	data, err := h.PostsRepo.GetById(postId)

	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't get updated post")
	}

	postUpdatedDTO, err := PostConvertToDTO(h.CommentRepo, h.VoteRepo, data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
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

	body, err := ioutil.ReadAll(r.Body)
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
		ID:      uuid.NewString(),
		Body:    commentRequest.Comment,
		PostId:  postId,
		UserId:  sess.UserID,
		Created: time.Now().Format(time.RFC3339),
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

	postUpdatedDTO, err := PostConvertToDTO(h.CommentRepo, h.VoteRepo, data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
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
	}

	fmt.Println("Delete comment")

	data, err := h.PostsRepo.GetById(postId)

	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't get updated post")
	}

	postUpdatedDTO, err := PostConvertToDTO(h.CommentRepo, h.VoteRepo, data)
	if err != nil {
		fmt.Println("can't convert post to dto", err)
		jsonError(w, http.StatusInternalServerError, "can't convert to dto")
	}

	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postUpdatedDTO)
}
