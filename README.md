# go-talk
A chat server writen in Golang

## How to run

docker build --file Dockerfile . --tag go-talk
docker run --publish 8080:8080 --volume $(pwd)/data:/data --detach go-talk