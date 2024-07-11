# Variables
DOCKER_COMPOSE = docker-compose
DOCKER_BUILD = docker build
DOCKER_RUN = docker run
DOCKER_EXEC = docker exec
DOCKER_STOP = docker-compose down

IMAGE_NAME = rashifal-go-scrapper

# Targets
build:
	$(DOCKER_BUILD) -t $(IMAGE_NAME) .

up:
	$(DOCKER_COMPOSE) up -d

down:
	$(DOCKER_STOP)

restart:
	$(DOCKER_STOP)
	$(DOCKER_COMPOSE) up -d

logs:
	$(DOCKER_COMPOSE) logs -f

exec:
	$(DOCKER_EXEC) -it $(CONTAINER_ID) /bin/bash

.PHONY: build up down restart logs exec