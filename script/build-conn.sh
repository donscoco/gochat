 docker build --build-arg IRONHEAD_PWD=${IRONHEAD_PWD} --build-arg IRONHEAD_OSS_SECRET=${IRONHEAD_OSS_SECRET} -t donscoco/conn_engine:v1 -f deploy/conn_engine/Dockerfile .