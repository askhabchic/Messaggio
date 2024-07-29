# Dockerfile
FROM golang:1.21-alpine
LABEL authors="ashilmira"

#ENV GO111MODULE=on

WORKDIR ./home
#
#COPY go.mod go.sum ./
#RUN go mod download

COPY . ./home

#RUN #CONFIG_PATH=./config/local.yaml ./server

EXPOSE 8080

CMD ["CONFIG_PATH=./config/local.yaml ./home/server"]

