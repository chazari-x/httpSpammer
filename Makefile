IMAGE_NAME=chazari/http_spammer
TAG=latest

.PHONY: build push

build:
	docker build -t $(IMAGE_NAME):$(TAG) .

push: build
	docker push $(IMAGE_NAME):$(TAG)
