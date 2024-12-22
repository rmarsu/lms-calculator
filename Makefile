IMAGE_NAME=my-go-app
CONTAINER_NAME=my-go-container

.PHONY: build run clean

build:
	docker build -t $(IMAGE_NAME) .

run: build
	docker run --rm --name $(CONTAINER_NAME) $(IMAGE_NAME)

clean:
	docker rmi $(IMAGE_NAME)
