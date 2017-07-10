FROM golang:1.8.2

MAINTAINER cnaize

RUN go get -u  github.com/cnaize/lifland

WORKDIR /go/src/github.com/cnaize/lifland
EXPOSE 8000
CMD go run main.go
