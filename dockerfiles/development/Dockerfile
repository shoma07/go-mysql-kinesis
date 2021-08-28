FROM golang:1.17.0-alpine3.14

RUN apk add --no-cache mariadb-dev mariadb-client

WORKDIR /go/src/app

COPY go.mod go.sum ./

RUN go mod download
