VERSION := $(shell git describe --tags --exact-match 2>/dev/null || echo latest)
DOCKERHUB_NAMESPACE ?= microkubes
IMAGE := ${DOCKERHUB_NAMESPACE}/microservice-user:${VERSION}

build:
	docker build -t ${IMAGE} .

push: build
	docker push ${IMAGE}

run: build
	docker run -p 8080:8080 ${IMAGE}
