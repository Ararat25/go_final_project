.PHONY: run
run:
	cd cmd && ./server.bin

.PHONY: build
build:
	cd cmd && go build -o server.bin

.PHONY: start
start: build run

.PHONY: test
test:
	go test -count=1 ./tests

.PHONY: format
format:
	goimports -w .

# Указываем переменную для файла .env
ENV_FILE := .env

# Импортируем переменные из .env файла
include $(ENV_FILE)
export $(shell sed 's/=.*//' $(ENV_FILE))

IMAGE_NAME := todo-list

.PHONY: docker_build
docker_build:
	docker build --tag $(IMAGE_NAME):v1.0 .

# Команда для запуска Docker контейнера
.PHONY: docker_run
docker_run:
	docker run -p $(TODO_PORT):$(TODO_PORT) --rm -it $(IMAGE_NAME):v1.0

.PHONY: all
all: docker_build docker_run