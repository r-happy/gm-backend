FROM golang:latest

WORKDIR /app

COPY ./src /app

RUN go mod tidy

EXPOSE 1323

CMD ["go", "run", "."]