FROM golang:1.18.2-alpine3.16 AS build

RUN mkdir -p /home/app

COPY . /home/app

WORKDIR /home/app

CMD  go build && ./hw5_redditclone
#CMD go run .