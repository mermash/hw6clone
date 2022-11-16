package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"reflect"
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

	//success
	postsRepoMock.EXPECT().GetByUserLogin(login).Return(multipleComplexDataByUserLogin, nil)
	dtoConverterMock.EXPECT().PostsConvertToDTO(multipleComplexDataByUserLogin).Return(postsDTOByUserLogin, nil)
	req := httptest.NewRequest("GET", "/api/user/test", nil)
	urlVars := map[string]string{
		"USER_LOGIN": login,
	}
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
}
