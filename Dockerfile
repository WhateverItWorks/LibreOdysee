# Golang docker: https://hub.docker.com/_/golang
FROM golang:1.21.3-alpine
RUN apk add musl-dev
RUN apk add libc-dev
RUN apk add gcc
WORKDIR /src
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY ./ /app
EXPOSE 3000
CMD ["/app/LibreOdysee"]
