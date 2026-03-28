#include .env
#export $(shell sed 's/=.*//' .env)

default: help

lint:
	golangci-lint run

build:
	docker build -t fasthttp-stream-server .

run:
	docker run --rm  -p 8080:8080 fasthttp-stream-server

help:
	@echo "Available commands:"
	@echo "  make lint — check linting"
