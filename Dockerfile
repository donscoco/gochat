FROM golang
MAINTAINER donscoco
WORKDIR /go/src/gochat
COPY . /go/src/gochat/
EXPOSE 30010
CMD ["/bin/bash", "/go/src/gochat/script/build-in-docker.sh"]
