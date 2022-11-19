package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
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

	postsDTOWithComments = []*PostDTO{
		{
			ID: "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
			Author: &AuthorDTO{
				UserName: "mer",
				ID:       "522cd619-841f-43d5-866d-f880e5f48d18",
			},
			Category: "fashion",
			Comments: []*CommentDTO{
				{
					Author: &AuthorDTO{
						UserName: "mer",
						ID:       "522cd619-841f-43d5-866d-f880e5f48d18",
					},
					Body:    "test comment fashion",
					Created: "2022-11-10T11:24:44Z",
					ID:      "dbed62a8-79c5-43bd-9594-92cddeb261ac",
				},
			},
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

	multipleExpectation           = `[{"id":"dc1e2f25-76a5-4aac-9212-96e2121c16f1","author":{"username":"mer","id":"522cd619-841f-43d5-866d-f880e5f48d18"},"category":"fashion","comments":[],"created":"2022-11-09T19:51:42Z","score":1,"text":"test fashion","title":"test fashion","type":"text","upvotepercentage":0,"votes":[],"views":0}]`
	singleExpectation             = `{"id":"dc1e2f25-76a5-4aac-9212-96e2121c16f1","author":{"username":"mer","id":"522cd619-841f-43d5-866d-f880e5f48d18"},"category":"fashion","comments":[],"created":"2022-11-09T19:51:42Z","score":1,"text":"test fashion","title":"test fashion","type":"text","upvotepercentage":0,"votes":[],"views":0}`
	singleExpectationWithComments = `{"id":"dc1e2f25-76a5-4aac-9212-96e2121c16f1","author":{"username":"mer","id":"522cd619-841f-43d5-866d-f880e5f48d18"},"category":"fashion","comments":[{"author":{"username":"mer","id":"522cd619-841f-43d5-866d-f880e5f48d18"},"body":"test comment fashion","created":"2022-11-10T11:24:44Z","id":"dbed62a8-79c5-43bd-9594-92cddeb261ac"}],"created":"2022-11-09T19:51:42Z","score":1,"text":"test fashion","title":"test fashion","type":"text","upvotepercentage":0,"votes":[],"views":0}`

	sess = &Session{
		ID:     "123",
		UserID: "522cd619-841f-43d5-866d-f880e5f48d18",
	}
)

func TestList(t *testing.T) {
	log.SetOutput(io.Discard)
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
	body, _ := io.ReadAll(resp.Body)
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
	log.SetOutput(io.Discard)
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
	body, _ := io.ReadAll(resp.Body)
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
	log.SetOutput(io.Discard)
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
	body, _ := io.ReadAll(resp.Body)
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

	//success
	postsRepoMock.EXPECT().Add(post).Return(&lastID, nil)
	postsRepoMock.EXPECT().GetById(lastID).Return(multipleComplexData[0], nil)
	dictionaryRepoMock.EXPECT().GetCategoryByName(categoryName).Return(category, nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(postsDTO[0], nil)
	timeGetterMock.EXPECT().GetCreated().Return("2022-11-09T19:51:42Z")
	uuidGetterMock.EXPECT().GetUUID().Return(lastID)

	req := httptest.NewRequest("POST", "/api/posts", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	ctx := context.WithValue(req.Context(), sessionKey, sess)
	service.Add(w, req.WithContext(ctx))

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)

	if !reflect.DeepEqual(bodyStr, singleExpectation) {
		t.Errorf("it's not matched; want: %#v; have: %#v", singleExpectation, bodyStr)
		return
	}

	//dictionary error
	dictionaryRepoMock.EXPECT().GetCategoryByName(categoryName).Return(nil, fmt.Errorf("dictionary error"))
	req = httptest.NewRequest("POST", "/api/posts", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	service.Add(w, req.WithContext(ctx))

	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code, got : %d", resp.StatusCode)
		return
	}

	//add error
	dictionaryRepoMock.EXPECT().GetCategoryByName(categoryName).Return(category, nil)
	uuidGetterMock.EXPECT().GetUUID().Return(lastID)
	timeGetterMock.EXPECT().GetCreated().Return("2022-11-09T19:51:42Z")
	postsRepoMock.EXPECT().Add(post).Return(nil, fmt.Errorf("add error"))
	req = httptest.NewRequest("POST", "/api/posts", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	service.Add(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code, got : %d", resp.StatusCode)
		return
	}

	//get by id error
	postsRepoMock.EXPECT().Add(post).Return(&lastID, nil)
	postsRepoMock.EXPECT().GetById(lastID).Return(nil, fmt.Errorf("get by id error"))
	dictionaryRepoMock.EXPECT().GetCategoryByName(categoryName).Return(category, nil)
	timeGetterMock.EXPECT().GetCreated().Return("2022-11-09T19:51:42Z")
	uuidGetterMock.EXPECT().GetUUID().Return(lastID)
	req = httptest.NewRequest("POST", "/api/posts", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	service.Add(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code, got : %d", resp.StatusCode)
		return
	}

	//converter error
	postsRepoMock.EXPECT().Add(post).Return(&lastID, nil)
	postsRepoMock.EXPECT().GetById(lastID).Return(multipleComplexData[0], nil)
	dictionaryRepoMock.EXPECT().GetCategoryByName(categoryName).Return(category, nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(nil, fmt.Errorf("converter error"))
	timeGetterMock.EXPECT().GetCreated().Return("2022-11-09T19:51:42Z")
	uuidGetterMock.EXPECT().GetUUID().Return(lastID)
	req = httptest.NewRequest("POST", "/api/posts", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	service.Add(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 status code, got : %d", resp.StatusCode)
		return
	}
}

func TestDelete(t *testing.T) {
	log.SetOutput(io.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	service := &PostsHandler{
		PostsRepo:    postsRepoMock,
		DTOConverter: dtoConverterMock,
		CommentRepo:  commentRepoMock,
	}
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"
	urlVars := map[string]string{
		"POST_ID": "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
	}

	//success
	expect := `{"message": "success"}`
	postsRepoMock.EXPECT().Delete(postId).Return(true, nil)
	req := httptest.NewRequest("DELETE", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1", nil)
	req = mux.SetURLVars(req, urlVars)
	w := httptest.NewRecorder()
	service.Delete(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, expect) {
		t.Errorf("it's not matched; want: %#v; have: %#v", expect, bodyStr)
		return
	}

	//query error
	postsRepoMock.EXPECT().Delete(postId).Return(false, fmt.Errorf("db_error"))
	req = httptest.NewRequest("DELETE", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.Delete(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}
}

func TestUpVote(t *testing.T) {
	log.SetOutput(io.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	service := &PostsHandler{
		PostsRepo:    postsRepoMock,
		DTOConverter: dtoConverterMock,
		CommentRepo:  commentRepoMock,
	}
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"
	urlVars := map[string]string{
		"POST_ID": "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
	}

	//success
	postsRepoMock.EXPECT().UpVote(postId).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(postsDTO[0], nil)
	req := httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/upvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w := httptest.NewRecorder()
	service.UpVote(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, singleExpectation) {
		t.Errorf("it's not matched; want: %#v; have: %#v", singleExpectation, bodyStr)
		return
	}

	//query error
	postsRepoMock.EXPECT().UpVote(postId).Return(false, fmt.Errorf("upvote db_error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/upvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.UpVote(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//get by id error
	postsRepoMock.EXPECT().UpVote(postId).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postId).Return(nil, fmt.Errorf("get by id error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/upvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.UpVote(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//converter error
	postsRepoMock.EXPECT().UpVote(postId).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(nil, fmt.Errorf("cconverter error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/upvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.UpVote(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}
}

func TestDownVote(t *testing.T) {
	log.SetOutput(io.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	service := &PostsHandler{
		PostsRepo:    postsRepoMock,
		DTOConverter: dtoConverterMock,
		CommentRepo:  commentRepoMock,
	}
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"
	urlVars := map[string]string{
		"POST_ID": "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
	}

	//success
	postsRepoMock.EXPECT().DownVote(postId).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(postsDTO[0], nil)
	req := httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/downvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w := httptest.NewRecorder()
	service.DownVote(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, singleExpectation) {
		t.Errorf("it's not matched; want: %#v; have: %#v", singleExpectation, bodyStr)
		return
	}

	//query errir
	postsRepoMock.EXPECT().DownVote(postId).Return(false, fmt.Errorf("downvote query error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/downvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.DownVote(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//query error
	postsRepoMock.EXPECT().DownVote(postId).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postId).Return(nil, fmt.Errorf("get by id error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/downvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.DownVote(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//converter error
	postsRepoMock.EXPECT().DownVote(postId).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(nil, fmt.Errorf("converter error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/downvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.DownVote(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}
}

func TestUnVote(t *testing.T) {
	log.SetOutput(io.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	service := &PostsHandler{
		PostsRepo:    postsRepoMock,
		DTOConverter: dtoConverterMock,
		CommentRepo:  commentRepoMock,
	}
	postId := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"
	urlVars := map[string]string{
		"POST_ID": "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
	}

	//success
	postsRepoMock.EXPECT().DownVote(postId).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(postsDTO[0], nil)
	req := httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/unvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w := httptest.NewRecorder()
	service.UnVote(w, req)
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, singleExpectation) {
		t.Errorf("it's not matched; want: %#v; have: %#v", singleExpectation, bodyStr)
		return
	}

	//query error
	postsRepoMock.EXPECT().DownVote(postId).Return(false, fmt.Errorf("unvote query error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/unvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.UnVote(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//get by id error
	postsRepoMock.EXPECT().DownVote(postId).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postId).Return(nil, fmt.Errorf("get by id error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/unvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.UnVote(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//converter error
	postsRepoMock.EXPECT().DownVote(postId).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(nil, fmt.Errorf("converter error"))
	req = httptest.NewRequest("GET", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/unvote", nil)
	req = mux.SetURLVars(req, urlVars)
	w = httptest.NewRecorder()
	service.UnVote(w, req)
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}
}

func TestAddComment(t *testing.T) {
	log.SetOutput(io.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	timeGetterMock := NewMockTimeGetterI(ctrl)
	uuidGetterMock := NewMockUUIDGetterI(ctrl)
	service := &PostsHandler{
		PostsRepo:    postsRepoMock,
		DTOConverter: dtoConverterMock,
		CommentRepo:  commentRepoMock,
		TimeGetter:   timeGetterMock,
		UUIDGetter:   uuidGetterMock,
	}

	lastID := "dbed62a8-79c5-43bd-9594-92cddeb261ac"
	newComment := &Comment{
		ID:      "dbed62a8-79c5-43bd-9594-92cddeb261ac",
		Body:    "test comment fashion",
		PostId:  "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
		UserId:  "522cd619-841f-43d5-866d-f880e5f48d18",
		Created: "2022-11-10T11:24:44Z",
	}
	reqBody := `{"comment":"test comment fashion"}`
	urlVars := map[string]string{
		"POST_ID": "dc1e2f25-76a5-4aac-9212-96e2121c16f1",
	}

	//success
	commentRepoMock.EXPECT().Add(newComment).Return(&lastID, nil)
	postsRepoMock.EXPECT().GetById(newComment.PostId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(postsDTOWithComments[0], nil)
	timeGetterMock.EXPECT().GetCreated().Return(newComment.Created)
	uuidGetterMock.EXPECT().GetUUID().Return(lastID)

	req := httptest.NewRequest("POST", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1", strings.NewReader(reqBody))
	w := httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	ctx := context.WithValue(req.Context(), sessionKey, sess)
	service.AddComment(w, req.WithContext(ctx))

	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, singleExpectationWithComments) {
		t.Errorf("it's not matched; want: %#v; have: %#v", singleExpectationWithComments, bodyStr)
		return
	}

	//query error
	timeGetterMock.EXPECT().GetCreated().Return(newComment.Created)
	uuidGetterMock.EXPECT().GetUUID().Return(lastID)
	commentRepoMock.EXPECT().Add(newComment).Return(nil, fmt.Errorf("add query error"))
	req = httptest.NewRequest("POST", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	service.AddComment(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//get by id error
	timeGetterMock.EXPECT().GetCreated().Return(newComment.Created)
	uuidGetterMock.EXPECT().GetUUID().Return(lastID)
	commentRepoMock.EXPECT().Add(newComment).Return(&lastID, nil)
	postsRepoMock.EXPECT().GetById(newComment.PostId).Return(nil, fmt.Errorf("get by id error"))
	req = httptest.NewRequest("POST", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	service.AddComment(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//converter error
	timeGetterMock.EXPECT().GetCreated().Return(newComment.Created)
	uuidGetterMock.EXPECT().GetUUID().Return(lastID)
	commentRepoMock.EXPECT().Add(newComment).Return(&lastID, nil)
	postsRepoMock.EXPECT().GetById(newComment.PostId).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(nil, fmt.Errorf("converter error"))
	req = httptest.NewRequest("POST", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1", strings.NewReader(reqBody))
	w = httptest.NewRecorder()
	req = mux.SetURLVars(req, urlVars)
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	service.AddComment(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}
}

func TestDeleteComment(t *testing.T) {
	log.SetOutput(io.Discard)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postsRepoMock := NewMockPostRepoI(ctrl)
	commentRepoMock := NewMockCommentRepoI(ctrl)
	dtoConverterMock := NewMockDTOConverterI(ctrl)
	timeGetterMock := NewMockTimeGetterI(ctrl)
	uuidGetterMock := NewMockUUIDGetterI(ctrl)
	service := &PostsHandler{
		PostsRepo:    postsRepoMock,
		DTOConverter: dtoConverterMock,
		CommentRepo:  commentRepoMock,
		TimeGetter:   timeGetterMock,
		UUIDGetter:   uuidGetterMock,
	}

	commentID := "dbed62a8-79c5-43bd-9594-92cddeb261ac"
	postID := "dc1e2f25-76a5-4aac-9212-96e2121c16f1"
	urlVars := map[string]string{
		"POST_ID":    postID,
		"COMMENT_ID": commentID,
	}

	//success
	commentRepoMock.EXPECT().Delete(commentID).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postID).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(postsDTO[0], nil)
	req := httptest.NewRequest("DELETE", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/dbed62a8-79c5-43bd-9594-92cddeb261ac", nil)
	req = mux.SetURLVars(req, urlVars)
	ctx := context.WithValue(req.Context(), sessionKey, sess)
	w := httptest.NewRecorder()
	service.DeleteComment(w, req.WithContext(ctx))
	resp := w.Result()
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	if !reflect.DeepEqual(bodyStr, singleExpectation) {
		t.Errorf("it's not match; want: %#v; have: %#v", singleExpectation, bodyStr)
		return
	}

	//query error
	commentRepoMock.EXPECT().Delete(commentID).Return(false, fmt.Errorf("delete query error"))
	req = httptest.NewRequest("DELETE", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/dbed62a8-79c5-43bd-9594-92cddeb261ac", nil)
	req = mux.SetURLVars(req, urlVars)
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	w = httptest.NewRecorder()
	service.DeleteComment(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//get by id error
	commentRepoMock.EXPECT().Delete(commentID).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postID).Return(nil, fmt.Errorf("get by id error"))
	req = httptest.NewRequest("DELETE", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/dbed62a8-79c5-43bd-9594-92cddeb261ac", nil)
	req = mux.SetURLVars(req, urlVars)
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	w = httptest.NewRecorder()
	service.DeleteComment(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}

	//converter error
	commentRepoMock.EXPECT().Delete(commentID).Return(true, nil)
	postsRepoMock.EXPECT().GetById(postID).Return(multipleComplexData[0], nil)
	dtoConverterMock.EXPECT().PostConvertToDTO(multipleComplexData[0]).Return(nil, fmt.Errorf("converter error"))
	req = httptest.NewRequest("DELETE", "/api/post/dc1e2f25-76a5-4aac-9212-96e2121c16f1/dbed62a8-79c5-43bd-9594-92cddeb261ac", nil)
	req = mux.SetURLVars(req, urlVars)
	ctx = context.WithValue(req.Context(), sessionKey, sess)
	w = httptest.NewRecorder()
	service.DeleteComment(w, req.WithContext(ctx))
	resp = w.Result()
	if resp.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected 500 statuscode; got %d", resp.StatusCode)
		return
	}
}
