FROM golang:1.19.4-alpine

RUN apk update && apk --no-cache add git build-base

WORKDIR /go/src/server

COPY ./server/ .

RUN go mod download

RUN go install golang.org/x/tools/gopls@latest && \ 
  go install github.com/cosmtrek/air@latest && \
  go install github.com/99designs/gqlgen@latest

CMD ["air"]