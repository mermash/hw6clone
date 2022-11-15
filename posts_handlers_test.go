package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

/*
mockgen -source=posts_handlers.go -destination=posts_handlers_mock.go -package=main
*/
var (
	multipleComplexData = []*PostComplexData{
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

	postsDTO = []*PostDTO{
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

	multipleExpectation = `[{"id":"dc1e2f25-76a5-4aac-9212-96e2121c16f1","author":{"username":"mer","id":"522cd619-841f-43d5-866d-f880e5f48d18"},"category":"fashion","comments":[],"created":"2022-11-09T19:51:42Z","score":1,"text":"test fashion","title":"test fashion","type":"text","upvotepercentage":0,"votes":[],"views":0}]`
	singleExpectation   = `{"id":"dc1e2f25-76a5-4aac-9212-96e2121c16f1","author":{"username":"mer","id":"522cd619-841f-43d5-866d-f880e5f48d18"},"category":"fashion","comments":[],"created":"2022-11-09T19:51:42Z","score":1,"text":"test fashion","title":"test fashion","type":"text","upvotepercentage":0,"votes":[],"views":0}`
)

func TestList(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dictionaryRepoMock := NewMockDictionaryRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	service := &PostsHandler{
		PostsRepo:      postsRepoMock,
		CommentRepo:    commentRepoMock,
		DictionaryRepo: dictionaryRepoMock,
		DTOConverter:   dtoConverterMock,
	}

	// success
	postsRepoMock.EXPECT().GetAll().Return(multipleComplexData, nil)
	dtoConverterMock.EXPECT().PostsConvertToDTO(multipleComplexData).Return(postsDTO, nil)
	req := httptest.NewRequest("GET", "/api/posts/", nil)
	w := httptest.NewRecorder()
	service.List(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, multipleExpectation) {
		t.Errorf("there aren't match; want: %#v; have: %#v", multipleExpectation, bodyStr)
		return
	}

	//getAll result error
	postsRepoMock.EXPECT().GetAll().Return(nil, fmt.Errorf("db_error"))
	req = httptest.NewRequest("GET", "/api/posts/", nil)
	w = httptest.NewRecorder()
	service.List(w, req)
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}

	//converter error
	postsRepoMock.EXPECT().GetAll().Return(multipleComplexData, nil)
	dtoConverterMock.EXPECT().PostsConvertToDTO(multipleComplexData).Return(nil, fmt.Errorf("converter error"))
	req = httptest.NewRequest("GET", "/api/posts/", nil)
	w = httptest.NewRecorder()
	service.List(w, req)
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status 500, got %d", resp.StatusCode)
		return
	}
}

