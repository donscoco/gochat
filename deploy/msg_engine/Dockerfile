FROM golang
MAINTAINER donscoco
WORKDIR /go/src/gochat
COPY . /go/src/gochat/
#EXPOSE 9980
RUN /bin/bash /go/src/gochat/script/build-msg-in-docker.sh
CMD ["/go/src/gochat/bin/msg_engine"]
