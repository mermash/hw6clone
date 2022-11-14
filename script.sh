go test -v -coverprofile=db_st.out && go tool cover -html=db_st.out -o db_st.html && rm db_st.out
open -a Safari ./db_st.html 


docker compose up
docker compose down
docker rmi mermash/redditclone-app   