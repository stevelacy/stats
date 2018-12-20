GOBUILD=go build .
DOCKER_USER=stevelacy
DOCKER_TAG=latest
NAME=stats
IMAGE=$(DOCKER_USER)/$(NAME):$(DOCKER_TAG)

all: docker

build:
	$(GOBUILD)

build_linux:
	GOOS=linux GOARCH=386 CGO_ENABLED=0 $(GOBUILD)

docker:
	docker build -t $(IMAGE) .

push:
	docker push $(IMAGE)

clean:
	rm -f $(NAME)