func TestGetById(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dictionaryRepoMock := NewMockDictionaryRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	service := &PostsHandler{
		PostsRepo:      postsRepoMock,
		DTOConverter:   dtoConverterMock,
		DictionaryRepo: dictionaryRepoMock,
		CommentRepo:    commentRepoMock,
	}
	var postId string = "dc1e2f25-76a5-4aac-9212-96e2121c16f1"
	urlVars := map[string]string{
		"POST_ID": "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
	}

	//success
	postsRepoMock.EXPECT().GetById(postId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(postsDTO[0], nil)
	req := httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	service.GetById(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, singleExpectation) {
		t.Errorf("it's not match; want :%#v; have: %#v", singleExpectation, bodyStr)
		return
	}

	//db error
	postsRepoMock.EXPECT().GetById(postId).Return(nil, fmt.Errorf("db_error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1", nil)
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	service.GetById(w, req)
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status code 500; got: %d", resp.StatusCode)
		return
	}

	//converter error
	postsRepoMock.EXPECT().GetById(postId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(nil, fmt.Errorf("converter error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1", nil)
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	service.GetById(w, req)
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected resp status code 500; got: %d", resp.StatusCode)
		return
	}
}

func TestGetByCategoryName(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dictionaryRepoMock := NewMockDictionaryRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	service := &PostsHandler{
		PostsRepo:      postsRepoMock,
		DTOConverter:   dtoConverterMock,
		DictionaryRepo: dictionaryRepoMock,
		CommentRepo:    commentRepoMock,
	}
	var categoryName string = "fashion"
	urlVars := map[string]string{
		"CATEGORY_NAME": "fashion",
	}

	//successed
	postsRepoMock.EXPECT().GetByCategoryName(categoryName).Return(multipleComplexData, nil)
	dtoConverterMock.EXPECT().PostsConvertToDTO(multipleComplexData).Return(postsDTO, nil)
	req := httptest.NewRequest("GET", "/api/posts/fashion", nil)
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	service.GetByCategoryName(w, req)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, multipleExpectation) {
		t.Errorf("it is not matched; want: %#v; have: %#v", multipleExpectation, bodyStr)
		return
	}

	//repository error
	postsRepoMock.EXPECT().GetByCategoryName(categoryName).Return(nil, fmt.Errorf("db_error"))
	req = httptest.NewRequest("GET", "/api/posts/fashion", nil)
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	service.GetByCategoryName(w, req)
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected statuscode 500; got: %d", resp.StatusCode)
		return
	}

	//converter error
	postsRepoMock.EXPECT().GetByCategoryName(categoryName).Return(multipleComplexData, nil)
	dtoConverterMock.EXPECT().PostsConvertToDTO(multipleComplexData).Return(nil, fmt.Errorf("converter error"))
	req = httptest.NewRequest("GET", "/api/posts/fashion", nil)
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	service.GetByCategoryName(w, req)
	resp = w.Result()
	if resp.StatusCode != 500 {
		t.Errorf("expected statuscode 500; got: %d", resp.StatusCode)
		return
	}
}

func TestAdd(t *testing.T) {
	log.SetOutput(io.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dictionaryRepoMock := NewMockDictionaryRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	timeGetterMock := NewMockTimeGetterI(ctrl)
	uuidGetterMock := NewMockUUIDGetterI(ctrl)
	service := &PostsHandler{
		PostsRepo:      postsRepoMock,
		DTOConverter:   dtoConverterMock,
		DictionaryRepo: dictionaryRepoMock,
		CommentRepo:    commentRepoMock,
		TimeGetter:     timeGetterMock,
		UUIDGetter:     uuidGetterMock,
	}

	post := &Post{
		ID:          "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
		Title:       "test fashion",
		Type:        "text",
		Description: "test fashion",
		Score:       1,
		UserID:      "522cd619-841f-43d5-866d-f880e5f48d18",
		CategoryID:  1,
		Created:     "2022-11-09T19:51:42Z",
	}
	lastID := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"
	reqBody := `{"category":"fashion","type":"text","title":"test fashion","text":"test fashion"}`
	categoryName := "fashion"
	category := &Category{
		ID:   1,
		Name: categoryName,
	}

	postsRepoMock.EXPECT().Add(post).Return(&lastID, nil)
	postsRepoMock.EXPECT().GetById(lastID).Return(multipleComplexData[0], nil)
	dictionaryRepoMock.EXPECT().GetCategoryByName(categoryName).Return(category, nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(postsDTO[0], nil)
	timeGetterMock.EXPECT().GetCreated().Return("2022-11-09T19:51:42Z")
	uuidGetterMock.EXPECT().GetUUID().Return(lastID)

	req := httptest.NewRequest("POST", "/api/posts", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	sess := &Session{
		ID:     "123",
		UserID: "522cd619-841f-43d5-866d-f880e5f48d18",
	}
	ctx := context.WithValue(req.Context(), sessionKey, sess)
	service.Add(w, req.WithContext(ctx))

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	bodyStr := string(body)

	if !reflect.DeepEqual(bodyStr, singleExpectation) {
		t.Errorf("it's not matched; want: %#v; have: %#v", singleExpectation, bodyStr)
		return
	}
}
