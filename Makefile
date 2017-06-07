DOCKER_IMAGE=jaschweder/websocket
DOCKER_RUN=docker run -it --rm -v ${CURDIR}:/go golang:alpine

all: go docker

docker:
	docker build -t ${DOCKER_IMAGE} .

clean:
	docker rmi ${DOCKER_IMAGE}
	rm -f ./bin/server

go:
	${DOCKER_RUN} go get -t -d
	${DOCKER_RUN} go build -o ./bin/server *.go
