# Dockerfile
FROM golang:1.21-alpine
LABEL authors="ashilmira"

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /server cmd/server/main.go

EXPOSE 8080

CMD ["/server"]

