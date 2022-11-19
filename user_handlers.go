package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type SessionManagerI interface {
	Check(*http.Request) (*Session, error)
	Create(http.ResponseWriter, *User) (*Session, error)
	DestroyCurrent(http.ResponseWriter, *http.Request) error
	DestroyAll(http.ResponseWriter, *User) error
}

type UserRepoI interface {
	GetById(id string) (*User, error)
	GetByLogin(login string) (*User, error)
	Create(user *User) (*string, error)
}

type UserUtilsI interface {
	GenerateJWT(user *User, sessID string) (string, error)
	GeneratePasswordHash(password string) (string, error)
	CheckPasswordHash(passwordReceived string, hash string) bool
}

type UserHandler struct {
	SessionManager SessionManagerI
	UserRepo       UserRepoI
	PostsRepo      PostRepoI
	DTOConverter   DTOConverterI
	UUIDGetter     UUIDGetterI
	TimeGetter     TimeGetterI
	UserUtils      UserUtilsI
	Logger         *log.Logger
}

func NewUserHandler(db *sql.DB, sm SessionManagerI) *UserHandler {
	return &UserHandler{
		SessionManager: sm,
		UserRepo:       NewUserRepo(db),
		PostsRepo:      NewPostsRepo(db),
		DTOConverter: &DTOConverter{
			CommentRepo: NewCommentRepo(db),
			VoteRepo:    NewVoteRepo(db),
		},
		UUIDGetter: &UUIDGetter{},
		TimeGetter: &TimeGetter{},
		UserUtils:  &UserUtils{},
		Logger:     nil,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if nil != err {
		fmt.Println("can't read request: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't read request")
		return
	}
	registerReuqest := &LoginDTO{}
	err = json.Unmarshal(body, registerReuqest)
	if nil != err {
		fmt.Println("can't unpack payload: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't unpack payload")
		return
	}
	passwordHash, err := h.UserUtils.GeneratePasswordHash(registerReuqest.Password)
	if nil != err {
		fmt.Println("can't generate a hash for the password: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't generate a hash for the password")
		return
	}
	user := &User{
		ID:       h.UUIDGetter.GetUUID(),
		Login:    registerReuqest.UserName,
		Password: passwordHash,
		Created:  h.TimeGetter.GetCreated(),
	}
	lastID, err := h.UserRepo.Create(user)
	if nil != err {
		fmt.Println("can't register a new user: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't register a new user")
		return
	}
	fmt.Println("Create user with id", lastID)
	userAdded, err := h.UserRepo.GetById(*lastID)
	if nil != err {
		fmt.Println("can't get new added user: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't get new added user")
		return
	}

	sess, err := h.SessionManager.Create(w, userAdded)

	if err != nil {
		fmt.Println("can't create session: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't create session")
		return
	}

	tokenString, err := h.UserUtils.GenerateJWT(userAdded, sess.ID)
	if nil != err {
		fmt.Println("can't generate jwt token: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't generate jwt token")
		return
	}
	data := map[string]string{
		"token": tokenString,
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	jsonResponse(w, data)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if nil != err {
		fmt.Println("can't read request body", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't read request body")
		return
	}
	loginRequest := &LoginDTO{}
	err = json.Unmarshal(data, loginRequest)
	if nil != err {
		fmt.Println("can't unpack payload: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't unpack payload")
		return
	}
	userStored, err := h.UserRepo.GetByLogin(loginRequest.UserName)
	if nil != err {
		fmt.Println("can't get user by login: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "user not found")
		return
	}
	if !h.UserUtils.CheckPasswordHash(loginRequest.Password, userStored.Password) {
		fmt.Println("invalid password")
		jsonError(w, http.StatusUnauthorized, "invalid password")
		return
	}

	sess, err := h.SessionManager.Create(w, userStored)

	if err != nil {
		fmt.Println("can't create session: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't create session")
		return
	}

	validToken, err := h.UserUtils.GenerateJWT(userStored, sess.ID)
	if nil != err {
		fmt.Println("can't generate jwt token: ", err.Error())
		jsonError(w, http.StatusInternalServerError, "can't generate jwt token")
		return
	}
	tokenData := map[string]string{
		"token": validToken,
	}
	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, tokenData)
}

func (h *UserHandler) GetPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	login := params["USER_LOGIN"]

	data, err := h.PostsRepo.GetByUserLogin(login)
	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't get posts by user login")
		return
	}

	postsDTO, err := h.DTOConverter.PostsConvertToDTO(data)
	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't convert posts by user login")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	jsonResponse(w, postsDTO)
}
