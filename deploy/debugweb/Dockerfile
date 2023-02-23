FROM golang
MAINTAINER donscoco
WORKDIR /go/src/gochat
COPY . /go/src/gochat/
EXPOSE 30010
RUN /bin/bash /go/src/gochat/script/build-in-docker.sh
CMD ["/go/src/gochat/bin/debugweb"]
