package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var (
	multipleComplexDataByUserLogin = []*PostComplexData{
		{
			Post: Post{
				ID:          "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
				Title:       "test fashion",
				Type:        "text",
				Description: "test fashion",
				Score:       1,
				UserID:      "522cd619-841f-43d5-866d-f880e5f48d18",
				CategoryID:  1,
				Created:     "2022-11-09T19:51:42Z",
			},
			User: User{
				ID:    "522cd619-841f-43d5-866d-f880e5f48d18",
				Login: "mer",
			},
			Category: Category{
				Name: "fashion",
			}},
	}

	postsDTOByUserLogin = []*PostDTO{
		{
			ID: "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
			Author: &AuthorDTO{
				UserName: "mer",
				ID:       "522cd619-841f-43d5-866d-f880e5f48d18",
			},
			Category:         "fashion",
			Comments:         []*CommentDTO{},
			Created:          "2022-11-09T19:51:42Z",
			Score:            1,
			Text:             "test fashion",
			Title:            "test fashion",
			Type:             "text",
			UpVotePercentage: 0,
			Votes:            []*VoteDTO{},
			Views:            0,
		},
	}

	multipleExpectationByUserLogin = `[{"id":"dc1e2f25-76a5-4aac-9212-96e2121c16f1","author":{"username":"mer","id":"522cd619-841f-43d5-866d-f880e5f48d18"},"category":"fashion","comments":[],"created":"2022-11-09T19:51:42Z","score":1,"text":"test fashion","title":"test fashion","type":"text","upvotepercentage":0,"votes":[],"views":0}]`

	user = &User{
		ID:       "522cd619-841f-43d5-866d-f880e5f48d18",
		Login:    "mer",
		Password: "$2a$14$JW9COT4Lbor8tt.hUABkrueH8bSlEju3FL/g1RruLD5CvjXoFKx1a",
		Created:  "2022-11-09T19:51:42Z",
	}

	loginDTO = &LoginDTO{
		UserName: "mer",
		Password: "testtest",
	}
	reqLogin  = `{"username":"mer","password":"testtest"}`
	respLogin = `{"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoibWVyIiwiaWQiOiI2MzU2ODk4ZTc1ZTE5YTAwMDkwZTk1OGYifSwiaWF0IjoxNjY4NjkzOTQ2LCJleHAiOjE2NjkyOTg3NDZ9.lMrP0D_MG2gychZ4GaIgFsefNc7Mbu86tdpJyUU5fSg"}`

	sessUser = &Session{
		ID:     "123",
		UserID: "522cd619-841f-43d5-866d-f880e5f48d18",
	}
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoibWVyIiwiaWQiOiI2MzU2ODk4ZTc1ZTE5YTAwMDkwZTk1OGYifSwiaWF0IjoxNjY4NjkzOTQ2LCJleHAiOjE2NjkyOTg3NDZ9.lMrP0D_MG2gychZ4GaIgFsefNc7Mbu86tdpJyUU5fSg"
)

