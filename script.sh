go test -v -coverprofile=db_st.out && go tool cover -html=repo_st.out -o repo_st.html && rm repo_st.out
open -a Safari ./db_st.html 

go install github.com/golang/mock/mockgen@v1.6.0

export PATH=$PATH:$HOME/go/bin
// PostsRepoI DTOConverterI should be nearby a realization
mockgen -source=posts_handlers.go -destination=posts_handlers_mock.go -package=main

go test -v -coverprofile=tests_cover.out && go tool cover -html=tests_cover.out -o tests_cover.html && rm tests_cover.out
open -a Safari ./tests_cover.html

docker compose up
docker compose down
docker rmi mermash/redditclone-app   