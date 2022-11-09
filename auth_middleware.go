package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	Sm SessionManager
}

func NewAuthMiddleware(sm SessionManager) AuthMiddleware {
	fmt.Println("Create authmiddleware")
	return AuthMiddleware{
		Sm: sm,
	}
}

func isAuthURL(r *http.Request) bool {
	authURLS := map[string]string{
		"/upvote":   "GET",
		"/downvote": "GET",
		"/unvote":   "GET",
	}
	authMethods := map[string]struct{}{
		"POST":   struct{}{},
		"DELETE": struct{}{},
	}
	if _, ok := authMethods[r.Method]; ok && strings.Contains(r.URL.Path, "/post") {
		return true
	}
	for path, method := range authURLS {
		if strings.Contains(r.URL.Path, path) && r.Method == method {
			return true
		}
	}
	return false
}

func (amw *AuthMiddleware) AuthMiddlewareSessionJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("isAuthorized")

		if !isAuthURL(r) {
			fmt.Println("shoudn't auth", r.URL.Path)
			next.ServeHTTP(w, r)
			return
		}

		sess, err := amw.Sm.Check(r)

		if err != nil {
			fmt.Println("error: no auth", err)
			jsonError(w, http.StatusUnauthorized, "No auth")
		}

		ctx := context.WithValue(r.Context(), sessionKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// func (amw *AuthMiddleware) AuthMiddlewareSession(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 		fmt.Println("isAuthorized by sessiondb")

// 		if !isAuthURL(r) {
// 			fmt.Println("shoudn't auth", r.URL.Path, r.Method)
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		sess, err := amw.Sm.Check(r)

// 		if err != nil {
// 			fmt.Println("error: no auth", err)
// 			jsonError(w, http.StatusUnauthorized, "No auth")
// 		}

// 		ctx := context.WithValue(r.Context(), sessionKey, sess)
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
