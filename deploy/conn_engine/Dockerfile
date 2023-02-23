# 第一版先简单做。todo 后续完善
FROM golang
WORKDIR /go/src/gochat
COPY . /go/src/gochat/

ARG IRONHEAD_PWD=deault_value
ENV IRONHEAD_PWD=$IRONHEAD_PWD

ARG IRONHEAD_OSS_SECRET=deault_value
ENV IRONHEAD_OSS_SECRET=$IRONHEAD_OSS_SECRET

ARG GOCHAT_ENV=deault_value
ENV GOCHAT_ENV=$GOCHAT_ENV

EXPOSE 9980
EXPOSE 9990

RUN /bin/bash /go/src/gochat/script/build-conn-in-docker.sh
CMD ["/go/src/gochat/bin/conn_engine"]

