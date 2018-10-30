DOCKER_REGISTRY ?= "iocaste"

.PHONY: build
build:
	mkdir -p bin/
	go build -o bin/slack-notifier ./main.go

.PHONY: docker-build
docker-build:
	mkdir -p rootfs
	GOOS=linux GOARCH=amd64 go build -o rootfs/bin/slack-notifier ./main.go
	docker build -t $(DOCKER_REGISTRY)/slack-notifier:latest .

.PHONY: docker-push
docker-push:
	docker push $(DOCKER_REGISTRY)/slack-notifier
