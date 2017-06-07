DOCKER_IMAGE=cbio/websocket

all: go docker

docker:
	docker build -t ${DOCKER_IMAGE} .

clean:
	docker rmi ${DOCKER_IMAGE}
	rm -f ./bin/server

go:
	docker run -it --rm -v ${CURDIR}:/go golang:alpine go build -o ./bin/server *.go
