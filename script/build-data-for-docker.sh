GOPROXY=https://goproxy.cn
go env

## 先下载好包，待会docker化时放到镜像里面
go mod vendor

docker build --build-arg IRONHEAD_PWD=${IRONHEAD_PWD} --build-arg GOCHAT_ENV=${GOCHAT_ENV} -t donscoco/data_engine:v1 -f deploy/data_engine/Dockerfile .
## docker login
#  docker build -t ${image}:${imageTag} deploy/debugweb/Dockerfile .
#  docker push ${image}:${imageTag}

