COMMIT?=$(shell git rev-parse --short HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

.PHONY: mockgen
mockgen:
	@echo "-- generate mocks"
	go install github.com/golang/mock/mockgen@v1.6.0
	export PATH=$PATH:$HOME/go/bin
	mockgen -source=posts_handlers.go -destination=posts_handlers_mock.go -package=main
	mockgen -source=user_handlers.go -destination=user_handlers_mock.go -package=main

.PHONY: test
test:
	@echo "-- run the tests"
	go test -v

.PHONY: test-cover
test-cover:
	@ecgo "-- generate info about test covering"
	go test -v -coverprofile=tests_cover.out && go tool cover -html=tests_cover.out -o tests_cover.html && rm tests_cover.out
	open -a Safari ./tests_cover.html

.PHONY: start
start:
	@echo "-- start app commit=${COMMIT} build_time=${BUILD_TIME}"
	docker compose up

.PHONY: stop
stop:
	@echo "-- stop app"
	docker compose down
	docker rmi mermash/redditclone-app

.PHONY: push
push:
	@echo "-- push to github ${MESSAGE}"
	git add .
	git status
	git commit -m "${MESSAGE}"
	git push
