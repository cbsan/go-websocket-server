DOCKER_IMAGE=jaschweder/websocket
DOCKER_IMAGE_BUILD=jaschweder/websocket-build
DOCKER_RUN=docker run -it --rm -v ${CURDIR}:/go ${DOCKER_IMAGE_BUILD}

all: go docker-run

docker: docker-build docker-run

docker-run:
	docker build -f Dockerfile.run -t ${DOCKER_IMAGE} .

docker-build:
	docker build -f Dockerfile.build -t ${DOCKER_IMAGE_BUILD} .

clean:
	docker rmi ${DOCKER_IMAGE} ${DOCKER_IMAGE}-build
	rm -f ./bin/server

go: docker-build
	${DOCKER_RUN} go get -t -d
	${DOCKER_RUN} go build -o ./bin/server *.go
