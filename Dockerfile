FROM golang:latest

COPY ./src /app

WORKDIR /app

RUN go mod tidy
RUN go build . -o back