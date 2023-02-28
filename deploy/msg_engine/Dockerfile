FROM golang
MAINTAINER donscoco
WORKDIR /go/src/gochat
COPY . /go/src/gochat/
#EXPOSE 9980

ARG IRONHEAD_PWD=deault_value
ENV IRONHEAD_PWD=$IRONHEAD_PWD

ARG IRONHEAD_OSS_SECRET=deault_value
ENV IRONHEAD_OSS_SECRET=$IRONHEAD_OSS_SECRET

ARG GOCHAT_ENV=deault_value
ENV GOCHAT_ENV=$GOCHAT_ENV

RUN /bin/bash /go/src/gochat/script/build-msg-in-docker.sh
CMD ["/go/src/gochat/bin/msg_engine"]
