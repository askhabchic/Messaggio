# Dockerfile
FROM golang:latest
LABEL authors="ashilmira"

WORKDIR /app

#COPY go.mod go.sum ./
COPY . ./

RUN go mod download

RUN go build -o subscriber consumer.go

EXPOSE 8081

ENTRYPOINT ["./subscriber"]

