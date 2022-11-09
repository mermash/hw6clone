package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func jsonResponse(w http.ResponseWriter, data interface{}) {
	respBody, err := json.Marshal(data)
	if nil != err {
		jsonError(w, http.StatusInternalServerError, "can't pack response in json")
	}
	w.Write(respBody)
}

func jsonError(w http.ResponseWriter, status int, msg string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})
	w.WriteHeader(status)
	w.Write(resp)
}

func PostToDTO(post *Post) *PostDTO {
	author := &AuthorDTO{
		UserName: "test author",
		ID:       post.UserID,
	}
	postDTO := &PostDTO{
		ID:       post.ID,
		Author:   author,
		Title:    post.Title,
		Category: "music",
		Text:     post.Description,
		Created:  post.Created,
		Type:     "link",
	}
	fmt.Println("post to dto", postDTO)
	return postDTO
}

func PostsToDTO(posts []*Post) []*PostDTO {
	postsDTO := make([]*PostDTO, 0)
	for _, post := range posts {
		postDTO := PostToDTO(post)
		postsDTO = append(postsDTO, postDTO)
	}
	return postsDTO
}
