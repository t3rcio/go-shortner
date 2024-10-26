# syntax=docker/dockerfile:1

FROM golang
WORKDIR /code
COPY . /code/
ENTRYPOINT [ "go", "run", "server_http.go", "--redisserver", "redis-db" ] 

