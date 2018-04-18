NAME=main
TAG=vraju/golang-$(NAME)
VER=v1.0

all: clean build run

build:
	go run main.go

run:
	docker run -d -p 80:80 -e PORT=80 --name=$(NAME) $(TAG)
	docker run -ti --rm --link $(NAME):$(NAME) qorbani/curl

clean:
	-docker stop $(NAME)
	-docker rm $(NAME)

push:
	docker push $(TAG)
	docker push $(TAG):$(VER)
