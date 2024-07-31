# Dockerfile
FROM golang:latest
LABEL authors="ashilmira"

WORKDIR /app
ENV CONFIG_PATH=./config/local.yaml

#COPY go.mod go.sum ./
COPY . ./

RUN go mod download

RUN go build -o server cmd/main/main.go

EXPOSE 8080

ENTRYPOINT ["./server"]