func TestGetPosts(t *testing.T) {
	log.SetOutput(io.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postsRepoMock := NewMockPostRepoI(ctrl)
	sessionManagerMock := NewMockSessionManagerI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	service := &UserHandler{
		PostsRepo:      postsRepoMock,
		SessionManager: sessionManagerMock,
		DTOConverter:   dtoConverterMock,
	}
	login := "test"
	urlVars := map[string]string{
		"USER_LOGIN": login,
	}

	//success
	postsRepoMock.EXPECT().GetByUserLogin(login).Return(multipleComplexDataByUserLogin, nil)
	dtoConverterMock.EXPECT().PostsConvertToDTO(multipleComplexDataByUserLogin).Return(postsDTOByUserLogin, nil)
	req := httptest.NewRequest("GET", "/api/user/test", nil)
	req = mux.SetURLVars(req, urlVars)
	w := httptest.NewRecorder()
	service.GetPosts(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, multipleExpectationByUserLogin) {
		t.Errorf("it's not matched; want: %#v; have: %#v", bodyStr, multipleExpectationByUserLogin)
		return
	}

	//query error
	postsRepoMock.EXPECT().GetByUserLogin(login).Return(nil, fmt.Errorf("GetByUserLogin: query error"))
	req = httptest.NewRequest("GET", "/api/user/test", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.GetPosts(w, req)
	resp = w.Result()

	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code; got: %d", resp.StatusCode)
		return
	}

	//converter error
	postsRepoMock.EXPECT().GetByUserLogin(login).Return(multipleComplexDataByUserLogin, nil)
	dtoConverterMock.EXPECT().PostsConvertToDTO(multipleComplexDataByUserLogin).Return(nil, fmt.Errorf("converter error"))
	req = httptest.NewRequest("GET", "/api/user/test", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.GetPosts(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code; got: %d", resp.StatusCode)
		return
	}

}

func TestLogin(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := NewMockUserRepoI(ctrl)
	sessionManagerMock := NewMockSessionManagerI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	timerGetterMock := NewMockTimeGetterI(ctrl)
	uuidGetterMock := NewMockUUIDGetterI(ctrl)
	userUtilsMock := NewMockUserUtilsI(ctrl)
	service := &UserHandler{
		UserRepo:       userRepoMock,
		SessionManager: sessionManagerMock,
		DTOConverter:   dtoConverterMock,
		UUIDGetter:     uuidGetterMock,
		TimeGetter:     timerGetterMock,
		UserUtils:      userUtilsMock,
	}

	//sucess
	userRepoMock.EXPECT().GetByLogin(loginDTO.UserName).Return(user, nil)
	userUtilsMock.EXPECT().CheckPasswordHash(loginDTO.Password, user.Password).Return(true)
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(reqLogin))
	w := httptest.NewRecorder()
	sessionManagerMock.EXPECT().Create(w, user).Return(sessUser, nil)
	userUtilsMock.EXPECT().GenerateJWT(user, sessUser.ID).Return(token, nil)
	service.Login(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, respLogin) {
		t.Errorf("it's not matched; want: %#v; have: %#v", respLogin, bodyStr)
		return
	}

	//query error
	userRepoMock.EXPECT().GetByLogin(loginDTO.UserName).Return(nil, fmt.Errorf("db error"))
	req = httptest.NewRequest("POST", "/api/login", strings.NewReader(reqLogin))
	w = httptest.NewRecorder()
	service.Login(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code; got: %d", resp.StatusCode)
		return
	}

	//check password error
	userRepoMock.EXPECT().GetByLogin(loginDTO.UserName).Return(user, nil)
	userUtilsMock.EXPECT().CheckPasswordHash(loginDTO.Password, user.Password).Return(false)
	req = httptest.NewRequest("POST", "/api/login", strings.NewReader(reqLogin))
	w = httptest.NewRecorder()
	service.Login(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401 status code; got: %d", resp.StatusCode)
		return
	}

	//sess create error
	userRepoMock.EXPECT().GetByLogin(loginDTO.UserName).Return(user, nil)
	userUtilsMock.EXPECT().CheckPasswordHash(loginDTO.Password, user.Password).Return(true)
	req = httptest.NewRequest("POST", "/api/login", strings.NewReader(reqLogin))
	w = httptest.NewRecorder()
	sessionManagerMock.EXPECT().Create(w, user).Return(nil, fmt.Errorf("sess create errror"))
	service.Login(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code; got: %d", resp.StatusCode)
		return
	}

	//jwt generate error
	userRepoMock.EXPECT().GetByLogin(loginDTO.UserName).Return(user, nil)
	userUtilsMock.EXPECT().CheckPasswordHash(loginDTO.Password, user.Password).Return(true)
	req = httptest.NewRequest("POST", "/api/login", strings.NewReader(reqLogin))
	w = httptest.NewRecorder()
	sessionManagerMock.EXPECT().Create(w, user).Return(sessUser, nil)
	userUtilsMock.EXPECT().GenerateJWT(user, sessUser.ID).Return("", fmt.Errorf("jwt generate error"))
	service.Login(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code; got: %d", resp.StatusCode)
		return
	}
}

func TestRegister(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userRepoMock := NewMockUserRepoI(ctrl)
	sessionManagerMock := NewMockSessionManagerI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	timerGetterMock := NewMockTimeGetterI(ctrl)
	uuidGetterMock := NewMockUUIDGetterI(ctrl)
	userUtilsMock := NewMockUserUtilsI(ctrl)
	service := &UserHandler{
		UserRepo:       userRepoMock,
		SessionManager: sessionManagerMock,
		DTOConverter:   dtoConverterMock,
		UUIDGetter:     uuidGetterMock,
		TimeGetter:     timerGetterMock,
		UserUtils:      userUtilsMock,
	}

	//sucess
	userRepoMock.EXPECT().Create(user).Return(&user.ID, nil)
	userRepoMock.EXPECT().GetById(user.ID).Return(user, nil)
	userUtilsMock.EXPECT().GeneratePasswordHash(loginDTO.Password).Return(user.Password, nil)
	uuidGetterMock.EXPECT().GetUUID().Return(user.ID)
	timerGetterMock.EXPECT().GetCreated().Return(user.Created)
	req := httptest.NewRequest("POST", "/api/register", strings.NewReader(reqLogin))
	w := httptest.NewRecorder()
	sessionManagerMock.EXPECT().Create(w, user).Return(sessUser, nil)
	userUtilsMock.EXPECT().GenerateJWT(user, sessUser.ID).Return(token, nil)
	service.Register(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, respLogin) {
		t.Errorf("it's not matched; want: %#v; have: %#v", respLogin, bodyStr)
		return
	}

	//generate password hash error
	userUtilsMock.EXPECT().GeneratePasswordHash(loginDTO.Password).Return("", fmt.Errorf("generate password hash error"))
	req = httptest.NewRequest("POST", "/api/register", strings.NewReader(reqLogin))
	w = httptest.NewRecorder()
	service.Register(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got: %d", resp.StatusCode)
		return
	}

	//create query error
	userRepoMock.EXPECT().Create(user).Return(nil, fmt.Errorf("create error"))
	userUtilsMock.EXPECT().GeneratePasswordHash(loginDTO.Password).Return(user.Password, nil)
	uuidGetterMock.EXPECT().GetUUID().Return(user.ID)
	timerGetterMock.EXPECT().GetCreated().Return(user.Created)
	req = httptest.NewRequest("POST", "/api/register", strings.NewReader(reqLogin))
	w = httptest.NewRecorder()
	service.Register(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got: %d", resp.StatusCode)
		return
	}

	//create query error
	userRepoMock.EXPECT().Create(user).Return(&user.ID, nil)
	userRepoMock.EXPECT().GetById(user.ID).Return(nil, fmt.Errorf("get by id error"))
	userUtilsMock.EXPECT().GeneratePasswordHash(loginDTO.Password).Return(user.Password, nil)
	uuidGetterMock.EXPECT().GetUUID().Return(user.ID)
	timerGetterMock.EXPECT().GetCreated().Return(user.Created)
	req = httptest.NewRequest("POST", "/api/register", strings.NewReader(reqLogin))
	w = httptest.NewRecorder()
	service.Register(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got: %d", resp.StatusCode)
		return
	}

	//session error
	userRepoMock.EXPECT().Create(user).Return(&user.ID, nil)
	userRepoMock.EXPECT().GetById(user.ID).Return(user, nil)
	userUtilsMock.EXPECT().GeneratePasswordHash(loginDTO.Password).Return(user.Password, nil)
	uuidGetterMock.EXPECT().GetUUID().Return(user.ID)
	timerGetterMock.EXPECT().GetCreated().Return(user.Created)
	req = httptest.NewRequest("POST", "/api/register", strings.NewReader(reqLogin))
	w = httptest.NewRecorder()
	sessionManagerMock.EXPECT().Create(w, user).Return(nil, fmt.Errorf("session error"))
	service.Register(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got: %d", resp.StatusCode)
		return
	}

	//session error
	userRepoMock.EXPECT().Create(user).Return(&user.ID, nil)
	userRepoMock.EXPECT().GetById(user.ID).Return(user, nil)
	userUtilsMock.EXPECT().GeneratePasswordHash(loginDTO.Password).Return(user.Password, nil)
	uuidGetterMock.EXPECT().GetUUID().Return(user.ID)
	timerGetterMock.EXPECT().GetCreated().Return(user.Created)
	req = httptest.NewRequest("POST", "/api/register", strings.NewReader(reqLogin))
	w = httptest.NewRecorder()
	sessionManagerMock.EXPECT().Create(w, user).Return(sessUser, nil)
	userUtilsMock.EXPECT().GenerateJWT(user, sessUser.ID).Return("", fmt.Errorf("generate jwt token error"))
	service.Register(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got: %d", resp.StatusCode)
		return
	}
}
