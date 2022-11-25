package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	fmt.Println("Hello, redditclone")

	rand.Seed(time.Now().UnixNano())

	templates := template.Must(template.ParseGlob("./template/*"))

	dsn := "root:root@tcp(mysql-db:3306)/redditclone?charset=utf8mb4&interpolateParams=true"

	db, err := sql.Open("mysql", dsn)
	if nil != err {
		fmt.Println(fmt.Errorf("can't connect to db"), err.Error())
		return
	}
	db.SetMaxOpenConns(10)
	err = db.Ping()
	if nil != err {
		fmt.Println(fmt.Errorf("can't connect to db: %s", err.Error()))
		return
	}

	sm := NewSessionDBManagerJWT(db)

	postsHandler := NewPostsHandler(db)
	userHandler := NewUserHandler(db, sm)

	router := mux.NewRouter()

	router.HandleFunc("/api/register", userHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", userHandler.Login).Methods("POST")
	router.HandleFunc("/api/user/{USER_LOGIN}", userHandler.GetPosts).Methods("GET")

	router.HandleFunc("/api/posts/", postsHandler.List).Methods("GET")
	router.HandleFunc("/api/posts/{CATEGORY_NAME}", postsHandler.GetByCategoryName).Methods("GET")
	router.HandleFunc("/api/post/{POST_ID}", postsHandler.GetById).Methods("GET")
	router.HandleFunc("/api/post/{POST_ID}/upvote", postsHandler.UpVote).Methods("GET")
	router.HandleFunc("/api/post/{POST_ID}/downvote", postsHandler.DownVote).Methods("GET")
	router.HandleFunc("/api/post/{POST_ID}/unvote", postsHandler.UnVote).Methods("GET")
	router.HandleFunc("/api/posts", postsHandler.Add).Methods("POST")
	router.HandleFunc("/api/post/{POST_ID}", postsHandler.Delete).Methods("DELETE")

	router.HandleFunc("/api/post/{POST_ID}", postsHandler.AddComment).Methods("POST")
	router.HandleFunc("/api/post/{POST_ID}/{COMMENT_ID}", postsHandler.DeleteComment).Methods("DELETE")

	router.Handle("/", Index(templates))

	staticHandler := http.StripPrefix(
		"/static/",
		http.FileServer(http.Dir("./static")),
	)
	router.PathPrefix("/static/").Handler(staticHandler)

	amw := NewAuthMiddleware(sm)
	router.Use(amw.AuthMiddlewareSessionJWT)

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("zap logger error: ", err)
	}
	defer logger.Sync()
	acmw := NewAccessLoggerMiddleware(logger)
	router.Use(acmw.AccessLog)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", router)
}

func Index(templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("index")
		err := templates.ExecuteTemplate(w, "index.html", nil)
		if nil != err {
			fmt.Println(fmt.Errorf("error templates: %s", err.Error()))
		}
	}
}
