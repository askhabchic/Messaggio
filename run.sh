#!/bin/bash

go build -o server cmd/main/main.go

docker-compose up -d