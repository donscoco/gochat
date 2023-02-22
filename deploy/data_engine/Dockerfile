FROM golang
MAINTAINER donscoco
WORKDIR /go/src/gochat
COPY . /go/src/gochat/
EXPOSE 7990

ARG IRONHEAD_PWD=deault_value
ENV IRONHEAD_PWD=$IRONHEAD_PWD

RUN /bin/bash /go/src/gochat/script/build-data-in-docker.sh
CMD ["/go/src/gochat/bin/data_engine"]
